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
	"math"
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
10
 */
func TestGetBinomialSortitionPriorityByHash(t *testing.T) {
	bd := new(binomialDistribution)
	hash := types.Hash{}
	hash[0] = 70
	ret := bd.getSortitionPriorityByHash(hash, 10, 10, 100)
	fmt.Println(ret)

	hash[0] = 128
	ret = bd.getSortitionPriorityByHash(hash, 10, 10, 100)
	fmt.Println(ret)

	hash[0] = 200
	ret = bd.getSortitionPriorityByHash(hash, 10, 10, 100)
	fmt.Println(ret)

	ret = bd.getSortitionPriorityByHash(TimeOut, 10, 10, 100)
	fmt.Println(ret)
}


/*
0.14592027257189427
0.5
0.8540797274281058
0.9824925094901688
0.9992172988709987
0.9999875866968306
0.9999999319598855
0.9999999998730186
0.9999999999999201
1
1
*/
func TestGetSumGaussian(t *testing.T) {
	w := 10
	p := 10.0/100.0
	e := float64(w) * p
	sigma := e * (1 - p)
	sigma = math.Sqrt(sigma)

	fmt.Println(p,e,sigma)

	for i := 0; i <= w; i++ {
		fmt.Println(i, normalCdf(e, sigma, float64(i)))
	}

	fmt.Println(normalInverseCdf(e, sigma, 0.8540797274281058))
	fmt.Println(normalInverseCdf(e, sigma, 0.9999999999999201))
	fmt.Println(normalInverseCdf(e, sigma, 0.14592027257189427))
	fmt.Println(normalInverseCdf(e, sigma, 0))
}

/*
0
1
2
9
*/
func TestGetGaussianSortitionPriorityByHash(t *testing.T) {
	gd := new(gaussianDistribution)
	hash := types.Hash{}
	hash[0] = 20
	ret := gd.getSortitionPriorityByHash(hash, 10, 10, 100)
	ret1 := gd.getSortitionPriorityByHashFloat(hash, 10, 10, 100)
	fmt.Println(ret,ret1)

	hash[0] = 127
	ret = gd.getSortitionPriorityByHash(hash, 10, 10, 100)
	ret1 = gd.getSortitionPriorityByHashFloat(hash, 10, 10, 100)
	fmt.Println(ret,ret1)

	hash[0] = 200
	ret = gd.getSortitionPriorityByHash(hash, 10, 10, 100)
	ret1 = gd.getSortitionPriorityByHashFloat(hash, 10, 10, 100)
	fmt.Println(ret,ret1)

	ret = gd.getSortitionPriorityByHash(TimeOut, 10, 10, 100)
	ret1 = gd.getSortitionPriorityByHashFloat(TimeOut, 10, 10, 100)
	fmt.Println(ret,ret1)
}

func TestPerformance(t *testing.T) {
	last := new(big.Float)
	logger.Info("no optimazation, start.time")
	w := 500
	for i := 0; i <= w; i++ {
		last = getSumBinomial(int64(w), 1000, 10000, int64(i))
		//fmt.Println(last)
	}
	logger.Info("no, optimazation end.time")

	logger.Info(" start.time")
	for i := 0; i <= w; i++ {
		last = getSumBinomialBasedLastSum(int64(w), 1000, 10000, int64(i), last)
		//fmt.Println(last)
	}
	logger.Info(" end.time")
}

func TestPerformance1(t *testing.T) {
	last := new(big.Float)
	w := 5000

	logger.Info(" start.time")
	for i := 0; i <= w; i++ {
		last = getSumBinomialBasedLastSum(int64(w), 2000, 10000, int64(i), last)
		//fmt.Println(last)
	}
	logger.Info(" end.time")
}

func TestPerformance1Gaussian(t *testing.T) {
	w := 5000
	p := 2000.0/10000.0
	e := float64(w) * p
	sigma := math.Sqrt(e * (1 - p))

	logger.Info(" start.time")
	for i := 0; i <= w; i++ {
		normalCdf(e, sigma, float64(i))
		//fmt.Println(last)
	}
	logger.Info(" end.time")
}



func TestGetBinomiaGaussianDiff(t *testing.T) {
	last := new(big.Float)
	w := 1000
	p := 2000.0/5000.0
	e := float64(w) * p
	sigma := e * (1 - p)

	fmt.Println(p,e,sigma)

	for i := 0; i <= w; i++ {
		last = getSumBinomialBasedLastSum(int64(w), 2000, 5000, int64(i), last)
		fmt.Println(last, normalCdf(e, math.Sqrt(sigma), float64(i)))
	}

}