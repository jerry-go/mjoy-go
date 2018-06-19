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
// @File: validate.go
// @Date: 2018/06/15 14:38:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"math/big"
	"mjoy.io/common/types"
	"errors"
	"mjoy.io/utils/crypto"
)

type MsgValidator struct {
	apos      *Apos
}

func NewMsgValidator(apos *Apos) *MsgValidator {
	validator := &MsgValidator{
		apos:            apos,
	}
	return validator
}

func (v *MsgValidator)ValidateCredential(cs *CredentialSig) error{

	srcBytes := []byte{}
	srcBytes = append(srcBytes , cs.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , cs.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , cs.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)

	leader := false
	if 1 == cs.Step.IntVal.Uint64() {
		leader = true
	}


	if isPotVerifier(h, leader) == false {
		return errors.New("credential has no right to verify")
	}

	cd := CredentialData{cs.Round,cs.Step, v.apos.commonTools.GetQr_k(1)}
	sig := &SignatureVal{&cs.R, &cs.S, &cs.V}
	//verify signature
	err := v.apos.commonTools.SigVerify(cd.Hash(), sig)
	if err != nil {
		logger.Info("verify CredentialSig fail", err)
		return err
	}
	return nil
}

func (v *MsgValidator)ValidateM1(msg *M1) error{
	//verify Credential
	err := v.ValidateCredential(msg.Credential)
	if err != nil {
		logger.Info("ValidateM1 fail", err)
		return err
	}

	//verify esig
	err = v.apos.commonTools.ESigVerify(msg.Block.Hash(), msg.Esig)
	if err != nil {
		logger.Info("verify M1 ephemeral signature fail", err)
		return err
	}

	//verify block
	//todo need verify header and body


	return nil
}

func (v *MsgValidator)ValidateM23(msg *M23) error{
	//verify Credential
	err := v.ValidateCredential(msg.Credential)
	if err != nil {
		logger.Info("ValidateM23 fail", err)
		return err
	}

	//verify esig
	err = v.apos.commonTools.ESigVerify(msg.Hash, msg.Esig)
	if err != nil {
		logger.Info("verify M23 ephemeral signature fail", err)
		return err
	}

	return nil
}

func (v *MsgValidator)ValidateMCommon(msg *MCommon) error{
	//verify Credential
	err := v.ValidateCredential(msg.Credential)
	if err != nil {
		logger.Info("Validate common message fail", err)
		return err
	}

	//verify esig for b
	bBash := types.BytesToHash(big.NewInt(int64(msg.B)).Bytes())
	err = v.apos.commonTools.ESigVerify(bBash, msg.EsigB)
	if err != nil {
		logger.Info("verify M23 ephemeral signature fail", err)
		return err
	}

	//verify esig for v
	err = v.apos.commonTools.ESigVerify(msg.Hash, msg.EsigV)
	if err != nil {
		logger.Info("verify M23 ephemeral signature fail", err)
		return err
	}

	return nil
}