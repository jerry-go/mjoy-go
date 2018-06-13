package apos

import (
	"math/big"
	"fmt"
	"strconv"
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
