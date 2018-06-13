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
// @File: message.go
// @Date: 2018/06/13 09:34:13
////////////////////////////////////////////////////////////////////////////////

package message

import (
	"sync"
	"mjoy.io/utils/event"
	"mjoy.io/communication/p2p/discover"
)

// the channel of consensus' message
type messageChannel struct {
	data chan interface{}
	stop chan struct{}
}

// each message must implement this interface
// external interface
type Message interface {
	Send() error	// send message
	Close()			// close message processing
}

// each message must implement this interface
// to handle message
type handler interface {
	handle(h handleable)
}

// each message must implement this interface
// to handle data and stop
type handleable interface {
	dataHandle(data interface{})
	stopHandle()
}

// private message struct
// The basic structure and interface of the message are implemented and could be inherited
type msgPriv struct {
	channel messageChannel
}

func newMsgPriv() *msgPriv {
	msg := msgPriv{
		channel: messageChannel{
			data: make(chan interface{}),
			stop: make(chan struct{}),
		},
	}
	return &msg
}

func (msg msgPriv) Send() error {
	msg.channel.data <- msg
	return nil
}

func (msg msgPriv) Close() {
	close(msg.channel.stop)
}

func (msg msgPriv) handle(h handleable) {
	for {
		select {
		case data := <-msg.channel.data:
			h.dataHandle(data)
		case <-msg.channel.stop:
			h.stopHandle()
			return
		}
	}
}

func isHandler(msg interface{}) bool {
	_, ok := msg.(handler)
	return ok
}

func getHandler(msg interface{}) handler {
	hd, ok := msg.(handler)
	if !ok {
		panic("not a handler")
	}
	return hd
}

func isHandleable(msg interface{}) bool {
	_, ok := msg.(handleable)
	return ok
}

func getHandleable(msg interface{}) handleable {
	handle, ok := msg.(handleable)
	if !ok {
		panic("not a handleable")
	}
	return handle
}

// TODO:
type msgcore struct {

}

// about msgcore singleton
var (
	instance *msgcore
	once sync.Once
)
// get the msgcore singleton
func Msgcore() *msgcore {
	once.Do(func() {
		instance = &msgcore{
		}
	})
	return instance
}

// handle msg
func (mc *msgcore) handle(msg interface{}) {
	handler := getHandler(msg)
	h := getHandleable(msg)
	go handler.handle(h)
}

type eventer struct {
	feed     *event.Feed
}