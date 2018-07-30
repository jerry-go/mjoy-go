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
// @File: assert.go
// @Date: 2018/07/30 19:25:30
////////////////////////////////////////////////////////////////////////////////

package interpreter

func assert(test bool) {
	if !test {
		panic("")
	}
}

func assertEx(test bool, msg string) {
	if !test {
		panic(msg)
	}
}

// assert msg
const (
	NotSupport = "not support"
)