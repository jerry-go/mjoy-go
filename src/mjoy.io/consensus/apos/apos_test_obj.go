package apos

import (
	"math/big"
	"mjoy.io/common/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"mjoy.io/utils/crypto"
	"sync"
	"fmt"
	"mjoy.io/core/blockchain/block"
	"time"
)

type outMsgManager struct {
	msgRcvChan chan dataPack    //rcv the msg from p2p,and broadcast to the nodes
	msgSndChan chan dataPack    //all node send msg
}

func newMsgManager()*outMsgManager{
	m := new(outMsgManager)
	m.msgRcvChan = make(chan dataPack , 100)//for node receiving msg
	m.msgSndChan = make(chan dataPack , 100)//for node sending msg by functions
	return m
}

func (this *outMsgManager)BroadCast(msg []byte)error{
	return nil
}

func (this *outMsgManager)GetMsg()<-chan dataPack{
	return this.msgRcvChan
}

func (this *outMsgManager)GetDataMsg()<-chan dataPack{
	return this.msgRcvChan
}

func (this *outMsgManager)SendCredential(c *CredentialSig)error{
	return nil
}

func (this *outMsgManager)PropagateCredential(c *CredentialSig)error{
	return nil
}

func (this *outMsgManager)SendMsg(data dataPack)error{
	this.msgSndChan<-data
	return nil
}

func (this *outMsgManager)PropagateMsg(data dataPack)error{
	return nil
}


//virtual Node Manager
type allNodeManager struct {
	lock sync.RWMutex
	vituals []*virtualNode
	msger *outMsgManager
	allVNodeChan chan dataPack  //all virtual node's data send to allVNodeChan
	//the true apos
	actualNode *Apos
}


func newAllNodeManager()*allNodeManager{
	v := new(allNodeManager)
	return v
}

func (this *allNodeManager)init(){
	this.allVNodeChan = make(chan dataPack , 1000)

	this.msger = newMsgManager()
	this.actualNode = NewApos(this.msger , newOutCommonTools())
	this.actualNode.outMsger = this.msger
	//100 virtual node
	for i := 1 ;i < 100 ; i++ {
		vNode := newVirtualNode(i , this.allVNodeChan)
		this.vituals = append(this.vituals , vNode)
		go vNode.run()
	}
	go this.actualNode.Run()
	go this.run()
}

func (this *allNodeManager)run(){
	for{
		select {
		case virtualData:=<-this.allVNodeChan:
			//send the data from all node to actual node
			this.msger.msgRcvChan <- virtualData

			case actualData := <-this.msger.msgSndChan:
				for _,vNode := range this.vituals{
					vNode.inChan <- actualData
				}

		}
	}
}


