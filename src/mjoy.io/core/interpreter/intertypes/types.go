package intertypes

import "mjoy.io/core/sdk"

type ActionResult struct {
	Key []byte
	Val []byte
}


//SystemParams contain all system running params
type SystemParams struct {
	SdkHandler *sdk.TmpStatusManager    //contain current
}

func MakeSystemParams(sdkHandler *sdk.TmpStatusManager)*SystemParams{
	s := new(SystemParams)
	s.SdkHandler = sdkHandler
	return s
}
