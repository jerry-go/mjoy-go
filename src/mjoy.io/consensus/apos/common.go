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
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
	"mjoy.io/params"
)

var TestPotVerifier = 0

// Determine a potential verifier(leader) by hash
func isPotVerifier(hash []byte, leader bool) bool {
	if TestPotVerifier != 0 {
		return true
	}
	//todo: all return false
	h := big.NewInt(0).SetBytes(hash)
	prVal := big.NewInt(0)
	if leader {
		prVal.SetUint64(Config().prLeader)
	} else {
		prVal.SetUint64(Config().prVerifier)
	}

	return h.Cmp(big.NewInt(0).Div(prVal.Mul(prVal, maxUint256), Config().precision())) < 0
}

func isHonest(vote, all int) bool {
	v := big.NewInt(int64(vote))
	a := big.NewInt(int64(all))
	pH := big.NewInt(0).SetUint64(Config().prH)
	return v.Div(v.Mul(v, honestPercision), a).Cmp(pH) >= 0
}

func isAbsHonest(vote int, leader bool) bool {
	a := Config().maxPotVerifiers
	logger.Debug("isAbsHonest maxPotVerifiers", a, "vote", vote)
	if leader {
		a = Config().maxPotLeaders
	}
	v := big.NewInt(int64(vote))
	pH := big.NewInt(0).SetUint64(Config().prH)
	return v.Div(v.Mul(v, honestPercision), a).Cmp(pH) >= 0
}

// priority queue Item
type pqItem struct {
	value    interface{}
	priority *big.Int
}

//priority Queue
type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority.Cmp(pq[j].priority) > 0
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*pqItem)
	*pq = append(*pq, item)
}

//pop the highest priority item
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

//H(SIGℓr (Qr−1), r)
func getQuantity(sigByte []byte, round uint64) (types.Hash, error) {
	q := Quantity{}
	q.Signature.init()
	err := q.Signature.get(sigByte)
	if err != nil {
		return types.Hash{}, err
	}
	q.Round = round
	return q.Hash(), nil
}

func makeEmptyBlockConsensusData(round uint64) *block.ConsensusData {
	bcd := &block.ConsensusData{}
	bcd.Id = ConsensusDataId

	cs := CredentialSign{}
	cs.init()
	cs.Round = round
	cs.Step = 1
	cs.sign(params.RewordPrikey)

	bcd.Para = cs.toBytes()
	return bcd
}

func makeBlockConsensusData(bp *BlockProposal) *block.ConsensusData {
	bcd := &block.ConsensusData{}
	bcd.Id = ConsensusDataId
	bcd.Para = bp.Credential.Signature.toBytes()
	return bcd
}

func SenderFromBlock(header *block.Header) (types.Address, error) {
	cs := &CredentialSign{}
	cs.init()
	err := cs.Signature.get(header.ConsensusData.Para)
	if err != nil {
		return types.Address{}, err
	}
	cs.Round = header.Number.IntVal.Uint64()
	cs.Step = 1
	return cs.sender()
}
