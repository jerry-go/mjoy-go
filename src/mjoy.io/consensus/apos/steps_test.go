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

	notHonestCnt := 2
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "NOT HONEST CNT:",notHonestCnt , COLOR_SHORT_RESET)
	for i := 1 ;i<=verifierCnt;i++{
		m2 := new(M23)
		m2.Hash = types.Hash{}
		if notHonestCnt > 0{
			m2.Hash[10] = m2.Hash[10]+1
			notHonestCnt--
		}else{
			m2.Hash[2] = m2.Hash[2]+1
		}

		v := newVirtualNode(i,nil)
		m2.Credential = v.makeCredential(2)
		m2.Esig = v.commonTools.ESIG(m2.Hash)

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
		m3 := new(M23)
		m3.Hash = types.Hash{}
		if notHonestCnt > 0{
			m3.Hash[10] = m3.Hash[10]+1
			notHonestCnt--
		}else{
			m3.Hash[2] = m3.Hash[2]+1
		}

		v := newVirtualNode(i,nil)
		m3.Credential = v.makeCredential(3)
		m3.Esig = v.commonTools.ESIG(m3.Hash)

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
		mc := new(MCommon)
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
		h := types.BytesToHash(big.NewInt(int64(mc.B)).Bytes())
		mc.EsigB = v.commonTools.ESIG(h)
		mc.EsigV = v.commonTools.ESIG(mc.Hash)

		an.SendDataPackToActualNode(mc)
	}


	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}
