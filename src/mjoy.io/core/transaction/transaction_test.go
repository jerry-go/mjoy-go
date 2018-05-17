package transaction

import (
	"testing"

	"math/big"
	"reflect"
	"mjoy.io/utils/crypto"
)

var (
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress  = crypto.PubkeyToAddress(testKey.PublicKey)
)

func TestTransactionNew(t *testing.T){
	var nonce uint64 =  10
	data := []byte{}
	data = append(data, 1, 4 ,5)
	actions := []Action{{
		Address: &testAddress,
		Params:data,
	},}

	tx := newTransaction(nonce , &testAddress , actions)

	sig := NewMSigner(big.NewInt(1))

	//h := sig.Hash(tx)
	//t.Logf("transaction hash = %x", h)
	//if !reflect.DeepEqual(h, testAddress) {
	//	t.Errorf("Error: have hash: %x, want: %v", h, testAddress.Hex())
	//}

	txWithSig, err := SignTx(tx, sig, testKey)
	if err != nil {
		t.Errorf("SignTx error: %v", err)
	}
	txWithSig.PrintVSR()

	addr, err :=sig.Sender(txWithSig)
	if err != nil {
		t.Errorf("Sender error: %v", err)
	}
	if !reflect.DeepEqual(addr, testAddress) {
		t.Errorf("Error: get addr: %v want addr: %v", addr, testAddress)
	}
}
