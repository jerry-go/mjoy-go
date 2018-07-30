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
// @File: context.go
// @Date: 2018/07/30 16:16:30
////////////////////////////////////////////////////////////////////////////////

package actioncontext

import (
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
	"mjoy.io/core/state"
	"mjoy.io/utils/database"
)

type Context struct {
	action *transaction.Action
	sender types.Address
	con contract

	// state db
	state *state.StateDB
	// database
	db    database.IDatabase
}

type contract struct {
	iid types.Hash				// interpreter id
	creator types.Address		// creator address
	self types.Address			// contract address
	code []byte					// code
}

func (ctx *Context) Init() {

}

func (ctx *Context) Exec() {

}
