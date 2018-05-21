package interpreter

import (
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
)

type WorkResult struct {
	Err error
	Logs []*transaction.Log
}


type Work struct {
	actions []transaction.Action
	from types.Address
	resultChan chan WorkResult
}

func NewWork(from types.Address , actions []transaction.Action)*Work{
	w := new(Work)

	//copy actions
	w.actions= make([]transaction.Action , len(actions))
	copy(w.actions , actions)
	w.resultChan = make(chan WorkResult , 1)
	return w
}








