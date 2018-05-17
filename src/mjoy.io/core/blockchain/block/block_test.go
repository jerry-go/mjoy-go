package block

import (
	"testing"
	"mjoy.io/common/types"
	"math/big"
	"mjoy.io/utils/crypto"
	"fmt"
	"bytes"
)

func TestHeaderSignature(t *testing.T) {
	header := &Header{Number:types.NewBigInt(*big.NewInt(334))}
	chainId := big.NewInt(101)
	singner := NewBlockSigner(chainId)

	var(
		key , _ = crypto.GenerateKey()
		address = crypto.PubkeyToAddress(key.PublicKey)
	)
	fmt.Println(address.Hex())
	signHeaer, err := SignHeader(header, singner, key)
	if err != nil {
		t.Fatalf("SignHeader fail")
	}

	fmt.Println(signHeaer)

	getaddress,err := singner.Sender(signHeaer)
	if err != nil {
		t.Fatalf("cann't get senser form header %v",err)
	}
	fmt.Println(getaddress.Hex())

	if !bytes.Equal(getaddress.Bytes(),address.Bytes())  {
		t.Fatalf("address is not same got:%v, want:%v",getaddress.Hex(), address.Hex())
	}
}