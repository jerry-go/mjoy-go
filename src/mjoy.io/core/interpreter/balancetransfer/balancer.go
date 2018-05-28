package balancetransfer

import (
	"encoding/json"
	"errors"
	"fmt"
	"mjoy.io/common/types"
	"mjoy.io/core/interpreter/intertypes"
)

var BalanceTransferAddress  = types.StringToAddress("0x0000000000000000000000000000000000000001")


type DoFunc func(map[string]interface{})([]intertypes.ActionResult , error)

type ContractBalancer struct {
	funcMapper map[int]DoFunc
}
//managed by vm
func NewContractBalancer()*ContractBalancer{
	b := new(ContractBalancer)
	b.init()
	return b
}

func (this *ContractBalancer)init(){
	//register call Back
	this.funcMapper = make(map[int]DoFunc)
	this.funcMapper[0] = TransferBalance
}



func ParseParms(param []byte)(map[string]interface{} , error){
	pResult:= make(map[string]interface{})
	err := json.Unmarshal(param , &pResult)
	if err != nil{
		return nil , err
	}
	return pResult , nil

}

func (this *ContractBalancer)DoFun( params []byte)([]intertypes.ActionResult , error){
	//unmarshal params
	jsonParams , err := ParseParms(params)
	if err != nil {
		return nil,err
	}
	var funcId int64
	if v,ok := jsonParams["funcId"];!ok{
		funcId = v.(int64)
		return nil , errors.New(fmt.Sprintf("ContractBalancer: Params not contain funcId" ))
	}

	if doFunc,ok := this.funcMapper[int(funcId)];ok {
		return doFunc(jsonParams)
	}

	return nil , errors.New(fmt.Sprintf("ContractBalancer: no Func Id:%d find in map" , funcId))
}






