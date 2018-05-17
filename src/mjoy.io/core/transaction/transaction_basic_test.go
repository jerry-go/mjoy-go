package transaction

import (
	"testing"
	"mjoy.io/common/types"

	"math/big"
	"fmt"
	"mjoy.io/utils/crypto"
)
var mSigner = NewMSigner(big.NewInt(1))


func TestTransactionNew(t *testing.T){
	var nonce uint64 =  10
	data := []byte{}
	data = append(data , 1 , 4 ,5)
	address := types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	actions := []Action{{
		Address:&address,
		Params:data,
	},}

	tx := newTransaction(nonce , &address , actions)
	_ = tx
}


func TestAsMessageGenerate(t *testing.T){
	var nonce uint64 =  10
	data := []byte{}
	data = append(data , 1 , 4 ,5)
	address := types.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	actions := []Action{{
		Address:&address,
		Params:data,
	},}
	//new transaction
	tx := newTransaction(nonce , &address , actions)
	//create key
	key , _ := crypto.GenerateKey()
	//Sign tx
	txSigned,_ := SignTx(tx,mSigner,key)


	msg , err :=txSigned.AsMessage(mSigner)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("msg:" , msg)
}