/*
//some out tools offered by Mjoy,such as signer and blockInfo getter
type CommonTools interface {
	//
	SIG([]byte )(R,S,V *big.Int)
	ESIG(hash types.Hash)([]byte)
	GetQr_k(k int)types.Hash
	GetNowBlockNum()int
	GetNextRound()int
}
*/
func generatePrivateKey()*ecdsa.PrivateKey{
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), reader)
	if err != nil {
		panic("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}
	return privateKeyECDSA
}

//each Node has a outCommonTools to sign or verify
type outCommonTools struct {
	pri *ecdsa.PrivateKey
	signer Signer
}

func newOutCommonTools()*outCommonTools{
	o := new(outCommonTools)
	//privateKey
	o.pri = generatePrivateKey()
	//signer with chainId
	o.signer = NewSigner(big.NewInt(1))

	return o
}

func (this *outCommonTools)SIG(hash types.Hash)(R,S,V *big.Int){

	sig , err := crypto.Sign(hash[:] , this.pri)
	if err != nil {
		logger.Errorf("outCommontTools SigErr:" , err.Error())
		return nil,nil,nil
	}
	R,S,V,err = this.signer.SignatureValues(sig)
	if err != nil {
		logger.Errorf("OutCommonTools Err:%s" , err.Error())
		return nil,nil,nil
	}
	return R,S,V

}

func (this *outCommonTools)SigVerify(h types.Hash , sig *SignatureVal)error{
	return nil
}

func (this *outCommonTools)Sender(h types.Hash , sig *SignatureVal)(types.Address , error){
	return types.Address{} , nil
}


func (this *outCommonTools)ESIG(hash types.Hash)([]byte){
	sig , err := crypto.Sign(hash[:] , this.pri)
	if err != nil{
		logger.Errorf("outCommonTools ESIG:" , err.Error())
		return nil
	}

	return sig
}

func (this *outCommonTools)ESigVerify(h types.Hash , sig []byte)error{
	return nil
}

func (this *outCommonTools)ESender(h types.Hash , sig []byte)(types.Address , error){
	return types.Address{} , nil
}


func (this *outCommonTools)GetQr_k(k int)types.Hash{
	qrKStr := "qrk=1"
	return types.BytesToHash([]byte(qrKStr))

}

func (this *outCommonTools)GetNowBlockNum()int{
	return 100
}

func (this *outCommonTools)GetNextRound()int{
	return 100
}

/*
virtual node:
what's works virtual node do:
1.generate privateKey
2.generate Br
3.generate M1 and send to Actual Node,Actual node judge the bigger one from two M1 data
4.when VN get msg from AN,add it's credential and sign the member of msg with it privateKey

*/
type virtualNode struct {

	//usefull member
	id int
	commonTools CommonTools
	inChan chan dataPack
	outChan chan dataPack
	exitChan chan interface{}

	lock sync.RWMutex
}

func newVirtualNode(id int,outChan chan dataPack)*virtualNode{
	v := new(virtualNode)
	v.commonTools = newOutCommonTools()
	v.id = id
	v.inChan = make(chan dataPack , 10)
	v.outChan = outChan
	v.exitChan = make(chan interface{} , 1)
	return v
}

//make credential
func (this *virtualNode)makeCredential(s int)*CredentialSig{
	this.lock.Lock()
	defer this.lock.Unlock()

	r := this.commonTools.GetNowBlockNum()
	k := 1

	Qr_k := this.commonTools.GetQr_k(k)
	str := fmt.Sprintf("%d%d%s" , r , k , Qr_k.Hex())
	//get sig
	R,S,V :=this.commonTools.SIG(types.BytesToHash([]byte(str)))

	c := new(CredentialSig)
	c.Round = types.BigInt{IntVal:*big.NewInt(int64(r))}
	c.Step = types.BigInt{IntVal:*big.NewInt(int64(s))}
	c.R = types.BigInt{IntVal:*R}
	c.S = types.BigInt{IntVal:*S}
	c.V = types.BigInt{IntVal:*V}

	return c
}

func (this *virtualNode)makeEmptyBlock()*block.Block{
	header := &block.Header{Number:types.NewBigInt(*big.NewInt(int64(this.commonTools.GetNowBlockNum()))),Time:types.NewBigInt(*big.NewInt(time.Now().Unix()))}
	//chainId := big.NewInt(100)
	//signer := block.NewBlockSigner(chainId)
	R,S,V := this.commonTools.SIG(header.Hash())
	header.R = &types.BigInt{*R}
	header.S = &types.BigInt{*S}
	header.V = &types.BigInt{*V}

	b := block.NewBlock(header , nil , nil)
	return b
}

func (this *virtualNode)makeM1(number int)dataPack{
	m := new(M1)
	m.Block = this.makeEmptyBlock()
	m.Credential = this.makeCredential(1)
	sigBytes := this.commonTools.ESIG(m.Block.Hash())
	m.Esig = append(m.Esig , sigBytes...)

	return m
}

func (this *virtualNode)dealM1(data dataPack)dataPack{
	m1 := data.(*M1)
	m2 := new(M23)
	m2.Credential = this.makeCredential(2)
	m2.Hash = m1.Block.Hash()
	m2.Esig = this.commonTools.ESIG(m2.Hash)

	return m2
}

func (this *virtualNode)dealM23(data dataPack)dataPack{
	m := data.(*M23)
	if 2 != m.Credential.Step.IntVal.Int64() || 3 != m.Credential.Step.IntVal.Int64() {
		return nil
	}

	if 2 == m.Credential.Step.IntVal.Int64() {
		// step 2,should make m3
		m3 := new(M23)
		m3.Credential = this.makeCredential(3)
		m3.Hash = m.Hash
		m3.Esig = this.commonTools.ESIG(m.Hash)

		return m3
	}else {
		// step 3,should make mCommon
		m4 := new(MCommon)
		m4.Credential = this.makeCredential(4)
		m4.B = 0
		m4.Hash = m.Hash

		m4.EsigV = this.commonTools.ESIG(m.Hash)
		str := fmt.Sprintf("%d" , m4.B)

		m4.EsigB = this.commonTools.ESIG(types.BytesToHash([]byte(str)))

		return m4
	}
	return nil

}

func (this *virtualNode)dealMCommon(data dataPack)dataPack{
	m := data.(*MCommon)
	mc := new(MCommon)
	mc.B = m.B
	mc.Hash = m.Hash
	mc.Credential = this.makeCredential(int(m.Credential.Step.IntVal.Int64())+1)

	mc.EsigV = this.commonTools.ESIG(mc.Hash)
	str := fmt.Sprintf("%d" , mc.B)
	mc.EsigB = this.commonTools.ESIG(types.BytesToHash([]byte(str)))

	return mc

}


//Focus:no matter what data the virtual
func (this *virtualNode)dataDeal(data dataPack)(dp dataPack){
	switch v := data.(type) {
	case *CredentialSig:
		dp = this.makeM1(int(v.Step.IntVal.Int64()))
	case *M1:
		dp = this.dealM1(v)
	case *M23:
		dp = this.dealM23(v)
	case *MCommon:
		dp = this.dealMCommon(v)
	}

	return
}






func (this *virtualNode)run(){
	for{
		select {
		case dataIn:=<-this.inChan:
			dp := this.dataDeal(dataIn)
			//the virtualNode's data should send to the actual node

			this.outChan <- dp
		case <-this.exitChan:
			return
		}
	}
}








