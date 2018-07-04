package apos

import (
	"sync"
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
	//k := 1
	//
	//Qr_k := this.commonTools.GetQr_k(k)
	//str := fmt.Sprintf("testHash")
	//hStr := types.BytesToHash([]byte(str))
	//
	//cd := CredentialData{Round:types.BigInt{*big.NewInt(int64(r))},Step:types.BigInt{*big.NewInt(int64(s))},Quantity:Qr_k}
	//_ = cd
	////h := cd.Hash()
	//h := hStr
	//get sig

	c := new(CredentialSign)
	c.Signature.init()
	c.Round = uint64(r)
	c.Step = uint64(s)

	err := this.commonTools.Sig(c)
	if err != nil{
		logger.Error(err.Error())
		return nil
	}

	return c
}

func (this *virtualNode)makeEmptyBlock()*block.Block{
	header := &block.Header{Number:types.NewBigInt(*big.NewInt(int64(this.commonTools.GetNowBlockNum()))),Time:types.NewBigInt(*big.NewInt(0))}
	//chainId := big.NewInt(100)
	//signer := block.NewBlockSigner(chainId)
	signature := MakeEmptySignature()

	sig := this.commonTools.SigHash(header.Hash())
	R,S,V,err:=signature.FillBySig(sig)
	if err != nil {
		logger.Error("makeEmptyBlock err:" , err.Error())
		return nil
	}
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

	this.commonTools.CreateTmpPriKey(int(m.Esig.step))

	err := this.commonTools.Esig(m.Esig)
	if err != nil{
		logger.Error(err.Error())
		return nil
	}

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

	this.commonTools.CreateTmpPriKey(int(m2.Credential.Step))
	err := this.commonTools.Esig(m2.Esig)
	if err!= nil{
		logger.Error(err.Error())
		return nil
	}

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

		this.commonTools.CreateTmpPriKey(int(m3.Credential.Step))
		err := this.commonTools.Esig(m3.Esig)
		if err != nil{
			logger.Error(err.Error())
			return nil
		}

		logger.Debug("\033[35m [V]In M2 Out M3 \033[0m ")
		return m3
	}else {
		// step 3,should make mCommon

		m4 := newBinaryByzantineAgreement()
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

		this.commonTools.CreateTmpPriKey(int(m4.Credential.Step))
		err := this.commonTools.Esig(m4.EsigV)
		if err != nil{
			logger.Error(err.Error())
			return nil
		}


		m4.EsigB.round = m4.Credential.Round
		m4.EsigB.step = m4.Credential.Step
		m4.EsigB.val = make([]byte , 0)
		m4.EsigB.val = append(m4.EsigB.val , big.NewInt(int64(m4.B)).Bytes()...)

		err = this.commonTools.Esig(m4.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return nil
		}

		logger.Debug("\033[35m [V]In M3 Out M4  \033[0m ")
		return m4
	}
	return nil

}

func (this *virtualNode)dealMCommon(data dataPack)dataPack{


	m := data.(*BinaryByzantineAgreement)
	mc := newBinaryByzantineAgreement()


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

	this.commonTools.CreateTmpPriKey(int(mc.Credential.Step))
	err := this.commonTools.Esig(mc.EsigV)
	if err != nil{
		logger.Error(err.Error())
		return nil
	}

	mc.EsigB.round = mc.Credential.Round
	mc.EsigB.step = mc.Credential.Step
	mc.EsigB.val = make([]byte , 0)
	mc.EsigB.val = append(mc.EsigB.val , big.NewInt(int64(mc.B)).Bytes()...)

	err = this.commonTools.Esig(mc.EsigB)
	if err != nil{
		logger.Error(err.Error())
		return nil
	}


	logger.Debug("\033[35m [V]In M",m.Credential.Step ,
		"  Out M",int(m.Credential.Step)+1 ,
		"  time:",time.Now().Format("2006-01-02 15:04:05"),"\033[0m ")
	return mc

}


//Focus:no matter what data the virtual
func (this *virtualNode)dataDeal(data dataPack)(dp dataPack){
	switch v := data.(type) {
	case *CredentialSign:
		//no need to deal CredentialSig
		//dp = this.makeM1(int(v.Step.IntVal.Int64()))
	case *BlockProposal:
		dp = this.dealM1(v)
	case *GradedConsensus:
		dp = this.dealM23(v)
	case *BinaryByzantineAgreement:
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



