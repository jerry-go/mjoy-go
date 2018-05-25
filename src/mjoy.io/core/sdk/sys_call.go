package sdk

import (
	"mjoy.io/common/types"
	"errors"
)

func Sys_GetValue(contractAddress types.Address , key []byte)[]byte{
	//nil check
	if nil == PtrSdkManager.pStatusManager {
		return nil
	}
	return PtrSdkManager.pStatusManager.GetValue(contractAddress , key)
}

func Sys_SetValue(contractAddress types.Address , key []byte , value []byte)error{
	//nil check
	if nil == PtrSdkManager.pStatusManager {
		return errors.New("ptr")
	}

	return PtrSdkManager.pStatusManager.SetValue(contractAddress , key , value)
}

