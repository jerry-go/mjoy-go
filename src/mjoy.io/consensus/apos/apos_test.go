package apos

import (
	"testing"
	"fmt"
	"mjoy.io/common/types"
	"math/big"
)

func TestAposRunning(t *testing.T){
	an := newAllNodeManager()
	Config().blockDelay = 2
	Config().verifyDelay = 1
	an.init(1)
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