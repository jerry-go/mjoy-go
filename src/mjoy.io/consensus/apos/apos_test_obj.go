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
	"mjoy.io/common"
)

type outMsgManager struct {
	msgRcvChan chan dataPack    //rcv the msg from p2p,and broadcast to the nodes
	msgSndChan chan dataPack    //all node send msg
}

func newMsgManager()*outMsgManager{
	m := new(outMsgManager)
	m.msgRcvChan = make(chan dataPack , 2000)//for node receiving msg
	m.msgSndChan = make(chan dataPack , 2000)//for node sending msg by functions
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

func (this *outMsgManager)SendInner(data dataPack)error{
	this.msgSndChan<-data
	this.msgRcvChan<-data
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
	//only one msger,for virtual  nodes and actual node
	this.msger = newMsgManager()
	this.actualNode = NewApos(this.msger , newOutCommonTools())
	this.actualNode.validate.fake = true
	this.actualNode.SetOutMsger(this.msger)
	//all nodes
	allNodesCnt := 99
	notHonestNodeCnt := 35
	//100 virtual node
	for i := 1 ;i <= allNodesCnt ; i++ {
		notHonestNodeCnt--
		vNode := newVirtualNode(i , this.allVNodeChan)
		if notHonestNodeCnt > 0{
			vNode.setIsHonest(false)
		}
		this.vituals = append(this.vituals , vNode)
		go vNode.run()
	}
	go this.actualNode.Run()
	go this.run()
	fmt.Println("allNodeManager Init ok...")
}

func (this *allNodeManager)run(){
	for{
		select {
		//deal all data from virtual] node,just send the virtualData to the chan(will send to actual node)
		case virtualData:=<-this.allVNodeChan:
			//send the data from all node to actual node
			this.msger.msgRcvChan <- virtualData

		case actualData := <-this.msger.msgSndChan:
			//continue
			for _,vNode := range this.vituals{
				vNode.inChan <- actualData
			}
		case <-this.actualNode.StopCh():
			//stop all virtualNode
			for _,n := range this.vituals{
				n.stop()
			}
			//the actualNode has a result ,exit the test
			fmt.Println("exit allNodeManager......")
			return
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
	V := &big.Int{}
	V = V.Sub(&sig.V.IntVal, big.NewInt(2))
	V.Sub(V, common.Big35)
	address , err := recoverPlain(h , &sig.R.IntVal , &sig.S.IntVal , V,true)
	return address , err
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







