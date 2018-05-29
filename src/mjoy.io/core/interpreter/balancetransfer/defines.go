package balancetransfer

import (
	"encoding/json"
	"mjoy.io/common/types"
)

//here for test,do not add msgp
type BalanceValue struct {
	Amount int    `json:"amount"`
}


func MakeActionParamsReword(producer types.Address)[]byte{
	a := make(map[string]interface{})
	a["funcId"] = "1"

	a["producer"] = producer.Hex()

	r , err :=json.Marshal(a)
	if err != nil {
		return nil
	}
	return r
}

