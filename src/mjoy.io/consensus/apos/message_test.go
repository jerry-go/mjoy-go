package apos

import (
	"math/big"
	"testing"
	"fmt"
	"mjoy.io/common/types"
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
	an := newAllNodeManager()
	verifierCnt := an.initTestCommonNew(1)
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
