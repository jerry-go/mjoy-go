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
// @File: message_test.go
// @Date: 2018/06/13 09:12:13
////////////////////////////////////////////////////////////////////////////////

package message

import (
	"fmt"
	"time"
)

type Testmsg struct {
	a int
	stop bool
	*msgPriv
}

// new a message
func NewTestmsg() *Testmsg {
	tm := Testmsg{
		a: 1,
		stop: false,
		msgPriv: newMsgPriv(),
	}
	Msgcore().handle(&tm)

	return &tm
}

func (tm *Testmsg) dataHandle(data interface{}) {
	go func() {
		for !tm.stop {
			fmt.Printf("handle: %v %v\n", tm.a, tm.stop)
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("handle stop\n")
	}()
}

func (tm *Testmsg) stopHandle() {
	fmt.Printf("stop ...\n")
	tm.stop = true
}

func ExampleMsg() {
	// new
	tm := NewTestmsg()
	// fix and process
	tm.a = 123
	// send
	tm.Send()
	time.Sleep(5 * time.Second)
	// close
	tm.Close()
	time.Sleep(1 * time.Second)
	fmt.Printf("End\n")

	// Output:
	// handle: 123 false
	// handle: 123 false
	// handle: 123 false
	// handle: 123 false
	// handle: 123 false
	// stop ...
	// handle stop
	// End
}

