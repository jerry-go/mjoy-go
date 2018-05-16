package txpool

import (
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
)

/*

*/
type interpreter interface {
	GetPriorityFromTransaction(tx *transaction.Transaction)*types.BigInt
}
