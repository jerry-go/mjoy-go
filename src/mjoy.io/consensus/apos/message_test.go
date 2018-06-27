package apos

import (
	"math/big"
	"testing"
	"fmt"
	"mjoy.io/common/types"
	"time"
	"mjoy.io/common"
)

func (s *Signature) init() {
	s.R = new(types.BigInt)
	s.S = new(types.BigInt)
	s.V = new(types.BigInt)
}
// End condition 0 for message bp bba
func TestBba_EndCondition0(t *testing.T){
	Config().blockDelay = 2
	Config().verifyDelay = 1
	Config().maxBBASteps = 12
	Config().prVerifier = 10000000000
	Config().prLeader = 10000000000
	an := newAllNodeManager()
	verifierCnt := an.initTestCommonNew(0)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_BLUE+COLOR_SUFFIX , "Verifier Cnt:" , verifierCnt , COLOR_SHORT_RESET)

	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 1
	cs.Signature.init()
	if _,_,_, err := cs.sign(priKey); err != nil {
		fmt.Println("111",err)
		return
	}

	bp := newBlockProposal()
	bp.Credential = cs
	bp.Block = an.actualNode.makeEmptyBlockForTest(bp.Credential)
	fmt.Println(bp.Block)
	hash := bp.Block.Hash()

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = hash.Bytes()
	bp.Esig.Signature.init()
	if _,_,_, err := bp.Esig.sign(priKey); err != nil {
		fmt.Println("2222",err)
	}

	//an.SendDataPackToActualNode(m1)
	msgbp := NewMsgBlockProposal(bp)
	msgbp.Send()

	for i := 1 ;i <= 4; i++ {
		time.Sleep(1 * time.Second)
		priKey := generatePrivateKey()
		cs := &CredentialSign{}
		cs.Round = 100
		cs.Step = 4 + 3
		cs.Signature.init()
		if _,_,_, err := cs.sign(priKey); err != nil {
			fmt.Println("333",err)
			return
		}
		bba := newBinaryByzantineAgreement()

		bba.Credential = cs
		bba.B = 0
		bba.Hash = hash
		//b
		bba.EsigB.round = bba.Credential.Round
		bba.EsigB.step = bba.Credential.Step
		bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
		bba.EsigB.Signature.init()
		bba.EsigB.sign(priKey)

		//hash
		bba.EsigV.round = bba.Credential.Round
		bba.EsigV.step = bba.Credential.Step
		bba.EsigV.val = hash.Bytes()
		bba.EsigV.Signature.init()
		bba.EsigV.sign(priKey)

		msgBba := NewMsgBinaryByzantineAgreement(bba)
		msgBba.Send()
	}

	select {
		case <-an.actualNode.StopCh():
	}
}

func TestCs_validate_success(t *testing.T){
	Config().prVerifier = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _,_,_, err := cs.sign(priKey); err != nil {
		fmt.Println("111",err)
		return
	}
	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	time.Sleep(2 * time.Second)
}

//credential has no right to verify
func TestCs_validate_fail_1(t *testing.T){
	Config().prVerifier = 1
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _,_,_, err := cs.sign(priKey); err != nil {
		fmt.Println("111",err)
		return
	}
	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	time.Sleep(2 * time.Second)
}

//verify CredentialSig fail: invalid chain id for signer
func TestCs_validate_fail_2(t *testing.T){
	Config().prVerifier = 10000000000
	priKey := generatePrivateKey()
	cs := &CredentialSign{}
	cs.Round = 100
	cs.Step = 2
	cs.Signature.init()
	if _,_,_, err := cs.sign(priKey); err != nil {
		fmt.Println("111",err)
		return
	}
	cs.V.IntVal.Add(&cs.V.IntVal, common.Big2)
	msgcs := NewMsgCredential(cs)
	msgcs.Send()
	time.Sleep(2 * time.Second)
}