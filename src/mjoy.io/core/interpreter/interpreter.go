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
// @File: interpreter.go
// @Date: 2018/05/19 10:50:19
////////////////////////////////////////////////////////////////////////////////

package interpreter

import "mjoy.io/common/types"

type Database interface {

}

type Executor interface {
	Exec() []byte
}

type Accesser interface {
	Format() string
}

type CodeContext interface {
	GetFrom() *types.Address
	GetTo() *types.Address
	GetProducer() *types.Address
	GetCode() *types.Address
	IsAccount(address *types.Address) bool
}
