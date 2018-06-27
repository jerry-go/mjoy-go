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
func (this *virtualNode)makeCredential(s int)*CredentialSign{
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

	c := new(CredentialSign)
	c.Round = uint64(r)
	c.Step = uint64(s)

	c.R = &types.BigInt{IntVal:*R}
	c.S = &types.BigInt{IntVal:*S}
	c.V = &types.BigInt{IntVal:*V}

	return c
}

func (this *virtualNode)makeEmptyBlock()*block.Block{
	header := &block.Header{Number:types.NewBigInt(*big.NewInt(int64(this.commonTools.GetNowBlockNum()))),Time:types.NewBigInt(*big.NewInt(0))}
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
	m := newBlockProposal()
	m.Block = this.makeEmptyBlock()
	m.Credential = this.makeCredential(1)

	m.Esig.round = m.Credential.Round
	m.Esig.step = m.Credential.Step
	m.Esig.val = make([]byte , 0 )
	h := m.Block.Hash()
	m.Esig.val = append(m.Esig.val , h[:]...)

	R,S,V := this.commonTools.ESIG(m.Block.Hash())

	m.Esig.R = new(types.BigInt)
	m.Esig.R.IntVal = *R

	m.Esig.S = new(types.BigInt)
	m.Esig.S.IntVal = *S

	m.Esig.V = new(types.BigInt)
	m.Esig.V.IntVal = *V



	return m
}

func (this *virtualNode)dealM1(data dataPack)dataPack{
	m1 := data.(*BlockProposal)

	m2 := newGradedConsensus()
	m2.Credential = this.makeCredential(2)


	//m2.Credential.PrintInfo()
	m2.Hash = m1.Block.Hash()
	if !this.isHonest{
		m2.Hash[10] = m2.Hash[10] + 1
	}
	m2.Esig.round = m2.Credential.Round
	m2.Esig.step = m2.Credential.Step
	m2.Esig.val = make([]byte , 0)
	m2.Esig.val = append(m2.Esig.val , m2.Hash[:]...)

	R,S,V := this.commonTools.ESIG(m2.Hash)

	m2.Esig.R = new(types.BigInt)
	m2.Esig.R.IntVal = *R

	m2.Esig.S = new(types.BigInt)
	m2.Esig.S.IntVal = *S

	m2.Esig.V = new(types.BigInt)
	m2.Esig.V.IntVal = *V

	logger.Debug("\033[35m [V] In M1 Out M2 \033[0m")
	return m2
}

func (this *virtualNode)dealM23(data dataPack)dataPack{
	m := data.(*GradedConsensus)

	if 2 != m.Credential.Step && 3 != m.Credential.Step {
		return nil
	}
	logger.Debug("\033[35m [V]In Mx step:",m.Credential.Step,"\033[0m ")
	if 2 == m.Credential.Step {
		// step 2,should make m3
		m3 := newGradedConsensus()
		m3.Credential = this.makeCredential(3)
		m3.Hash = m.Hash
		if !this.isHonest{
			m3.Hash[10] = m3.Hash[10] + 1
		}

		m3.Esig.round = m3.Credential.Round
		m3.Esig.step = m3.Credential.Step
		m3.Esig.val = make([]byte , 0 )
		m3.Esig.val = append(m3.Esig.val , m3.Hash[:]...)

		R,S,V := this.commonTools.ESIG(m3.Hash)

		m3.Esig.R = new(types.BigInt)
		m3.Esig.R.IntVal = *R

		m3.Esig.S = new(types.BigInt)
		m3.Esig.S.IntVal = *S

		m3.Esig.V = new(types.BigInt)
		m3.Esig.V.IntVal = *V

		logger.Debug("\033[35m [V]In M2 Out M3 \033[0m ")
		return m3
	}else {
		// step 3,should make mCommon
		m4 := new(BinaryByzantineAgreement)
		m4.Credential = this.makeCredential(4)
		if this.isHonest{
			m4.B = 0
			m4.Hash = m.Hash
		}else{
			m4.Hash = m.Hash
			m4.B = 1
		}

		m4.EsigV.round = m4.Credential.Round
		m4.EsigV.step  = m4.Credential.Step
		m4.EsigV.val = make([]byte , 0)
		m4.EsigV.val = append(m4.EsigV.val , m4.Hash[:]...)

		R,S,V := this.commonTools.ESIG(m4.Hash)

		m4.EsigV.R = new(types.BigInt)
		m4.EsigV.R.IntVal = *R

		m4.EsigV.S = new(types.BigInt)
		m4.EsigV.S.IntVal = *S

		m4.EsigV.V = new(types.BigInt)
		m4.EsigV.V.IntVal = *V


		m4.EsigB.round = m4.Credential.Round
		m4.EsigB.step = m4.Credential.Step
		m4.EsigB.val = make([]byte , 0)
		m4.EsigB.val = append(m4.EsigB.val , big.NewInt(int64(m4.B)).Bytes()...)
		h := types.BytesToHash(big.NewInt(int64(m4.B)).Bytes())
		m4.EsigB.val = append(m4.EsigB.val , h[:]...)

		R,S,V = this.commonTools.ESIG(types.BytesToHash(big.NewInt(int64(m4.B)).Bytes()))

		m4.EsigB.R = new(types.BigInt)
		m4.EsigB.R.IntVal = *R

		m4.EsigB.S = new(types.BigInt)
		m4.EsigB.S.IntVal = *S

		m4.EsigB.V = new(types.BigInt)
		m4.EsigB.V.IntVal = *V

		logger.Debug("\033[35m [V]In M3 Out M4  \033[0m ")
		return m4
	}
	return nil

}

func (this *virtualNode)dealMCommon(data dataPack)dataPack{


	m := data.(*BinaryByzantineAgreement)
	mc := new(BinaryByzantineAgreement)


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
	mc.Credential = this.makeCredential(int(m.Credential.Step)+1)


	mc.EsigV.round = mc.Credential.Round
	mc.EsigV.step  = mc.Credential.Step
	mc.EsigV.val = make([]byte , 0)
	mc.EsigV.val = append(mc.EsigV.val , mc.Hash[:]...)

	R,S,V := this.commonTools.ESIG(mc.Hash)

	mc.EsigV.R = new(types.BigInt)
	mc.EsigV.R.IntVal = *R

	mc.EsigV.S = new(types.BigInt)
	mc.EsigV.S.IntVal = *S

	mc.EsigV.V = new(types.BigInt)
	mc.EsigV.V.IntVal = *V


	mc.EsigB.round = mc.Credential.Round
	mc.EsigB.step = mc.Credential.Step
	mc.EsigB.val = make([]byte , 0)
	mc.EsigB.val = append(mc.EsigB.val , big.NewInt(int64(mc.B)).Bytes()...)
	h := types.BytesToHash(big.NewInt(int64(mc.B)).Bytes())
	mc.EsigB.val = append(mc.EsigB.val , h[:]...)

	R,S,V = this.commonTools.ESIG(types.BytesToHash(big.NewInt(int64(mc.B)).Bytes()))

	mc.EsigB.R = new(types.BigInt)
	mc.EsigB.R.IntVal = *R

	mc.EsigB.S = new(types.BigInt)
	mc.EsigB.S.IntVal = *S

	mc.EsigB.V = new(types.BigInt)
	mc.EsigB.V.IntVal = *V

	logger.Debug("\033[35m [V]In M",m.Credential.Step ,
		"  Out M",int(m.Credential.Step)+1 ,
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



