package transaction

import (
	"testing"
	"mjoy.io/common/types"

	"math/big"
	"reflect"
)

func TestTransactionNew(t *testing.T){
	want := "0xaa7f4fb00b7ad45824a36645ddb5b317265a1dde9901160594ac54c68a315625"

	var nonce uint64 =  10
	data := []byte{}
	data = append(data, 1, 4 ,5)
	address := types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	actions := []Action{{
		Address:&address,
		Params:data,
	},}

	tx := newTransaction(nonce , &address , actions)

	sig := NewMSigner(big.NewInt(1))
	h := sig.Hash(tx)
	t.Logf("transaction hash = %x", h)
	if !reflect.DeepEqual(h, types.HexToHash(want)) {
		t.Errorf("Error: have hash: %x, want: %v", h, want)
	}
}