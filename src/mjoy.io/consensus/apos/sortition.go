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
// @File: sortition.go
// @Date: 2018/07/27 14:11:27
//
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"mjoy.io/common/types"
	"math/big"
	"math"
	"github.com/ematvey/gostat"
)

type SortitionPriority interface {
	getSortitionPriorityByHash(hash types.Hash, w, tao, W int64) (j int64)
}

func sortition(tools CommonTools, tao, round, step ,w ,W uint64) (types.Hash, []byte, int) {
	return types.Hash{}, nil, 0
}

type gaussianDistribution struct {

}

func normalCdf(μ, σ , x float64) float64 {
	return (1.0 / 2.0) * (1 + math.Erf((x-μ)/(σ*math.Sqrt2)))
}

func normalInverseCdf(μ, σ , p float64) float64 {
	return σ * (stat.Z_InvCDF_For(p))  + μ
}

func (gs *gaussianDistribution) getSortitionPriorityByHash(hash types.Hash, w, tao, W int64) (j int64)  {
	p := float64(tao)/float64(W)
	e := float64(w) * p
	sigma := math.Sqrt(e * (1 - p))

	hashBig := new(big.Int).SetBytes(hash.Bytes())
	hashP := new(big.Float).Quo(new(big.Float).SetInt(hashBig), new(big.Float).SetInt(maxUint256))

	for j = 0; j < w; j++{
		if hashP.Cmp(big.NewFloat(normalCdf(e, sigma, float64(j)))) < 0 {
			break
		}
	}
	return j
}

func (gs *gaussianDistribution) getSortitionPriorityByHashFloat(hash types.Hash, w, tao, W int64) (j float64)  {
	p := float64(tao)/float64(W)
	e := float64(w) * p
	sigma := math.Sqrt(e * (1 - p))

	hashBig := new(big.Int).SetBytes(hash.Bytes())
	hashP := new(big.Float).Quo(new(big.Float).SetInt(hashBig), new(big.Float).SetInt(maxUint256))
	hashPf,_ := hashP.Float64()

	j = normalInverseCdf(e, sigma, hashPf)
	if j > float64(w) {
		j = float64(w)
	}
	return j
}

type binomialDistribution struct {

}

func (bs *binomialDistribution) getSortitionPriorityByHash(hash types.Hash, w, tao, W int64) (j int64)  {
	hashBig := new(big.Int).SetBytes(hash.Bytes())
	hashP := new(big.Float).Quo(new(big.Float).SetInt(hashBig), new(big.Float).SetInt(maxUint256))

	last := new(big.Float)

	for j = 0; j < w; j++{
		last = getSumBinomialBasedLastSum(w, tao, W, j, last)
		if hashP.Cmp(last) < 0 {
			break
		}
	}
	return j
}

func getSumBinomialBasedLastSum(w, tao, W, j int64, last *big.Float)  *big.Float {
	ret := new(big.Float)
	ret.Add(last, getBinomial(j, w, tao, W))
	return ret
}

func getSumBinomial(w, tao, W, j int64)  *big.Float {
	ret := new(big.Float)
	i := j
	for i=0; i <= j; i++{
		ret.Add(ret, getBinomial(i, w, tao, W))
	}
	return ret
}

//k < w
// Binomial(w,k) *(p**k) * ((1−p)**(w−k))
// p =tao/W
func getBinomial(k, w, tao, W int64)  *big.Float {
	binomial := new(big.Float).SetInt(new(big.Int).Binomial(w, k))
	pRet := new(big.Float).Mul(getPexpK(k, tao, W), getPexpK(w - k, W-tao, W))
	ret := binomial.Mul(binomial, pRet)
	return ret
}

func getPexpK(k, tao, W int64) *big.Float{
	taoBig := big.NewInt(tao)
	kBig := big.NewInt(k)
	Wbig := big.NewInt(W)
	numerator := new(big.Int).Exp(taoBig, kBig, nil)
	denominator := new(big.Int).Exp(Wbig, kBig, nil)

	//p**K
	ret := new(big.Float).Quo(new(big.Float).SetInt(numerator), new(big.Float).SetInt(denominator))
	return ret
}