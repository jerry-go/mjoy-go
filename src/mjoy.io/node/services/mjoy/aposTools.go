package mjoy

import (
	"math/big"
	"crypto/ecdsa"
	"mjoy.io/consensus/apos"
	"bytes"
	"crypto/rand"
	"mjoy.io/utils/crypto"
	"sync"
	"errors"
	"fmt"
	"mjoy.io/common/types"
	"mjoy.io/common"
	"mjoy.io/core/blockchain/block"
	"reflect"
	"mjoy.io/accounts/keystore"
)

type PriKeyHandler interface {
	GetBasePriKey(kind reflect.Type)*ecdsa.PrivateKey
}

type BlockChainHandler interface {
	CurrentBlockNum()uint64
	GetNowBlockHash()types.Hash
	InsertChain(chain block.Blocks) (int, error)
}

type BlockProducerHandler interface {
	GetProducerNewBlock()*block.Block
}

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

type aposTools struct {
	lock sync.RWMutex
	chanId *big.Int
	tmpPriKeys map[int]*ecdsa.PrivateKey
	basePrivKeyHandler PriKeyHandler
	blockChainHandler BlockChainHandler
	producerHandler BlockProducerHandler
	signer apos.Signer
}

func newAposTools(chanId *big.Int , keyHandler PriKeyHandler , bcHandler BlockChainHandler , producerHander BlockProducerHandler)*aposTools{
	a := new(aposTools)
	//a.chanId = big.Int{}.Set(chanId)
	a.chanId = big.NewInt(0).Set(chanId)
	a.basePrivKeyHandler = keyHandler
	a.blockChainHandler = bcHandler
	a.producerHandler = producerHander
	return a
}



func (this *aposTools)CreateTmpPriKey(step int){
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.tmpPriKeys == nil {
		this.tmpPriKeys = make(map[int]*ecdsa.PrivateKey)
	}


	if _,ok:=this.tmpPriKeys[step];ok{
		return
	}

	tmpKey:= generatePrivateKey()
	this.tmpPriKeys[step] = tmpKey
}


func (this *aposTools)Sig(pCs *apos.CredentialSign)error{
	priKey := this.basePrivKeyHandler.GetBasePriKey(keystore.KeyStoreType)
	if priKey == nil {
		fmt.Println("******Sig PriKey == nil")
	}
	_,_,_,err := pCs.Sign(priKey)
	return err
}

func (this *aposTools)Esig(pEphemeralSign *apos.EphemeralSign)error{
	this.lock.Lock()
	defer this.lock.Unlock()

	step := int(pEphemeralSign.GetStep())

	if pri , ok := this.tmpPriKeys[step];ok{
		pEphemeralSign.Signature.Init()
		_,_,_,err := pEphemeralSign.Sign(pri)
		return err
	}

	return errors.New(fmt.Sprintf("Not Find TmpPrivKey About:%d" , step))
}

func (this *aposTools)DelTmpKey(step int){
	this.lock.Lock()
	defer this.lock.Unlock()

	if _,ok := this.tmpPriKeys[step];ok{
		delete(this.tmpPriKeys , step)
	}
}

func (this *aposTools)ClearTmpKeys(){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.tmpPriKeys = nil
}

func (this *aposTools)SigHash(hash types.Hash)(R,S,V *big.Int){
	sig , err := crypto.Sign(hash[:] , this.basePrivKeyHandler.GetBasePriKey(keystore.KeyStoreType))
	if err != nil{
		logger.Error("aposTools SigErr:" , err.Error())
		return nil , nil , nil
	}

	R,S,V , err  = this.signer.SignatureValues(sig)
	if err != nil{
		logger.Errorf("aposTools Err:%s" , err.Error())
		return nil , nil , nil
	}
	return R,S,V
}

func (this *aposTools)SigVerify(h types.Hash , sig *apos.SignatureVal)error{
	return nil
}


func (this *aposTools)Sender(h types.Hash , sig *apos.SignatureVal)(types.Address , error){
	V := &big.Int{}
	V = V.Sub(&sig.V.IntVal, big.NewInt(2))
	V.Sub(V, common.Big35)

	address , err := apos.RecoverPlain(h , &sig.R.IntVal , &sig.S.IntVal , V,true)
	return address , err
}

func (this *aposTools)ESigVerify(h types.Hash , sig []byte)error{
	return nil
}


func (this *aposTools)ESender(hash types.Hash , sig []byte)(types.Address , error){
	return types.Address{} , nil
}

func (this *aposTools)GetQr_k(k int)types.Hash{
	qrKStr := "qrk=1"
	return types.BytesToHash([]byte(qrKStr))
}

func (this *aposTools)GetNowBlockNum()int{
	return int(this.blockChainHandler.CurrentBlockNum())
}

func (this *aposTools)GetNextRound()int{
	return int(this.blockChainHandler.CurrentBlockNum() + 1)
}


func (this *aposTools)GetNowBlockHash()types.Hash{
	return this.blockChainHandler.GetNowBlockHash()
}


func (this *aposTools)GetProducerNewBlock()*block.Block{
	return this.producerHandler.GetProducerNewBlock()
}


func (this *aposTools)InsertChain(chain block.Blocks) (int, error){
	return this.blockChainHandler.InsertChain(chain)
}





