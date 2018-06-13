package apos

import (
	"mjoy.io/common/types"
	"math/big"
)

/*
For out caller
*/

type OutMsger interface {
	SendMsg([]byte)error
	BroadCast([]byte)error
	GetMsg()<-chan []byte
}
//some out tools offered by Mjoy,such as signer and blockInfo getter
type CommonTools interface {
	//
	SIG([]byte )(R,S,V *big.Int)

	GetQr_k(k int)types.Hash
	GetNowBlockNum()int
}


