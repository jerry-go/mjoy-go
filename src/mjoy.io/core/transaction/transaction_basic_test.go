package transaction

import (
	"testing"
	"mjoy.io/common/types"

)

func TestTransactionNew(t *testing.T){
	var nonce uint64 =  10
	address := types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	actions := []*Action{{
		Address:&address,
		Params:make([]byte , 0),
	},}

	tx := newTransaction(nonce , &address , actions)
	_ = tx
}

