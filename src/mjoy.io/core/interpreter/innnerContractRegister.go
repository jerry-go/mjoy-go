/*
All Innercontract implements should be added into InnerRegister slice.
*/

package interpreter

import "mjoy.io/common/types"

type innerRegisterMap struct {
	address types.Address
	inner   InnerContract
}

type InnersRegister []innerRegisterMap

var allInnerRegister InnersRegister = InnersRegister{
	{types.Address{} , nil},
	{types.Address{} , nil},
}

