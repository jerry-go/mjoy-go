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
// @File: aaa.go
// @Date: 2018/06/15 17:12:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"bytes"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

const(
	Type_Credential = iota
	Type_BrCredential
)

//ConsensusData:the data type for sending and receiving
type ConsensusData struct{
	Step   int
	Type   int  //0:just credential data 1:credential with other info
	Para   []byte
}

func PackConsensusData(s , t int , data []byte)[]byte{
	c := new(ConsensusData)
	c.Step = s
	c.Type = t
	c.Para = append(c.Para , data...)

	var buf bytes.Buffer
	err := msgp.Encode(&buf, c)
	if err != nil{
		return nil
	}

	return buf.Bytes()
}

func UnpackConsensusData(data []byte)*ConsensusData{
	c := new(ConsensusData)
	var buf bytes.Buffer
	buf.Write(data)

	err := msgp.Decode(&buf , c)
	if err != nil{
		logger.Errorf("UnpackConsensusData Err:%s",err.Error())
		return nil
	}
	return c
}