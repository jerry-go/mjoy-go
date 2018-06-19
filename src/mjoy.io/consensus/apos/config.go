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
// @File: types.go
// @Date: 2018/06/12 11:01:51
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"sync"
	"math/big"
	"fmt"
)

var (
	decimal = big.NewInt(10)
	honestPercision = big.NewInt(100)
)

//go:generate gencodec -type config -field-override configMarshaling -out gen_config.go

//some system param(algorand system param) for step goroutine.
type config struct {
	lookback            	int			`json:"lookback"`             // lookback val, r - k
	prPrecision				uint64		`json:"precision"`            // the precision
	prLeader				uint64		`json:"probability-leader"`   // the probability of Leaders
	prVerifier				uint64		`json:"probability-verifier"` // the probability of Verifiers
	maxBBASteps         	int			`json:"max-steps"`            // the max number of BBA steps
	maxNodesPerRound    	int			`json:"max-nodes-per-round"`  // the max number of nodes per round
	prH						uint64		`json:"probability-honest"`   // the probability of honest
	blockDelay          	int  		`json:"block-delay"`          // time A, sec
	verifyDelay         	int  		`json:"verify-delay"`         // time λ, sec

	prP						*big.Int	`json:"-"`                    // 10 ^ prPrecision
	maxPotLeaders			*big.Int	`json:"-"`                    // the max number of potential leaders
	maxPotVerifiers			*big.Int	`json:"-"`                    // the max number of potential verifiers
}

func (c *config) setDefault() {
	c.lookback = 100
	c.prPrecision = 10
	c.prLeader = 1000000000		// 0.1
	c.prVerifier = 5000000000 	// 0.5
	c.maxBBASteps = 180
	c.maxNodesPerRound = 10
	c.prH = 34
	c.blockDelay = 60
	c.verifyDelay = 10
}

// about msgcore singleton
var (
	instance	*config
	once		sync.Once
)
// get the msgcore singleton
func Config() *config {
	once.Do(func() {
		instance = &config{
		}
	})
	instance.setDefault()
	instance.Verify()
	instance.verifier()
	return instance
}

func (c *config) precision() *big.Int {
	if c.prP == nil {
		c.prP = decimal.Not(big.NewInt(0).SetUint64(c.prPrecision))
	}
	return c.prP
}

func (c *config) verifier() (uint64, uint64, uint64, *big.Int, *big.Int) {
	if c.maxPotLeaders == nil {
		c.maxPotLeaders = big.NewInt(int64(c.maxNodesPerRound))
		c.maxPotLeaders.Mul(c.maxPotLeaders, c.precision())
		c.maxPotLeaders.Div(c.maxPotLeaders, big.NewInt(0).SetUint64(c.prLeader))
	}

	if c.maxPotVerifiers == nil {
		c.maxPotVerifiers = big.NewInt(int64(c.maxNodesPerRound))
		c.maxPotVerifiers.Mul(c.maxPotVerifiers, c.precision())
		c.maxPotVerifiers.Div(c.maxPotVerifiers, big.NewInt(0).SetUint64(c.prVerifier))
	}

	return c.prPrecision, c.prLeader, c.prVerifier, c.maxPotLeaders, c.maxPotVerifiers
}

func (c *config) Verify() {
	if c.lookback <= 0 {
		panic(fmt.Errorf("lookback <= 0 \n"))
	}

	if c.maxBBASteps <= 0 {
		panic(fmt.Errorf("maxBBASteps <= 0 \n"))
	}

	if c.maxNodesPerRound <= 0 {
		panic(fmt.Errorf("maxNodesPerRound <= 0 \n"))
	}

	if c.blockDelay <= 0 {
		panic(fmt.Errorf("blockDelay <= 0 \n"))
	}

	if c.verifyDelay <= 0 {
		panic(fmt.Errorf("verifyDelay <= 0 \n"))
	}

	if c.prH == 0 {
		panic(fmt.Errorf("prH == 0 \n"))
	}
	if big.NewInt(0).SetUint64(c.prH).Cmp(honestPercision) > 0 {
		panic(fmt.Errorf("prH > 100 \n"))
	}

	if c.precision().Cmp(big.NewInt(2).Not(big.NewInt(256))) > 0 {
		panic(fmt.Errorf("PrLeader > precision \n"))
	}

	if c.precision().Cmp(big.NewInt(0).SetUint64(c.prLeader)) < 0 {
		panic(fmt.Errorf("prLeader < precision \n"))
	}

	if c.precision().Cmp(big.NewInt(0).SetUint64(c.prVerifier)) < 0 {
		panic(fmt.Errorf("prVerifier < Precision \n"))
	}
}
