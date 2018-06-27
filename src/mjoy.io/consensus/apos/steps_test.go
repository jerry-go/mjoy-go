package apos

import (
	"testing"
	"time"
	"mjoy.io/common/types"
	"math/big"
)

//instructions:
//this test file just test the steps running right or not,we just send the data which the steps need
//That is , we just need focus on the result of the test


/*
each step obj like below:
type stepInterface interface {
	sendMsg(dataPack,*Round) error
	stop()
	run(wg *sync.WaitGroup)
}





*/




func TestStep3Result(t *testing.T){
	//open the Flag_StepTest
	Flag_StepTest = true
	LessTimeDelayFlag = true
	LessTimeDelayCnt = 5

	an := newAllNodeManager()
	verifierCnt := an.initTestSteps(3)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	notHonestCnt := 1
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "NOT HONEST CNT:",notHonestCnt , COLOR_SHORT_RESET)
	for i := 1 ;i<=verifierCnt;i++{

		m2 := newGradedConsensus()
		m2.Hash = types.Hash{}
		if notHonestCnt > 0{
			m2.Hash[10] = m2.Hash[10]+1
			notHonestCnt--
		}else{
			m2.Hash[2] = m2.Hash[2]+1
		}

		v := newVirtualNode(i,nil)
		m2.Credential = v.makeCredential(2)

		//sig
		m2.Esig.round = m2.Credential.Round
		m2.Esig.step = m2.Credential.Step
		m2.Esig.val = make([]byte , 0)
		m2.Esig.val = append(m2.Esig.val  , m2.Hash[:]...)

		R,S,V := v.commonTools.ESIG(m2.Hash)

		m2.Esig.R = new(types.BigInt)
		m2.Esig.R.IntVal = *R

		m2.Esig.S = new(types.BigInt)
		m2.Esig.S.IntVal = *S

		m2.Esig.V = new(types.BigInt)
		m2.Esig.V.IntVal = *V

		an.SendDataPackToActualNode(m2)
	}


	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

func TestStep4Result(t *testing.T){
	//open the Flag_StepTest
	Flag_StepTest = true
	LessTimeDelayFlag = true
	LessTimeDelayCnt = 5

	an := newAllNodeManager()
	verifierCnt := an.initTestSteps(4)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	notHonestCnt := 2
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "NOT HONEST CNT:",notHonestCnt , COLOR_SHORT_RESET)

	for i := 1 ;i<=verifierCnt;i++{
		m3 := newGradedConsensus()
		m3.Hash = types.Hash{}
		if notHonestCnt > 0{
			m3.Hash[10] = m3.Hash[10]+1
			notHonestCnt--
		}else{
			m3.Hash[2] = m3.Hash[2]+1
		}

		v := newVirtualNode(i,nil)
		m3.Credential = v.makeCredential(3)
		//sig
		m3.Esig.round = m3.Credential.Round
		m3.Esig.step = m3.Credential.Step
		m3.Esig.val = make([]byte , 0)
		m3.Esig.val = append(m3.Esig.val , m3.Hash[:]...)

		R,S,V := v.commonTools.ESIG(m3.Hash)

		m3.Esig.R = new(types.BigInt)
		m3.Esig.R.IntVal = *R

		m3.Esig.S = new(types.BigInt)
		m3.Esig.S.IntVal = *S

		m3.Esig.V = new(types.BigInt)
		m3.Esig.V.IntVal = *V

		an.SendDataPackToActualNode(m3)
	}


	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

func TestStepCommonResult_ChangeHashAndB(t *testing.T){
	//open the Flag_StepTest
	Flag_StepTest = true
	LessTimeDelayFlag = true
	LessTimeDelayCnt = 5

	var CheckStep int64 = 12

	an := newAllNodeManager()
	verifierCnt := an.initTestSteps(CheckStep+1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	notHonestCnt := 2   //change majoraty Hash
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "NOT HONEST CNT:",notHonestCnt , COLOR_SHORT_RESET)

	for i := 1 ;i<=verifierCnt;i++{
		mc := newBinaryByzantineAgreement()
		mc.Hash = types.Hash{}

		if notHonestCnt > 0{
			mc.Hash[10] = mc.Hash[10]+1
			notHonestCnt--
		}else{
			mc.Hash[2] = mc.Hash[2]+1
		}
		mc.B = 1    //change B
		v := newVirtualNode(i,nil)
		mc.Credential = v.makeCredential(int(CheckStep))
		//b sig
		mc.EsigB.round = mc.Credential.Round
		mc.EsigB.step = mc.Credential.Step
		mc.EsigB.val = make([]byte , 0)
		mc.EsigB.val = append(mc.EsigB.val , big.NewInt(int64(mc.B)).Bytes()...)

		h := types.BytesToHash(big.NewInt(int64(mc.B)).Bytes())
		R,S,V := v.commonTools.ESIG(h)

		mc.EsigB.R = new(types.BigInt)
		mc.EsigB.R.IntVal = *R

		mc.EsigB.S = new(types.BigInt)
		mc.EsigB.S.IntVal = *S

		mc.EsigB.V = new(types.BigInt)
		mc.EsigB.V.IntVal = *V

		//v sig
		mc.EsigV.round = mc.Credential.Round
		mc.EsigV.step = mc.Credential.Step
		mc.EsigV.val = make([]byte , 0)
		mc.EsigV.val = append(mc.EsigV.val , mc.Hash[:]...)

		R,S,V = v.commonTools.ESIG(mc.Hash)

		mc.EsigV.R = new(types.BigInt)
		mc.EsigV.R.IntVal = *R

		mc.EsigV.S = new(types.BigInt)
		mc.EsigV.S.IntVal = *S

		mc.EsigV.V = new(types.BigInt)
		mc.EsigV.V.IntVal = *V
		an.SendDataPackToActualNode(mc)
	}


	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}
