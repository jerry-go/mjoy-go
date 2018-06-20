package apos

import (
	"sync"
	"fmt"
	"mjoy.io/common/types"
	"math/big"
	"mjoy.io/core/blockchain/block"
	"time"
)

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
	isHonest bool
	lock sync.RWMutex
}

func newVirtualNode(id int,outChan chan dataPack)*virtualNode{
	v := new(virtualNode)
	v.commonTools = newOutCommonTools()
	v.id = id
	v.inChan = make(chan dataPack , 10)
	v.outChan = outChan
	v.exitChan = make(chan interface{} , 1)
	v.isHonest = true
	return v
}

func (this *virtualNode)setIsHonest(isHonest bool){
	this.isHonest = isHonest
}




//make credential
func (this *virtualNode)makeCredential(s int)*CredentialSig{
	this.lock.Lock()
	defer this.lock.Unlock()

	r := this.commonTools.GetNowBlockNum()
	k := 1

	Qr_k := this.commonTools.GetQr_k(k)
	str := fmt.Sprintf("testHash")
	hStr := types.BytesToHash([]byte(str))

	cd := CredentialData{Round:types.BigInt{*big.NewInt(int64(r))},Step:types.BigInt{*big.NewInt(int64(s))},Quantity:Qr_k}
	_ = cd
	//h := cd.Hash()
	h := hStr
	//get sig
	R,S,V :=this.commonTools.SIG(h)

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


	//m2.Credential.PrintInfo()
	m2.Hash = m1.Block.Hash()
	if !this.isHonest{
		m2.Hash[10] = m2.Hash[10] + 1
	}

	m2.Esig = this.commonTools.ESIG(m2.Hash)
	logger.Debug("\033[35m [V] In M1 Out M2 \033[0m")
	return m2
}

func (this *virtualNode)dealM23(data dataPack)dataPack{

	m := data.(*M23)
	if 2 != m.Credential.Step.IntVal.Int64() && 3 != m.Credential.Step.IntVal.Int64() {
		return nil
	}
	logger.Debug("\033[35m [V]In Mx step:",m.Credential.Step.IntVal.Int64(),"\033[0m ")
	if 2 == m.Credential.Step.IntVal.Int64() {
		// step 2,should make m3
		m3 := new(M23)
		m3.Credential = this.makeCredential(3)
		m3.Hash = m.Hash
		if !this.isHonest{
			m3.Hash[10] = m3.Hash[10] + 1
		}

		m3.Esig = this.commonTools.ESIG(m.Hash)
		logger.Debug("\033[35m [V]In M2 Out M3 \033[0m ")
		return m3
	}else {
		// step 3,should make mCommon

		m4 := new(MCommon)
		m4.Credential = this.makeCredential(4)
		if this.isHonest{
			m4.B = 0
			m4.Hash = m.Hash
		}else{
			m4.Hash = m.Hash
			m4.B = 1
		}


		m4.EsigV = this.commonTools.ESIG(m.Hash)
		str := fmt.Sprintf("%d" , m4.B)

		m4.EsigB = this.commonTools.ESIG(types.BytesToHash([]byte(str)))
		logger.Debug("\033[35m [V]In M3 Out M4  \033[0m ")
		return m4
	}
	return nil

}

func (this *virtualNode)dealMCommon(data dataPack)dataPack{

	m := data.(*MCommon)

	mc := new(MCommon)
	if this.isHonest{
		mc.B = m.B
	}else{
		if m.B == 0 {
			mc.B = 1
		}else{
			mc.B = 0
		}
	}


	mc.Hash = m.Hash
	mc.Credential = this.makeCredential(int(m.Credential.Step.IntVal.Int64())+1)

	mc.EsigV = this.commonTools.ESIG(mc.Hash)
	str := fmt.Sprintf("%d" , mc.B)
	mc.EsigB = this.commonTools.ESIG(types.BytesToHash([]byte(str)))
	logger.Debug("\033[35m [V]In M",m.Credential.Step.IntVal.Int64() ,
		"  Out M",int(m.Credential.Step.IntVal.Int64())+1 ,
		"  time:",time.Now().Format("2006-01-02 15:04:05"),"\033[0m ")
	return mc

}


//Focus:no matter what data the virtual
func (this *virtualNode)dataDeal(data dataPack)(dp dataPack){
	switch v := data.(type) {
	case *CredentialSig:
		//no need to deal CredentialSig
		//dp = this.makeM1(int(v.Step.IntVal.Int64()))
	case *M1:
		dp = this.dealM1(v)
	case *M23:
		dp = this.dealM23(v)
	case *MCommon:
		dp = this.dealMCommon(v)
	}

	return
}


func (this *virtualNode)stop(){
	this.exitChan<-1
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



