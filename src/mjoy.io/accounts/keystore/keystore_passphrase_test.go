////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software.
//
// @File: keystore_passphrase_test.go
// @Date: 2018/05/08 17:10:08
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
	"io/ioutil"
	"testing"

	"mjoy.io/common/types"
)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)

// Tests that a json key file can be decrypted and encrypted in multiple rounds.
func TestKeyEncryptDecrypt(t *testing.T) {
	keyjson, err := ioutil.ReadFile("testdata/very-light-scrypt.json")
	if err != nil {
		t.Fatal(err)
	}
	password := ""
	address := types.HexToAddress("45dea0fb0bba44f4fcf290bba71fd57d7117cbb8")

	// Do a few rounds of decryption and encryption
	for i := 0; i < 3; i++ {
		// Try a bad password first
		if _, err := DecryptKey(keyjson, password+"bad"); err == nil {
			t.Errorf("test %d: json key decrypted with bad password", i)
		}
		// Decrypt with the correct password
		key, err := DecryptKey(keyjson, password)
		if err != nil {
			t.Fatalf("test %d: json key failed to decrypt: %v", i, err)
		}
		if key.Address != address {
			t.Errorf("test %d: key address mismatch: have %x, want %x", i, key.Address, address)
		}
		// Recrypt with a new password and start over
		password += "new data appended"
		if keyjson, err = EncryptKey(key, password, veryLightScryptN, veryLightScryptP); err != nil {
			t.Errorf("test %d: failed to recrypt key %v", i, err)
		}
	}
}
