package apos

import (
	"testing"
	"fmt"
	"mjoy.io/common/types"
	"math/big"
	"time"
)

func TestAposRunning(t *testing.T){
	fmt.Println("TestAposRunning.........")
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
	an.init(0)
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
	fmt.Println("round:",vnCredential.Round ,
					"step:",vnCredential.Step)

	cd := CredentialData{Round:types.BigInt{*big.NewInt(int64(vnCredential.Round))},Step:types.BigInt{*big.NewInt(int64(vnCredential.Step))},Quantity:types.Hash{}}
	sig := &SignatureVal{vnCredential.R, vnCredential.S, vnCredential.V}

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

		m23 := newGradedConsensus()
		m23.Hash = hash

		m23.Credential = v.makeCredential(2)

		m23.Esig.round = m23.Credential.Round
		m23.Esig.step = m23.Credential.Step
		m23.Esig.val = make([]byte , 0)
		m23.Esig.val = append(m23.Esig.val , m23.Hash[:]...)

		v.commonTools.CreateTmpPriKey(int(m23.Credential.Step))
		err := v.commonTools.Esig(m23.Esig)
		if err != nil{
			logger.Error(err.Error())
			return
		}
		v.commonTools.DelTmpKey(int(m23.Credential.Step))
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

		m23 := newGradedConsensus()
		m23.Hash = hash
		m23.Credential = v.makeCredential(2)

		m23.Esig.round = m23.Credential.Round
		m23.Esig.step = m23.Credential.Step

		m23.Esig.val = make([]byte , 0)
		m23.Esig.val = append(m23.Esig.val , m23.Hash[:]...)

		v.commonTools.CreateTmpPriKey(int(m23.Credential.Step))
		err := v.commonTools.Esig(m23.Esig)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(m23.Credential.Step))
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

		m23 := newGradedConsensus()
		m23.Hash = hash
		m23.Credential = v.makeCredential(2)

		m23.Esig.round = m23.Credential.Round
		m23.Esig.step = m23.Credential.Step
		m23.Esig.val = make([]byte , 0)

		m23.Esig.val = append(m23.Esig.val , m23.Hash[:]...)
		v.commonTools.CreateTmpPriKey(int(m23.Credential.Step))

		err := v.commonTools.Esig(m23.Esig)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		an.SendDataPackToActualNode(m23)

		hash1 := types.Hash{}
		hash1[0] = 2
		m23_1 := newGradedConsensus()

		m23_1.Hash = hash
		m23.Credential = m23.Credential
		m23_1.Esig = m23.Esig

		v.commonTools.DelTmpKey(int(m23.Credential.Step))

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
		mcommon := newBinaryByzantineAgreement()


		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.B = 0
		mcommon.Hash = hash

		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0 )
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))

		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0)
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}


func TestMCommon_filter_duplicate(t *testing.T){
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


		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.Hash = hash
		mcommon.B = 0

		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0 )
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))

		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil {
			logger.Error(err.Error())
			return
		}


		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0 )
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)


		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

func TestMCommon_filter_duplicate2(t *testing.T){
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


		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.Hash = hash
		mcommon.B = 0

		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0 )
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))

		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0 )
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}


		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)

		mcommonXX := newBinaryByzantineAgreement()

		mcommonXX.Credential = mcommon.Credential
		mcommonXX.B = 1
		mcommonXX.EsigB = mcommon.EsigB
		mcommonXX.Hash = hash
		mcommonXX.EsigV = mcommon.EsigV

		//logger.Info("receive different vote common message!", msg.B)
		an.SendDataPackToActualNode(mcommonXX)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

func TestMCommon_filter_duplicate3(t *testing.T){
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


		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.Hash = hash
		mcommon.B = 0

		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0 )
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))
		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil {
			logger.Error(err.Error())
			return
		}


		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0 )
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)

		mcommonXX := newBinaryByzantineAgreement()

		mcommonXX.Credential = mcommon.Credential
		mcommonXX.B = 1
		mcommonXX.EsigB = mcommon.EsigB
		mcommonXX.Hash = hash
		mcommonXX.EsigV = mcommon.EsigV

		//receive different hash in BBA message, it must a malicious peer
		an.SendDataPackToActualNode(mcommonXX)

	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

// End condition 0
func TestMCommon_EndCondition0(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)

	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()

	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0 )
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))

	err := v.commonTools.Esig(m1.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}
	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)


		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.Hash = hash
		mcommon.B = 0

		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0 )
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))

		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil {
			logger.Error(err.Error())
			return
		}


		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0 )
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
	}

	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

