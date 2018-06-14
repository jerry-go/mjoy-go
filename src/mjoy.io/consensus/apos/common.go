package apos

import (
	"math/big"
	"fmt"
	"strconv"
	"mjoy.io/utils/crypto"
)
var (
	// maxUint256 is a big integer representing 2^256-1
	maxUint256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
)

func BytesToFloat(b []byte)(float64,error){
	bigI := new(big.Int)
	bigI.SetBytes(b[:])

	s := fmt.Sprintf("0.%d" , bigI.Uint64())

	endFloat , err := strconv.ParseFloat(s , 64)
	if err != nil {
		endFloat = 0.0
	}
	return endFloat , err
}



func BytesToDifficulty(b []byte) (*big.Int){
	bigI := new(big.Int).SetBytes(b)
	target := new(big.Int).Div(maxUint256, bigI)
	return target
}

func GetDifficulty(pCredentialSig *CredentialSig) *big.Int {
	srcBytes := []byte{}
	srcBytes = append(srcBytes , pCredentialSig.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	return BytesToDifficulty(h)
}

func EndConditon(voteNum, target int) bool {
	if (3 * voteNum) > (2 * target) {
		return true
	} else {
		return false
	}
}

