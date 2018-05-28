package balancetransfer

import (
	"errors"
	"fmt"
	"mjoy.io/common/types"
	"mjoy.io/core/sdk"
	"encoding/json"
	"mjoy.io/core/interpreter/intertypes"
)


func CreateTransferBalance(){}




func TransferBalance(param map[string]interface{})([]intertypes.ActionResult , error){
	var from string
	var fromAddress types.Address
	var to string
	var toAddress types.Address
	var amount int64


	//get params
	//from
	if fromi,ok := param["from"];ok{
		from = fromi.(string)
		fromAddress = types.StringToAddress(from)
	}else{
		return nil ,errors.New(fmt.Sprintf("TransferBalance:param no index:from"))
	}

	//to
	if toi , ok := param["to"];ok{
		to = toi.(string)
		toAddress = types.StringToAddress(to)
	}else {
		return nil , errors.New(fmt.Sprintf("TransferBalance:param no index:to"))
	}

	//amount
	if amounti , ok := param["amount"];ok{
		amount = amounti.(int64)
	}

	//logicDeal
	//get sender's Balance
	dataFrom := sdk.Sys_GetValue(BalanceTransferAddress , fromAddress[:])
	if nil == dataFrom{
		return nil , errors.New("TransferBalance:Do not find data:From")
	}

	balanceFrom := new(BalanceValue)
	err := json.Unmarshal(dataFrom , balanceFrom)
	if err != nil {
		return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
	}

	//balance value check
	if balanceFrom.Amount < amount{
		return nil , errors.New(fmt.Sprintf("TransferBalance:has %d , but want %d" , balanceFrom.Amount , amount))
	}

	//get receiver's Balance
	dataTo := sdk.Sys_GetValue(BalanceTransferAddress , toAddress[:])
	if nil == dataTo{
		return nil , errors.New("TransferBalance:Do not find data:To")
	}

	balanceTo := new(BalanceValue)
	err = json.Unmarshal(dataTo , balanceTo)
	if err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
	}

	//balance modify
	balanceFrom.Amount -= amount
	balanceTo.Amount += amount
	//set value to database(by sys_xxx call,setting into memery)
	// 1. marshal data
	bytesFrom , err := json.Marshal(balanceFrom)
	if err != nil {
		return nil , errors.New(fmt.Sprintf("TransferBalance:Marshal json:%s" , err.Error()))
	}

	bytesTo , err := json.Marshal(balanceTo)
	if err != nil {
		return nil , errors.New(fmt.Sprintf("TransferBalance:Marshal json:%s" , err.Error()))
	}
	if err = sdk.Sys_SetValue(BalanceTransferAddress , fromAddress[:] , bytesFrom);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set From :%s" , err.Error()))
	}

	if err = sdk.Sys_SetValue(BalanceTransferAddress , toAddress[:] , bytesTo);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set To :%s" , err.Error()))
	}

	//make a result
	results := make([]intertypes.ActionResult , 2)
	results = append(results , intertypes.ActionResult{Key:fromAddress[:] , Val:bytesFrom})
	results = append(results , intertypes.ActionResult{Key:toAddress[:] , Val:bytesTo})

	return results , nil
}






