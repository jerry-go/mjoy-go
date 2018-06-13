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
// @File: types.go
// @Date: 2018/06/12 11:01:51
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"sync"
	"bytes"
	"github.com/tinylib/msgp/msgp"
	"math/big"
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


//some system param(algorand system param) for step goroutine.
//goroutine can set param by SetXXXX,and get param by GetXXXX
type algoParam struct {

	lock                sync.RWMutex
	k                   int
	pLeader             float64
	leaderDifficulty    *big.Int
	pVerifier           float64
	verifierDifficulty  *big.Int
	maxSteps            int
	m                   int
	nNodes              int
	timeDelayA int  //time A
	timeDelayY int  //time λ



}
func (this *algoParam)SetDefault(){
	this.lock.Lock()
	defer this.lock.Unlock()
	this.k = 1
	this.pLeader = 0.1
	this.leaderDifficulty = big.NewInt(10)
	this.pVerifier = 0.2
	this.leaderDifficulty = big.NewInt(5)
	this.maxSteps = 183
	this.m = this.maxSteps - 3
	this.nNodes = 100

}

//set param k
func (this *algoParam)SetK(k int){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.k = k
}
//get param k
func (this *algoParam)GetK()int{
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.k
}

//Set pLeader
func (this *algoParam)SetPleader(pLeader float64){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.pLeader = pLeader
}
//Get pLeader
func (this *algoParam)GetPleader()float64{
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.pLeader
}
//Set pVerifier
func (this *algoParam)SetPverifier(pVerifier float64){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.pVerifier = pVerifier
}

//Get pVerifier
func (this *algoParam)GetPverifier()float64{
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.pVerifier
}

//Set maxSteps
func (this *algoParam)SetMaxSteps(ms int){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.maxSteps = ms
}
//Get maxSteps
func (this *algoParam)GetMaxSteps()int{
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.maxSteps
}

//Set nNodes
func (this *algoParam)SetNnodes(n int){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.nNodes = n
}

//Get nNodes
func (this *algoParam)GetNnodes()int{
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.nNodes
}

func newAlgoParam()*algoParam{
	n := new(algoParam)
	n.SetDefault()
	return n
}



