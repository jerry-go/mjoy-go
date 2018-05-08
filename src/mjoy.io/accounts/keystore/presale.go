////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software.
//
// @File: presale.go
// @Date: 2018/05/08 17:11:08
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  <http://www.apache.org/licenses/>
//
////////////////////////////////////////////////////////////////////////////////

package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mjoy.io/accounts"
	"mjoy.io/utils/crypto"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/pbkdf2"

)

// creates a Key and stores that in the given KeyStore by decrypting a presale key JSON
func importPreSaleKey(keyStore keyStore, keyJSON []byte, password string) (accounts.Account, *Key, error) {
	key, err := decryptPreSaleKey(keyJSON, password)
	if err != nil {
		return accounts.Account{}, nil, err
	}
	key.Id = uuid.NewRandom()
	a := accounts.Account{Address: key.Address, URL: accounts.URL{Scheme: KeyStoreScheme, Path: keyStore.JoinPath(keyFileName(key.Address))}}
	err = keyStore.StoreKey(a.URL.Path, key, password)
	return a, key, err
}

func decryptPreSaleKey(fileContent []byte, password string) (key *Key, err error) {
	preSaleKeyStruct := struct {
		EncSeed string
		MjoyAddr string
		Email   string
		BtcAddr string
	}{}
	err = json.Unmarshal(fileContent, &preSaleKeyStruct)
	if err != nil {
		return nil, err
	}
	encSeedBytes, err := hex.DecodeString(preSaleKeyStruct.EncSeed)
	if err != nil {
		return nil, errors.New("invalid hex in encSeed")
	}
	if len(encSeedBytes) < 16 {
		return nil, errors.New("invalid encSeed, too short")
	}
	iv := encSeedBytes[:16]
	cipherText := encSeedBytes[16:]

	passBytes := []byte(password)
	derivedKey := pbkdf2.Key(passBytes, passBytes, 2000, 16, sha256.New)
	plainText, err := aesCBCDecrypt(derivedKey, cipherText, iv)
	if err != nil {
		return nil, err
	}
	mjoyPriv := crypto.Keccak256(plainText)
	ecKey := crypto.ToECDSAUnsafe(mjoyPriv)

	key = &Key{
		Id:         nil,
		Address:    crypto.PubkeyToAddress(ecKey.PublicKey),
		PrivateKey: ecKey,
	}
	derivedAddr := hex.EncodeToString(key.Address.Bytes()) // needed because .Hex() gives leading "0x"
	expectedAddr := preSaleKeyStruct.MjoyAddr
	if derivedAddr != expectedAddr {
		err = fmt.Errorf("decrypted addr '%s' not equal to expected addr '%s'", derivedAddr, expectedAddr)
	}
	return key, err
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func aesCBCDecrypt(key, cipherText, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(aesBlock, iv)
	paddedPlaintext := make([]byte, len(cipherText))
	decrypter.CryptBlocks(paddedPlaintext, cipherText)
	plaintext := pkcs7Unpad(paddedPlaintext)
	if plaintext == nil {
		return nil, ErrDecrypt
	}
	return plaintext, err
}

// From https://leanpub.com/gocrypto/read#leanpub-auto-block-cipher-modes
func pkcs7Unpad(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	padding := in[len(in)-1]
	if int(padding) > len(in) || padding > aes.BlockSize {
		return nil
	} else if padding == 0 {
		return nil
	}

	for i := len(in) - 1; i > len(in)-int(padding)-1; i-- {
		if in[i] != padding {
			return nil
		}
	}
	return in[:len(in)-int(padding)]
}
