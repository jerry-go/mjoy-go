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
// @File: genesis_test.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package genesis

import (
	"testing"
	"mjoy.io/common/types/util"
	"mjoy.io/params"
	"math/big"
	"reflect"
	"github.com/davecgh/go-spew/spew"
	"mjoy.io/common/types"
	"mjoy.io/utils/database"
	"mjoy.io/core/blockchain"
)

var defaultGenesisHexHash = "52c98064aaa0bc6fe8ccc14db0ded6325223fda3ab376fb8ac2d3b08d4f4303c"

func TestDefaultGenesisBlock(t *testing.T) {
	block, _ := DefaultGenesisBlock().ToBlock()

	if hexHash := util.Bytes2Hex(block.Hash().Bytes()); hexHash != defaultGenesisHexHash {
		t.Errorf("wrong mainnet genesis hash, got %v,", hexHash)
	}
}


func TestSetupGenesis(t *testing.T) {
	var (
		customghash = types.HexToHash("0x5f054f020b4d217853fbfdbbcdc4af9bd4b31fb25661c03fe6c4b32c281470ea")
		customg     = Genesis{
			Config:  &params.ChainConfig{big.NewInt(500)},
			Alloc: GenesisAlloc{
				{1}: {Balance: big.NewInt(1), Storage: map[types.Hash]types.Hash{{1}: {1}}},
			},
		}
		oldcustomg = customg

		customghash2 = types.HexToHash("0x4b94a6ffcf0b1611cc12405315dd4815126460a911b8ccfa535ff9bb3d226e85")
		customg2     = Genesis{
			Config:  &params.ChainConfig{big.NewInt(700)},
			Alloc: GenesisAlloc{
				{1}: {Balance: big.NewInt(2), Storage: map[types.Hash]types.Hash{{2}: {2}}},
			},
		}
	)
	oldcustomg.Config = &params.ChainConfig{ChainId: big.NewInt(2)}
	tests := []struct {
		name       string
		fn         func(database.IDatabase) (*params.ChainConfig, types.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   types.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				return SetupGenesisBlock(db, new(Genesis))
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.DefaultChainConfig,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   types.HexToHash(defaultGenesisHexHash),
			wantConfig: params.DefaultChainConfig,
		},
		{
			name: "test block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   types.HexToHash(defaultGenesisHexHash),
			wantConfig: params.DefaultChainConfig,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "custom block in DB, genesis == custom2",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, &customg2)
			},
			wantErr:    &GenesisMismatchError{Stored: customghash, New: customghash2},
			wantHash:   customghash2,
			wantConfig: customg2.Config,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db database.IDatabase) (*params.ChainConfig, types.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},

	}

	for _, test := range tests {
		db, _ := database.OpenMemDB()
		config, hash, err := test.fn(db)
		// Check the return values.
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
			// Check database content.
			stored := blockchain.GetBlock(db, test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}
