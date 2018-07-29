package actioncontext

import (
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
)

type Context struct {
	action *transaction.Action
	sender types.Address
	con contract
	//db TODO:
}

type contract struct {
	iid types.Hash				// interpreter id
	creator types.Address		// creator address
	self types.Address			// contract address
	code []byte					// code
}