// End condition 0
// b =1 ignore
// vote number sum 0
func TestMCommon_EndCondition0_B1(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)

	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)

	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()

	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0 )
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))

	err := v.commonTools.Esig(m1.Esig)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3)
		mcommon.B = 1
		mcommon.Hash = hash
		//b
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)

		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))

		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return
		}


		//hash
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0)

		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))

		an.SendDataPackToActualNode(mcommon)

	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}


// End condition 1
func TestMCommon_EndCondition1(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)
	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()

	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0)

	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))

	err := v.commonTools.Esig(m1.Esig)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()


		mcommon.Credential = v.makeCredential(4 + 3 + 1)
		mcommon.B = 1
		mcommon.Hash = hash

		//sig b
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)

		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))
		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		//sig v
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0 )

		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err  = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

// End condition 1
//b = 0 ignore
// vote number sum 0
func TestMCommon_EndCondition1_b0(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)
	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()
	//sig
	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0)
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))
	err := v.commonTools.Esig(m1.Esig)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()


		mcommon.Credential = v.makeCredential(4 + 3 + 1)
		mcommon.B = 0
		mcommon.Hash = hash

		//sig B
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))

		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		//sig V
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val = make([]byte , 0)
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

// End condition
//s = 7 ignore
// vote number sum 0
func TestMCommon_EndCondition_s7_b0(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)
	m1 := newBlockProposal()
	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()

	//sig
	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0)
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))
	err := v.commonTools.Esig(m1.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3 + 2)
		mcommon.B = 0
		mcommon.Hash = hash

		//sig B
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))
		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil {
			logger.Error(err.Error())
			return
		}


		//sig V
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val= make([]byte , 0)
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

// End condition
//s = 7 ignore
// vote number sum 0
func TestMCommon_EndCondition_s7_b1(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)
	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()
	//sig
	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0)
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))
	err := v.commonTools.Esig(m1.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(4 + 3 + 2)
		mcommon.B = 1
		mcommon.Hash = hash

		//sig B
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))
		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return
		}


		//sig V
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val= make([]byte , 0)
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))

		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

// End condition max
//OK Consensus....ret: 3
func TestMCommon_EndConditionMax(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)
	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()
	//sig
	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0)
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))
	err := v.commonTools.Esig(m1.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(15)
		mcommon.B = 1
		mcommon.Hash = hash

		//sig B
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))
		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return
		}


		//sig V
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val= make([]byte , 0)
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
		an.SendDataPackToActualNode(mcommon)
	}
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
}

//verify msg common fail m + 3 message b is not equal 1
func TestMCommon_EndConditionMax_validate(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	an := newAllNodeManager()
	verifierCnt := an.initTestCommon(1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	//send m1
	v := newVirtualNode(1,nil)
	m1 := newBlockProposal()

	m1.Credential = v.makeCredential(1)
	m1.Block = an.actualNode.makeEmptyBlockForTest(m1.Credential)
	fmt.Println(m1.Block)
	hash := m1.Block.Hash()
	//sig
	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = m1.Credential.Step
	m1.Esig.val = make([]byte , 0)
	m1.Esig.val = append(m1.Esig.val , hash[:]...)

	v.commonTools.CreateTmpPriKey(int(m1.Credential.Step))
	err := v.commonTools.Esig(m1.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	v.commonTools.DelTmpKey(int(m1.Credential.Step))
	an.SendDataPackToActualNode(m1)


	for i := 1 ;i <= 4; i++ {
		v := newVirtualNode(i,nil)
		mcommon := newBinaryByzantineAgreement()

		mcommon.Credential = v.makeCredential(15)
		mcommon.B = 0
		mcommon.Hash = hash

		//sig B
		mcommon.EsigB.round = mcommon.Credential.Round
		mcommon.EsigB.step = mcommon.Credential.Step
		mcommon.EsigB.val = make([]byte , 0)
		mcommon.EsigB.val = append(mcommon.EsigB.val , big.NewInt(int64(mcommon.B)).Bytes()...)

		v.commonTools.CreateTmpPriKey(int(mcommon.Credential.Step))
		err := v.commonTools.Esig(mcommon.EsigB)
		if err != nil{
			logger.Error(err.Error())
			return
		}

		//sig V
		mcommon.EsigV.round = mcommon.Credential.Round
		mcommon.EsigV.step = mcommon.Credential.Step
		mcommon.EsigV.val= make([]byte , 0)
		mcommon.EsigV.val = append(mcommon.EsigV.val , mcommon.Hash[:]...)

		err = v.commonTools.Esig(mcommon.EsigV)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		v.commonTools.DelTmpKey(int(mcommon.Credential.Step))
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

