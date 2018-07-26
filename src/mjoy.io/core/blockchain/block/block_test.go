package block

import (
	"bytes"
	"fmt"
	"math/big"
	"mjoy.io/common/types"
	"mjoy.io/utils/crypto"
	"testing"
)

func TestHeaderSignatureRamdomkey(t *testing.T) {
	header := &Header{Number: types.NewBigInt(*big.NewInt(334)), Time: types.NewBigInt(*big.NewInt(1212121))}
	chainId := big.NewInt(101)
	singner := NewBlockSigner(chainId)

	var (
		key, _  = crypto.GenerateKey()
		address = crypto.PubkeyToAddress(key.PublicKey)
	)
	signHeaer, err := SignHeader(header, singner, key)
	if err != nil {
		t.Fatalf("SignHeader fail")
	}

	getaddress, err := singner.Sender(signHeaer)
	if err != nil {
		t.Fatalf("cann't get senser form header %v", err)
	}
	fmt.Println(signHeaer)

	if !bytes.Equal(getaddress.Bytes(), address.Bytes()) {
		t.Fatalf("address is not same got:%v, want:%v", getaddress.Hex(), address.Hex())
	}
}

func TestHeaderSigantureFixkey(t *testing.T) {
	conData := make([]byte, 10)
	conData[4] = 7
	Producer := types.Address{}
	Producer[10] = 1

	header := &Header{Number: types.NewBigInt(*big.NewInt(333)), Time: types.NewBigInt(*big.NewInt(1212121)), Producer: Producer, Cdata: ConsensusData{"test", conData}, Extra:[]byte("asdfasd")}
	chainId := big.NewInt(101)
	singner := NewBlockSigner(chainId)

	var (
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
		address = crypto.PubkeyToAddress(key.PublicKey)
	)

	//fmt.Println(header)

	signHeaer, err := SignHeader(header, singner, key)
	if err != nil {
		t.Fatalf("SignHeader fail")
	}

	getaddress, err := singner.Sender(signHeaer)
	if err != nil {
		t.Fatalf("cann't get senser form header %v", err)
	}
	fmt.Println(signHeaer)

	if !bytes.Equal(getaddress.Bytes(), address.Bytes()) {
		t.Fatalf("address is not same got:%v, want:%v", getaddress.Hex(), address.Hex())
	}
}
