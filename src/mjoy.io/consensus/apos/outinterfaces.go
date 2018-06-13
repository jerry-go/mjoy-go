package apos

import (
	"mjoy.io/common/types"
	"math/big"
)

/*
For out caller
*/

type dataPack interface {

}

type OutMsger interface {
	SendMsg([]byte)error
	BroadCast([]byte)error
	GetMsg() <-chan []byte

	GetDataMsg() <-chan dataPack

	SendCredential(*CredentialSig) error
}
//some out tools offered by Mjoy,such as signer and blockInfo getter
type CommonTools interface {
	//
	SIG(r,s int ,Qr_k types.Hash )(R,S,V *big.Int)

	GetQr_k(k int)types.Hash
	GetNowBlockNum()int
}


