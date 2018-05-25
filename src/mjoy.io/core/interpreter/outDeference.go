/*
This file define a OutDeference
*/

package interpreter

import (
	"mjoy.io/core/state"
	"sync"
)

type OutDeference struct {
	mu sync.RWMutex

	mjoyState *state.StateDB
}

func NewOutDeference()*OutDeference{
	o := new(OutDeference)

	return o
}

func (this *OutDeference)RegStateDB(s *state.StateDB){
	this.mjoyState = s
}

