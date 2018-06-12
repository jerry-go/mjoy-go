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

package algorand

import (
	"mjoy.io/common/types"
	"mjoy.io/common"
)

// for algorand1, fill block header ConsensusData filed
// id = "algorand1"
// para = Q(r) = hash(Sig(Q(r-1)), r) where r = block number
var (
	ConsensusDataId = "algorand1"
)

// credentialData is the data for generating credential
// credential = sig(credentialData)
type CredentialData struct {
	Round         types.BigInt
	Step          types.BigInt
	Quantity      types.Hash    //the seed of round r
}

type CredentialSig struct {
	Round         types.BigInt
	Step          types.BigInt
	R             types.BigInt
	S             types.BigInt
	V             types.BigInt
}

type SignatureVal struct {
	R             types.BigInt
	S             types.BigInt
	V             types.BigInt
}

func (h *CredentialData) Hash() types.Hash {
	hash, err := common.MsgpHash(h)
	if err != nil {
		return types.Hash{}
	}
	return hash
}


