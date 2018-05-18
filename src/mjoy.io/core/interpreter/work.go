package interpreter

import "mjoy.io/core/transaction"

type Work struct {
	transaction *transaction.Transaction
	errCh <-chan error
}

func NewWork(tx *transaction.Transaction)*Work{
	w := new(Work)
	w.transaction = tx
	w.errCh = make(chan error , 1)
	return w
}








