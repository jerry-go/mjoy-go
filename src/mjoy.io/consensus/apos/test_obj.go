package apos

//instructions:
//test_obj implements the CommonTools and outMsg

/*
type OutMsger interface {
	//SendMsg([]byte)error
	BroadCast([]byte)error
	GetMsg() <-chan []byte

	GetDataMsg() <-chan dataPack

	// send msg means that the implement must send this message to apos (loopback) as a plus step
	// Propagate msg means that the implement just send msg to p2p
	SendCredential(*CredentialSig) error
	PropagateCredential(*CredentialSig) error

	SendMsg(dataPack) error
	PropagateMsg(dataPack) error
}
//some out tools offered by Mjoy,such as signer and blockInfo getter
type CommonTools interface {
	//
	SIG([]byte )(R,S,V *big.Int)
	ESIG(hash types.Hash)([]byte)
	GetQr_k(k int)types.Hash
	GetNowBlockNum()int
	GetNextRound()int
}
*/