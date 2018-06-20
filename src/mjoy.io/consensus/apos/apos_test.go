package apos

import (
	"testing"
	"time"
	"fmt"
	"mjoy.io/common/types"
	"math/big"
)

func TestAposRunning(t *testing.T){
	an := newAllNodeManager()
	an.init()
	for{
		time.Sleep(3*time.Second)
		//fmt.Println("apos_test doing....")
	}
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
