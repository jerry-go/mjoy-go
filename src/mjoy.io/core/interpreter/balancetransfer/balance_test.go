package balancetransfer

import (
	"testing"
	"encoding/json"
	"fmt"
)

type Para struct {
	FuncName string
	Para map[string]interface{}
}

func TestJosnInterface(t *testing.T) {
	para := &Para{}
	para.Para = make(map[string]interface{})

	para.Para["from"] = "0x214343434343"
	para.Para["to"] = "0x2222222222"
	para.Para["amount"] = "100"
	para.FuncName = "TransferBalance"

	jsonResult,err := json.Marshal(para)
	fmt.Println(jsonResult,err)


	paraUnmarshal := &Para{}

	err = json.Unmarshal(jsonResult , paraUnmarshal)
	fmt.Println(paraUnmarshal,err)

}
