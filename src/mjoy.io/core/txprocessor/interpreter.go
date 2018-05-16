package txprocessor

import (
	"mjoy.io/core/transaction"
	"math/big"
)

type Interpreter interface {
	/*
	if a transaction want join in txpool , a priority we should have,
	what is the priority?it should determined by miner's config file.
	*/
	GetPriorityFromTransaction(tx *transaction.Transaction)*big.Int
	/*
	Why Interpreter should check a transaction is legal,because we do not know
	what transaction meaning for.Our transaction should just a vm running data,what's
	the mean of transaction should be determined by interpreter and maker.
	*/
	CheckTxLegal(tx *transaction.Transaction)error

	/*
	GetBasicInfo return a interpreter unmarshal info,ex. for adding in txpool
	*/
	GetBasicInfo(tx *transaction.Transaction)interface{}

}
