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
// @File: outinterfaces.go
// @Date: 2018/06/15 10:26:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"crypto/ecdsa"
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
)

/*
For out caller
*/

type dataPack interface {
}

type OutMsger interface {
	GetDataMsg() <-chan dataPack
	GetSubDataMsg() <-chan dataPack //for test
	SendInner(dataPack) error
	Send2Apos(dataPack)
	PropagateMsg(dataPack) error
}

//some out tools offered by Mjoy,such as signer and blockInfo getter
type CommonTools interface {
	//
	Sig(pCs *CredentialSign) error
	Esig(pEphemeralSign *EphemeralSign) error
	SigHash(hash types.Hash) []byte

	SigVerify(hash types.Hash, sig *SignatureVal) error
	Sender(hash types.Hash, sig *SignatureVal) (types.Address, error)

	ESigVerify(hash types.Hash, sig []byte) error
	ESender(hash types.Hash, sig []byte) (types.Address, error)

	GetLastQrSignature() []byte
	GetQrSignature(round uint64) []byte
	GetNowBlockNum() int
	GetNextRound() int
	GetNowBlockHash() types.Hash

	SetPriKey(priKey *ecdsa.PrivateKey)
	CreateTmpPriKey(step int)
	DelTmpKey(step int)
	ClearTmpKeys()

	GetProducerNewBlock(data *block.ConsensusData) *block.Block //get a new block from block producer
	MakeEmptyBlock(data *block.ConsensusData) *block.Block
	InsertChain(chain block.Blocks) (int, error)
	GetCurrentBlock()*block.Block
	GetBlockByNum(num uint64)*block.Block

	//version 1.1
	GetBlockByHash(hash types.Hash) *block.Block
}
