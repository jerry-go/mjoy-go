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
	"reflect"
	"errors"
	"mjoy.io/core/blockchain/block"
)

var (
	Flag_StepTest bool = false  //stop SendInner msg transfer
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

func (this *outMsgManager)SendCredential(c *CredentialSign)error{
	return nil
}

func (this *outMsgManager)PropagateCredential(c *CredentialSign)error{
	return nil
}

func (this *outMsgManager)SendInner(data dataPack)error{
	this.msgSndChan<-data
	if Flag_StepTest{

	}else{
		this.msgRcvChan<-data
	}
	return nil
}

func (this *outMsgManager)PropagateMsg(data dataPack)error{
	return nil
}


//virtual Node Manager
type allNodeManager struct {
	lock sync.RWMutex
	vituals []*virtualNode
	msger OutMsger
	allVNodeChan chan dataPack  //all virtual node's data send to allVNodeChan
	//the true apos
	actualNode *Apos
}


func newAllNodeManager()*allNodeManager{
	v := new(allNodeManager)
	return v
}

func (this *allNodeManager)init(maliciousNodeCnt int){
	this.allVNodeChan = make(chan dataPack , 1000)
	//only one msger,for virtual  nodes and actual node
	this.msger = MsgTransfer()
	this.actualNode = NewApos(this.msger , newOutCommonTools())
	this.actualNode.SetOutMsger(this.msger)
	TestPotVerifier = 1
	//all nodes
	allNodesCnt := Config().maxPotVerifiers.Uint64() -1


	//100 virtual node
	for i := 1 ;i <= int(allNodesCnt) ; i++ {
		vNode := newVirtualNode(i , this.allVNodeChan)

		if maliciousNodeCnt > 0{
			vNode.setIsHonest(false)
			maliciousNodeCnt--
		}
		this.vituals = append(this.vituals , vNode)
		go vNode.run()
	}
	go this.actualNode.Run()
	//go this.run()
	fmt.Println("allNodeManager Init ok...")
}

func (this *allNodeManager)initTestCommon(testPotVerifier int) int{
	this.allVNodeChan = make(chan dataPack , 100)
	//only one msger,for virtual  nodes and actual node
	this.msger = MsgTransfer()
	this.actualNode = NewApos(this.msger , newOutCommonTools())
	//this.actualNode.validate.fake = true
	this.actualNode.SetOutMsger(this.msger)
	TestPotVerifier = testPotVerifier


	go this.actualNode.Run()
	go this.runTestStep(0)
	fmt.Println("allNodeManager Init ok...")
	allNodesCnt := Config().maxPotVerifiers.Uint64() -1
	if Flag_StepTest{
		//why here do that,because no data from n step to n+1 step,
		allNodesCnt += 1
	}
	return int(allNodesCnt)
}

func (this *allNodeManager)initTestCommonNew(testPotVerifier int) int{
	this.msger = MsgTransfer()
	this.actualNode = NewApos(this.msger , newOutCommonTools())
	//this.actualNode.SetOutMsger(this.msger)
	TestPotVerifier = testPotVerifier


	go this.actualNode.Run()

	fmt.Println("allNodeManager Init ok...")
	allNodesCnt := Config().maxPotVerifiers.Uint64() -1

	return int(allNodesCnt)
}



func (this *allNodeManager)SendDataPackToActualNode(dp dataPack){
	this.allVNodeChan <- dp
}

func (this *allNodeManager)run(){
	for{
		select {
		//deal all data from virtual] node,just send the virtualData to the chan(will send to actual node)
		case virtualData:=<-this.allVNodeChan:
			fmt.Println("Run:allVNodeChan Type:" , reflect.TypeOf(virtualData))
			this.msger.Send2Apos(virtualData)
			//send the data from all node to actual node


		case actualData := <-this.msger.GetSubDataMsg():
			fmt.Println("Run:getSubDataMSG TYPE:" , reflect.TypeOf(actualData))
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


func (this *allNodeManager)initTestSteps(checkStep int64)int{
	this.allVNodeChan = make(chan dataPack , 1000)
	//only one msger,for virtual  nodes and actual node
	this.msger = MsgTransfer()
	this.actualNode = NewApos(this.msger , newOutCommonTools())

	this.actualNode.SetOutMsger(this.msger)
	TestPotVerifier = 1


	go this.actualNode.Run()
	go this.runTestStep(checkStep)
	fmt.Println("allNodeManager Init ok...")
	allNodesCnt := Config().maxPotVerifiers.Uint64() -1
	if Flag_StepTest{
		//why here do that,because no data from n step to n+1 step,
		allNodesCnt += 1
	}
	return int(allNodesCnt)
}


func (this *allNodeManager)runTestStep(checkStep int64){
	for {
		select {
		//deal all data from virtual] node,just send the virtualData to the chan(will send to actual node)
		case virtualData:=<-this.allVNodeChan:
			//send the data from all node to actual node
			this.msger.Send2Apos(virtualData)
			//this.msger.msgRcvChan <- virtualData

		case actualData := <-this.msger.GetSubDataMsg():

			switch v := actualData.(type) {
			case *CredentialSign:

			case *BlockProposal:

			case *GradedConsensus:
				if v.Credential.Step == uint64(checkStep){

					logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX,"Actual Step:",v.Credential.Step,"  ,Return:",v.Hash,COLOR_SHORT_RESET)
					return
				}

			case *BinaryByzantineAgreement:
				if v.Credential.Step == uint64(checkStep){
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX,"Actual Step:",v.Credential.Step,
						"\r\nReturn: Hash:",v.Hash,"\r\nBStatus:",v.B,COLOR_SHORT_RESET)
					return
				}

			default:
				logger.Warn("invalid message type ",reflect.TypeOf(v))
			}

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
	pri *ecdsa.PrivateKey   //stable key
	tmpPriKeys map[int]*ecdsa.PrivateKey
	signer Signer

	lock sync.RWMutex
}

func newOutCommonTools()*outCommonTools{
	o := new(outCommonTools)
	//privateKey
	o.pri = generatePrivateKey()
	//signer with chainId
	o.signer = NewSigner(big.NewInt(1))
	o.tmpPriKeys = make(map[int]*ecdsa.PrivateKey)

	return o
}

//if the nod is leader or verfier,call the method to create a key
func (this *outCommonTools)CreateTmpPriKey(step int){
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.tmpPriKeys == nil {
		this.tmpPriKeys = make(map[int]*ecdsa.PrivateKey)
	}


	if _,ok:=this.tmpPriKeys[step];ok{
		return
	}
	tmpKey := generatePrivateKey()
	this.tmpPriKeys[step] = tmpKey

	return
}

func (this *outCommonTools)Sig(pCs *CredentialSign)error{
	_,_,_,err:=pCs.sign(this.pri)
	return err
}

func (this *outCommonTools)Esig(pEphemeralSign *EphemeralSign)error{
	this.lock.RLock()
	defer this.lock.RUnlock()

	step := int(pEphemeralSign.step)
	if pri,ok:= this.tmpPriKeys[step];ok{
		//sign
		pEphemeralSign.Signature.init()
		_,_,_,err := pEphemeralSign.sign(pri)

		return err
	}
	return errors.New(fmt.Sprintf("Not Find TmpPriKey About:%d" , step))
}

func (this *outCommonTools)DelTmpKey(step int){
	this.lock.Lock()
	defer this.lock.Unlock()

	if _,ok := this.tmpPriKeys[step];ok{
		delete(this.tmpPriKeys , step)
	}
}

func (this *outCommonTools)ClearTmpKeys(){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.tmpPriKeys = nil
}

func (this *outCommonTools)SigHash(hash types.Hash)[]byte{

	sig , err := crypto.Sign(hash[:] , this.pri)
	if err != nil {
		logger.Errorf("outCommontTools SigErr:" , err.Error())
		return nil
	}
	return sig

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


func (this *outCommonTools)ESIG(hash types.Hash)(R,S,V *big.Int){
	sig , err := crypto.Sign(hash[:] , this.pri)
	if err != nil{
		logger.Errorf("outCommonTools ESIG:" , err.Error())
		return nil,nil,nil
	}

	R,S,V,err = this.signer.SignatureValues(sig)
	if err != nil{
		logger.Error(err.Error())
		R,S,V = nil,nil,nil
	}
	return
}



func (this *outCommonTools)ESigVerify(h types.Hash , sig []byte)error{
	return nil
}

func (this *outCommonTools)ESender(h types.Hash , sig []byte)(types.Address , error){
	return types.Address{} , nil
}


func (this *outCommonTools) GetLastQrSignature() []byte{
	qrKStr := "qrk=1"
	return []byte(qrKStr)

}

func (this *outCommonTools)GetNowBlockNum()int{
	return 100
}

func (this *outCommonTools)GetNextRound()int{
	return 100
}


func (this *outCommonTools)GetProducerNewBlock(data *block.ConsensusData)*block.Block{
	return nil
}

func (this *outCommonTools)MakeEmptyBlock(data *block.ConsensusData)*block.Block{
	return nil
}

func (this *outCommonTools)GetNowBlockHash()types.Hash{
	return types.Hash{}
}

func (this *outCommonTools)InsertOneBlock(b *block.Block)(int , error){
	return 0 , nil
}
func (this *outCommonTools)InsertChain(chain block.Blocks) (int, error){
	return 0 , nil
}

func (this *outCommonTools)SetPriKey(priKey *ecdsa.PrivateKey){

}
