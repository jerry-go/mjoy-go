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
// @File: apos_signing.go
// @Date: 2018/06/13 11:12:13
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"math/big"
	"fmt"
	"errors"
	"mjoy.io/common/types"
	"crypto/ecdsa"
	"mjoy.io/common"
	"mjoy.io/utils/crypto"
)

//go:generate msgp

var (
	ErrInvalidSig = errors.New("invalid  v, r, s values")
	ErrInvalidChainId = errors.New("invalid chain id for signer")
)

// Signer encapsulates apos signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type signer interface {
	// sign the obj
	sign(prv *ecdsa.PrivateKey) (R *big.Int, S *big.Int, V *big.Int, err error)

	// Sender returns the sender address of the Credential.
	sender() (types.Address, error)

	// hash
	hash() types.Hash
}

// signature R, S, V
type signature struct {
	R *big.Int
	S *big.Int
	V *big.Int
}

type signValue interface {
	// check the signature obj is initialized, if not, throw painc
	checkObj()

	// get() computes R, S, V values corresponding to the
	// given signature.
	get(sig []byte) (err error)
}

func (s *signature) checkObj() (err error) {
	if s.R == nil || s.S == nil || s.V == nil {
		panic(fmt.Errorf("signature obj is not initialized"))
	}
}

func (s *signature) get(sig []byte) (err error) {
	s.checkObj()

	if len(sig) != 65 {
		return errors.New(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	} else {
		s.R = new(big.Int).SetBytes(sig[:32])
		s.S = new(big.Int).SetBytes(sig[32:64])

		if Config().chainId != nil && Config().chainId.Sign() != 0 {
			s.V = big.NewInt(int64(sig[64] + 35))
			s.V.Add(s.V, Config().chainIdMul)
		} else {
			s.V = new(big.Int).SetBytes([]byte{sig[64] + 27})
		}
	}
	return nil
}

// long-term key singer
type Credential struct {
	Round  		uint64		// round
	Step   		uint64		// step
	Quantity    []byte		// quantity(seed, Qr-1)

	signature
}

func (cret *Credential) sign(prv *ecdsa.PrivateKey) (R *big.Int, S *big.Int, V *big.Int, err error) {
	if prv == nil {
		err := errors.New(fmt.Sprintf("private key is empty"))
		return nil, nil, nil, err
	}

	hash := cret.hash()
	if (hash == types.Hash{}) {
		err := errors.New(fmt.Sprintf("the hash of credential is empty"))
		return nil, nil, nil, err
	}

	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, nil, nil, err
	}

	err = cret.signature.get(sig)
	if err != nil {
		return nil, nil, nil, err
	}
	R = cret.signature.R
	S = cret.signature.S
	V = cret.signature.V

	return R, S, V, nil
}

func (cret *Credential) sender() (types.Address, error) {
	cret.signature.checkObj()

	if Config().chainId != nil && deriveChainId(cret.signature.V).Cmp(Config().chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}

	V := &big.Int{}
	if Config().chainId.Sign() != 0 {
		V = V.Sub(cret.signature.V, Config().chainIdMul)
		V.Sub(V, common.Big35)
	} else{
		V = V.Sub(cret.signature.V, common.Big27)
	}
	address, err :=  recoverPlain(cret.hash(), cret.signature.R, cret.signature.S, V, true)
	return address, err
}

func (cret *Credential) hash() types.Hash {
	hash, err := common.MsgpHash(cret)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

func recoverPlain(sighash types.Hash, R, S, Vb *big.Int, homestead bool) (types.Address, error) {
	if Vb.BitLen() > 8 {
		return types.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64())
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return types.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return types.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return types.Address{}, errors.New("invalid public key")
	}
	var addr types.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
