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
// @File: count_vote_test.go
// @Date: 2018/07/19 14:15:19
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"testing"
	"time"
	"mjoy.io/common/types"
)

func commitVote(s int , hash types.Hash) {

}

func TestCvRun(t *testing.T) {
	cv := newCountVote(commitVote, types.Hash{})
	Config().maxStep = 150
	go cv.run()
	time.Sleep(5 * time.Second)
	cv.stopCh <- 1
}