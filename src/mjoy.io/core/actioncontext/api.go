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
// @File: api.go
// @Date: 2018/07/30 16:14:30
////////////////////////////////////////////////////////////////////////////////

package actioncontext

import (
	"mjoy.io/common/types"
	"mjoy.io/utils/crypto"
	"bytes"
	"fmt"
	"runtime"
)

type api interface {
	setCtx(ctx *Context)
	assert(ret bool, msg string)
}

type baseApi struct {
	ctx *Context
}

func (self *baseApi) setCtx(ctx *Context) {
	self.ctx = ctx
}

func (self *baseApi) assert(test bool, msg string) {
	if test {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	_ = file
	_ = line
	if ok {
		f := runtime.FuncForPC(pc)
		if msg == "" {
			msg = "unknown"
		}
		assertEx(false, fmt.Sprintf("%s error : %s", f.Name(), msg))
	} else {
		assert(false)
	}
}

type authorizationApi struct {
	baseApi
}

func (self *authorizationApi) RequireAuth(address types.Address) {
	self.assert(self.ctx.sender == address, fmt.Sprintf("missing authority of %s", address))
}

func (self *authorizationApi) RequireAccount(address types.Address) bool {
	return false
}

func (self *authorizationApi) RequireContract(address types.Address) bool {
	return false
}

/*func (self *authorizationApi) pushRecipient(address types.Address) {

}*/

type systemApi struct {
	baseApi
}

type assertApi struct {
	baseApi
}

func (self *assertApi) Assert(test bool, msg string) {
	assertEx(test, msg)
}

type consoleApi struct {
	baseApi
}

type crpytoApi struct {
	baseApi
}

type producerApi struct {
	baseApi
}

func (self *producerApi) Producer() types.Address {
	return types.Address{}
}

type actionApi struct {
	baseApi
}

func (self *actionApi) Contract() types.Address {
	return self.ctx.con.self
}

type databaseApi struct {
	baseApi
}

func (self *databaseApi) Emplace(key []byte, val []byte) {
	keyHash := types.BytesToHash(crypto.Keccak256(key))
	valHash := types.BytesToHash(crypto.Keccak256(val))

	ret := self.ctx.state.GetState(self.ctx.con.self, keyHash)
	self.assert(bytes.Compare(ret.Bytes(), types.Hash{}.Bytes()) == 0, "key is exist")

	self.ctx.state.SetState(self.ctx.con.self, keyHash, valHash)
	self.ctx.state.AddPreimge(valHash, val)
}

func (self *databaseApi) Modify(key []byte, val []byte) {
	keyHash := types.BytesToHash(crypto.Keccak256(key))
	valHash := types.BytesToHash(crypto.Keccak256(val))

	ret := self.ctx.state.GetState(self.ctx.con.self, keyHash)
	self.assert(bytes.Compare(ret.Bytes(), types.Hash{}.Bytes()) != 0, "key is not exist")

	self.ctx.state.SetState(self.ctx.con.self, keyHash, valHash)
	self.ctx.state.AddPreimge(valHash, val)
}

func (self *databaseApi) Erase(key []byte) {
	// TODO: ??
	// nothing to do
	self.assert(false, NotSupport)
}

func (self *databaseApi) Find(key []byte) []byte {
	keyHash := types.BytesToHash(crypto.Keccak256(key))
	ret := self.ctx.state.GetState(self.ctx.con.self, keyHash)
	if bytes.Compare(ret.Bytes(), types.Hash{}.Bytes()) == 0 {
		return nil
	}
	// TODO: maybe use self.ctx.db.Get(preimagePrefix + ret);  Refactor statedb later, may be adjustments here.
	value, err := self.ctx.db.Get(ret.Bytes())
	if err != nil {
		return nil
	}
	return value
}

func (self *databaseApi) Get(key []byte) []byte {
	ret := self.Find(key)
	self.assert(ret != nil, "unable to find key")
	return ret
}
