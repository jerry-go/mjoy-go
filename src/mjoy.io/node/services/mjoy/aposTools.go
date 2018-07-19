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
	"mjoy.io/params"
)

type PriKeyHandler interface {
	GetBasePriKey(kind reflect.Type)*ecdsa.PrivateKey
}

type BlockChainHandler interface {
	CurrentBlock() *block.Block
	CurrentBlockNum()uint64
	GetNowBlockHash()types.Hash
	InsertChain(chain block.Blocks) (int, error)
	GetBlockByNumber(number uint64) *block.Block
}

type BlockProducerHandler interface {
	GetProducerNewBlock(data *block.ConsensusData)*block.Block
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
	basePriKey *ecdsa.PrivateKey
	tmpPriKeys map[int]*ecdsa.PrivateKey
	blockChainHandler BlockChainHandler
	producerHandler BlockProducerHandler
}

func newAposTools(chanId *big.Int  , bcHandler BlockChainHandler , producerHander BlockProducerHandler)*aposTools{
	a := new(aposTools)
	//a.chanId = big.Int{}.Set(chanId)
	a.chanId = big.NewInt(0).Set(chanId)
	a.basePriKey = nil
	a.blockChainHandler = bcHandler
	a.producerHandler = producerHander
	return a
}

func (this *aposTools)CreateTmpPriKey(step int){
	this.lock.Lock()
	defer this.lock.Unlock()
	return
	if this.tmpPriKeys == nil {
		this.tmpPriKeys = make(map[int]*ecdsa.PrivateKey)
	}


	if _,ok:=this.tmpPriKeys[step];ok{
		return
	}

	tmpKey:= generatePrivateKey()
	this.tmpPriKeys[step] = tmpKey
}

func (this *aposTools)SetPriKey(priKey *ecdsa.PrivateKey){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.basePriKey = priKey
}

func (this *aposTools)Sig(pCs *apos.CredentialSign)error{
	this.lock.RLock()
	defer this.lock.RUnlock()

	_,_,_,err := pCs.Sign(this.basePriKey)
	return err
}

func (this *aposTools)Esig(pEphemeralSign *apos.EphemeralSign)error{
	this.lock.Lock()
	defer this.lock.Unlock()

	step := int(pEphemeralSign.GetStep())

	if true{
		pEphemeralSign.Signature.Init()
		_,_,_,err := pEphemeralSign.Sign(this.basePriKey)
		return err
	}else{
		if pri , ok := this.tmpPriKeys[step];ok{
			pEphemeralSign.Signature.Init()
			_,_,_,err := pEphemeralSign.Sign(pri)
			return err
		}
	}


	return errors.New(fmt.Sprintf("Not Find TmpPrivKey About:%d" , step))
}

func (this *aposTools)DelTmpKey(step int){
	this.lock.Lock()
	defer this.lock.Unlock()
	return
	if _,ok := this.tmpPriKeys[step];ok{
		delete(this.tmpPriKeys , step)
	}
}

func (this *aposTools)ClearTmpKeys(){
	this.lock.Lock()
	defer this.lock.Unlock()

	return
	this.tmpPriKeys = nil
}

func (this *aposTools)SigHash(hash types.Hash)[]byte{
	this.lock.RLock()
	defer this.lock.RUnlock()

	sig , err := crypto.Sign(hash[:] , this.basePriKey)
	if err != nil{
		logger.Error("aposTools SigErr:" , err.Error())
		return nil
	}

	return sig
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

func (this *aposTools) GetLastQrSignature() []byte{
	blk := this.blockChainHandler.CurrentBlock()
	if blk == nil {
		return nil
	}
	return blk.B_header.ConsensusData.Para
}

func (this *aposTools) GetQrSignature(round uint64) []byte {
	blk := this.blockChainHandler.GetBlockByNumber(round)
	if blk == nil {
		return nil
	}
	return blk.B_header.ConsensusData.Para
}

func (this *aposTools)GetNowBlockNum()int{
	return int(this.blockChainHandler.CurrentBlockNum())
}

func (this *aposTools)GetNextRound()int{
	return int(this.blockChainHandler.CurrentBlockNum() + 1)
}

func (this *aposTools)MakeEmptyBlock(data *block.ConsensusData)*block.Block{
	parent := this.blockChainHandler.CurrentBlock()

	header := block.CopyHeader(parent.B_header)
	header.ParentHash = parent.Hash()

	//r = r-1 + 1
	header.Number = types.NewBigInt(*big.NewInt(header.Number.IntVal.Int64() + 1))
	header.ConsensusData = *data

	header.Bloom = types.Bloom{}

	b := block.NewBlock(header , nil , nil)
	//use system private key to sign the block
	err := block.SignHeaderInner(b.B_header, block.NewBlockSigner(apos.Config().GetChainId()), params.RewordPrikey)
	if err != nil {
		logger.Error("makeEmptyBlock error:", err)
		return nil
	}
	return b
}

func (this *aposTools)GetNowBlockHash()types.Hash{
	return this.blockChainHandler.GetNowBlockHash()
}


func (this *aposTools)GetProducerNewBlock(data *block.ConsensusData)*block.Block{
	return this.producerHandler.GetProducerNewBlock(data)
}


func (this *aposTools)InsertChain(chain block.Blocks) (int, error){
	return this.blockChainHandler.InsertChain(chain)
}





