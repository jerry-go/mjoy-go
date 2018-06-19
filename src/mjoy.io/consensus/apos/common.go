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
)

// Determine a potential verifier(leader) by hash
func isPotVerifier(hash []byte, leader bool) bool {
	h := big.NewInt(0).SetBytes(hash)
	prVal := big.NewInt(0)
	if leader {
		prVal.SetUint64(Config().prLeader)
	} else {
		prVal.SetUint64(Config().prVerifier)
	}

	return h.Div(h, Config().precision()).Cmp(prVal) < 0
}

func isHonest(vote, all int) bool {
	v := big.NewInt(int64(vote))
	a := big.NewInt(int64(all))
	pH := big.NewInt(0).SetUint64(Config().prH)
	return v.Div(v.Mul(v, honestPercision), a).Cmp(pH) >= 0
}

func isAbsHonest(vote int, leader bool) bool {
	a := Config().maxPotVerifiers
	if leader {
		a = Config().maxPotLeaders
	}
	v := big.NewInt(int64(vote))
	pH := big.NewInt(0).SetUint64(Config().prH)
	return v.Div(v.Mul(v, honestPercision), a).Cmp(pH) >= 0
}
