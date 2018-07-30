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
// @File: sortition_test.go
// @Date: 2018/07/27 14:11:27
//
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"testing"
	"fmt"
	"mjoy.io/common/types"
	"math/big"
)

func TestGetExpK(t *testing.T) {
	aa := getPexpK(10, 1, 10)
	fmt.Println(aa)

	bb := getPexpK(2, 9, 10)
	fmt.Println(bb)
}

/*
0.3486784401
0.387420489
0.1937102445
0.057395627999999999998
0.0111602609999999999995
0.0014880348000000000001
0.00013778100000000000001
8.748e-06
3.645e-07
9e-09
*/
func TestGetBinomial(t *testing.T) {
	for i := 0; i <= 10; i++ {
		fmt.Println(getBinomial(int64(i), 10, 10, 100))
	}
}


/*
0.3486784401
0.7360989291
0.9298091736
0.9872048016
0.9983650626
0.9998530974
0.9999908784
0.9999996264
0.9999999909
0.99999999989999999997
0.99999999999999999995
*/
func TestGetSumBinomial(t *testing.T) {
	for i := 0; i <= 10; i++ {
		fmt.Println(getSumBinomial(10, 10, 100, int64(i)))
	}
}

func TestGetSumBinomialBasedLastSum(t *testing.T) {
	last := new(big.Float)
	for i := 0; i <= 10; i++ {
		last = getSumBinomialBasedLastSum(10, 10, 100, int64(i), last)
		fmt.Println(last)
	}
}

/**
0
1
2
10/
 */
func TestGetSortitionPriorityByHash(t *testing.T) {
	hash := types.Hash{}
	hash[0] = 70
	ret := getSortitionPriorityByHash(hash, 10, 10, 100)
	fmt.Println(ret)

	hash[0] = 128
	ret = getSortitionPriorityByHash(hash, 10, 10, 100)
	fmt.Println(ret)

	hash[0] = 200
	ret = getSortitionPriorityByHash(hash, 10, 10, 100)
	fmt.Println(ret)

	ret = getSortitionPriorityByHash(TimeOut, 10, 10, 100)
	fmt.Println(ret)
}