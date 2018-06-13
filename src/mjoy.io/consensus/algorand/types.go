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

package algorand

import "sync"

type ConsensusData struct{
	Id     string
	Para   []byte
}

//some system param(algorand system param) for step goroutine.
//goroutine can set param by SetXXXX,and get param by GetXXXX
type algoParam struct {
	lock sync.RWMutex
	k int

}
func (this *algoParam)SetDefault(){
	this.lock.Lock()
	defer this.lock.Unlock()
	this.k = 1
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




func newAlgoParam()*algoParam{
	n := new(algoParam)
	n.SetDefault()
	return n
}



