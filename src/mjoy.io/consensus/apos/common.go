////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: common.go
// @Date: 2018/06/14 14:14:14
////////////////////////////////////////////////////////////////////////////////

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

