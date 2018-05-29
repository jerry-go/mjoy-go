package sdk

import (
	"mjoy.io/common/types"
	"errors"
)

func Sys_GetValue(handlePtr *TmpStatusManager , contractAddress types.Address , key []byte)[]byte{
	//nil check
	if nil == handlePtr {
		return nil
	}
	return handlePtr.GetValue(contractAddress , key)
}

func Sys_SetValue(handlePtr *TmpStatusManager ,contractAddress types.Address , key []byte , value []byte)error{
	//nil check
	if nil == handlePtr {
		return errors.New("ptr")
	}

	return handlePtr.SetValue(contractAddress , key , value)
}

