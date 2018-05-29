package interpreter

import (
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
	"mjoy.io/core/interpreter/intertypes"
)

type WorkResult struct {
	Err error
	Results []intertypes.ActionResult
}


type Work struct {
	actions []transaction.Action
	contractAddress types.Address
	sysParams *intertypes.SystemParams
	resultChan chan WorkResult
}

func NewWork(contractAddress types.Address , actions []transaction.Action  , sysParams *intertypes.SystemParams)*Work{
	w := new(Work)

	//copy actions
	w.contractAddress = contractAddress     //who deal the transaction
	w.actions= make([]transaction.Action , len(actions))
	w.sysParams = sysParams
	copy(w.actions , actions)
	w.resultChan = make(chan WorkResult , 1)
	return w
}








