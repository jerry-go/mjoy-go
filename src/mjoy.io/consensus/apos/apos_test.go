package apos

import (
	"testing"
	"fmt"
	"mjoy.io/common/types"
	"math/big"
	"time"
)

func TestAposRunning(t *testing.T){
	an := newAllNodeManager()
	Config().blockDelay = 2
	Config().verifyDelay = 1
	an.init(0)
	an.run()
}

func TestApos1(t *testing.T){
	an := newAllNodeManager()
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an.init(2)
	an.run()
}

func TestApos2(t *testing.T){
	an := newAllNodeManager()
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an.init(3)
	an.run()
}

func TestApos3(t *testing.T){
	an := newAllNodeManager()
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an.init(4)
	an.run()
}

func TestRSV(t *testing.T){
	vn := newVirtualNode(1,nil)

	vnCredential := vn.makeCredential(2)
	fmt.Println("round:",vnCredential.Round.IntVal.Int64() ,
					"step:",vnCredential.Step.IntVal.Int64())
	//testStr := "testStr"
	//h := types.BytesToHash([]byte(testStr))
	//esig := vn.commonTools.ESIG(h)
	//_ = esig
	//
	//cd := CredentialData{vnCredential.Round,vnCredential.Step, vn.commonTools.GetQr_k(1)}
	cd := CredentialData{Round:types.BigInt{*big.NewInt(int64(vnCredential.Round.IntVal.Int64()))},Step:types.BigInt{*big.NewInt(int64(vnCredential.Step.IntVal.Int64()))},Quantity:vn.commonTools.GetQr_k(1)}
	sig := &SignatureVal{&vnCredential.R, &vnCredential.S, &vnCredential.V}

	str := fmt.Sprintf("testHash")
	hStr := types.BytesToHash([]byte(str))

	_ = cd
	_ ,err :=  vn.commonTools.Sender(hStr, sig)

	fmt.Println("err:",err)

}

func TestColor(t *testing.T){
	fmt.Println("\033[35mThis text is red \033[0mThis text has default color\n");
}

//m0 verify and Propagate
func TestM0(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(0)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	for i := 1 ;i <= 10; i++ {
		v := newVirtualNode(i,nil)
		m0 := v.makeCredential(i)
		an.SendDataPackToActualNode(m0)
	}

	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

//m0 filter : duplicate
func TestM0fail(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)


	v := newVirtualNode(0,nil)
	m0 := v.makeCredential(2)
	for i := 1 ;i <= 10; i++ {
		an.SendDataPackToActualNode(m0)
	}

	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

//m23 verify and Propagate
func TestM23(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	for i := 1 ;i <= 10; i++ {
		v := newVirtualNode(i,nil)
		hash := types.Hash{}
		hash[0] = 1
		m23 := &M23{
			Hash:hash,
		}
		m23.Credential = v.makeCredential(2)
		m23.Esig = v.commonTools.ESIG(m23.Hash)

		an.SendDataPackToActualNode(m23)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

//m23 duplicate message
func TestM23filter(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	for i := 1 ;i <= 10; i++ {
		v := newVirtualNode(i,nil)
		hash := types.Hash{}
		hash[0] = 1
		m23 := &M23{
			Hash:hash,
		}
		m23.Credential = v.makeCredential(2)
		m23.Esig = v.commonTools.ESIG(m23.Hash)

		an.SendDataPackToActualNode(m23)
		an.SendDataPackToActualNode(m23)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

//m23 malicious message
func TestM23filter_malicious(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	for i := 1 ;i <= 5; i++ {
		v := newVirtualNode(i,nil)
		hash := types.Hash{}
		hash[0] = 1
		m23 := &M23{
			Hash:hash,
		}
		m23.Credential = v.makeCredential(2)
		m23.Esig = v.commonTools.ESIG(m23.Hash)

		an.SendDataPackToActualNode(m23)

		hash1 := types.Hash{}
		hash1[0] = 2
		m23_1 := &M23{
			Hash:hash1,
		}
		m23_1.Credential = m23.Credential
		m23_1.Esig = m23.Esig
		//receive different  vote message m23, it must a malicious peer
		an.SendDataPackToActualNode(m23_1)

		//not honesty peer
		an.SendDataPackToActualNode(m23_1)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

//m common verify and Propagate
func TestMCommon(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		hash := types.Hash{}
		hash[0] = 1
		mcommon := &MCommon{
			Hash:hash,
		}

		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.B = 0
		mcommon.EsigB = v.commonTools.ESIG(types.BytesToHash(big.NewInt(int64(mcommon.B)).Bytes()))
		mcommon.Hash = hash
		mcommon.EsigV = v.commonTools.ESIG(hash)

		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}


func TestBp(t *testing.T) {
	bp := &BlockProposal{}
	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()
}