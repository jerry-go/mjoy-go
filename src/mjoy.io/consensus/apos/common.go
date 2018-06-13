package apos

import (
	"math/big"
	"fmt"
	"strconv"
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
