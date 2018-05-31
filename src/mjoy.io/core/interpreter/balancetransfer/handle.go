package balancetransfer

import (
	"errors"
	"fmt"
	"mjoy.io/common/types"
	"mjoy.io/core/sdk"
	"encoding/json"
	"mjoy.io/core/interpreter/intertypes"
	"strconv"
)

func TransferFee(param map[string]interface{} , sysparam *intertypes.SystemParams)([]intertypes.ActionResult , error){
	var from string
	var fromAddress types.Address
	var to string
	var toAddress types.Address
	var feeAmount int

	logger.Debug(">>>>>>>>StartTransferBalanceDeal.............................................")
	//get params
	//from
	if fromi,ok := param["from"];ok{
		from = fromi.(string)
		fromAddress = types.HexToAddress(from)
	}else{
		return nil ,errors.New(fmt.Sprintf("TransferBalance:param no index:from"))
	}

	//to
	if ptoAddr  := sdk.Sys_GetCoinbase(sysparam.SdkHandler);ptoAddr == nil {
		return nil , errors.New("Sys_GetCoinbase return Nil")
	}else{
		toAddress = *ptoAddr
	}


	//Fee amount
	if amounti , ok := param["amount"];ok{
		feeAmount,_ = strconv.Atoi(amounti.(string))
	}

	//logicDeal
	//get sender's Balance
	dataFrom := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , fromAddress[:])
	if nil == dataFrom{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Do not find data:From:%x" , fromAddress))
	}

	balanceFrom := new(BalanceValue)
	err := json.Unmarshal(dataFrom , balanceFrom)
	if err != nil {
		return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
	}
	fmt.Printf(".........................................Sender Address:%s   Balance:%d\n" , from , balanceFrom.Amount)
	//balance value check
	if balanceFrom.Amount < feeAmount{
		return nil , errors.New(fmt.Sprintf("TransferBalance:has %d , but want %d" , balanceFrom.Amount , feeAmount))
	}

	//get receiver's Balance
	dataTo := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , toAddress[:])
	balanceTo := new(BalanceValue)

	if nil == dataTo{
		//return nil , errors.New("TransferBalance:Do not find data:To")
		logger.Warnf("DataTo %s  dataTo==nil,mean no balance before!!!!" , toAddress.Hex())
		balanceTo.Amount = 0
	}else{
		err = json.Unmarshal(dataTo , balanceTo)
		if err != nil{
			return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
		}
	}

	fmt.Printf("..........................Receiver Address:%s   Balance:%d\n" , to , balanceTo.Amount)
	//balance modify
	balanceFrom.Amount -= feeAmount
	balanceTo.Amount += feeAmount
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
	if err = sdk.Sys_SetValue(sysparam.SdkHandler ,  BalanceTransferAddress , fromAddress[:] , bytesFrom);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set From :%s" , err.Error()))
	}

	if err = sdk.Sys_SetValue(sysparam.SdkHandler ,  BalanceTransferAddress , toAddress[:] , bytesTo);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set To :%s" , err.Error()))
	}

	newdataTo := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , toAddress[:])
	if newdataTo == nil {
		logger.Error("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!newDataTo == nil")
	}else{
		newbalanceTo := new(BalanceValue)
		err  := json.Unmarshal(newdataTo , newbalanceTo)
		if err != nil {
			logger.Errorf("!!!!!!!!!!!!!!!!!!!NewBalanceTo Unmarshal Err:%s" ,  err.Error())
		}else{
			fmt.Printf("...........>>>>>>...............New Receiver Address:%s   Balance:%d\n" , to , newbalanceTo.Amount)
		}

	}

	//make a result
	results := make([]intertypes.ActionResult , 2)
	results = results[:0]
	results = append(results , intertypes.ActionResult{Key:fromAddress[:] , Val:bytesFrom})
	results = append(results , intertypes.ActionResult{Key:toAddress[:] , Val:bytesTo})

	return results , nil
}

func TransferBalance(param map[string]interface{},sysparam *intertypes.SystemParams)([]intertypes.ActionResult , error){
	var from string
	var fromAddress types.Address
	var to string
	var toAddress types.Address
	var amount int

	logger.Debug(">>>>>>>>StartTransferBalanceDeal.............................................")
	//get params
	//from
	if fromi,ok := param["from"];ok{
		from = fromi.(string)
		fromAddress = types.HexToAddress(from)
	}else{
		return nil ,errors.New(fmt.Sprintf("TransferBalance:param no index:from"))
	}

	//to
	if toi , ok := param["to"];ok{
		to = toi.(string)
		toAddress = types.HexToAddress(to)
	}else {
		return nil , errors.New(fmt.Sprintf("TransferBalance:param no index:to"))
	}

	//amount
	if amounti , ok := param["amount"];ok{
		amount,_ = strconv.Atoi(amounti.(string))
	}

	//logicDeal
	//get sender's Balance
	dataFrom := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , fromAddress[:])
	if nil == dataFrom{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Do not find data:From:%x" , fromAddress))
	}

	balanceFrom := new(BalanceValue)
	err := json.Unmarshal(dataFrom , balanceFrom)
	if err != nil {
		return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
	}
	fmt.Printf(".........................................Sender Address:%s   Balance:%d\n" , from , balanceFrom.Amount)
	//balance value check
	if balanceFrom.Amount < amount{
		return nil , errors.New(fmt.Sprintf("TransferBalance:has %d , but want %d" , balanceFrom.Amount , amount))
	}

	//get receiver's Balance
	dataTo := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , toAddress[:])
	balanceTo := new(BalanceValue)

	if nil == dataTo{
		//return nil , errors.New("TransferBalance:Do not find data:To")
		logger.Warnf("DataTo %s  dataTo==nil,mean no balance before!!!!" , toAddress.Hex())
		balanceTo.Amount = 0
	}else{
		err = json.Unmarshal(dataTo , balanceTo)
		if err != nil{
			return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
		}
	}

	fmt.Printf("..........................Receiver Address:%s   Balance:%d\n" , to , balanceTo.Amount)
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
	if err = sdk.Sys_SetValue(sysparam.SdkHandler ,  BalanceTransferAddress , fromAddress[:] , bytesFrom);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set From :%s" , err.Error()))
	}

	if err = sdk.Sys_SetValue(sysparam.SdkHandler ,  BalanceTransferAddress , toAddress[:] , bytesTo);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set To :%s" , err.Error()))
	}

	newdataTo := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , toAddress[:])
	if newdataTo == nil {
		logger.Error("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!newDataTo == nil")
	}else{
		newbalanceTo := new(BalanceValue)
		err  := json.Unmarshal(newdataTo , newbalanceTo)
		if err != nil {
			logger.Errorf("!!!!!!!!!!!!!!!!!!!NewBalanceTo Unmarshal Err:%s" ,  err.Error())
		}else{
			fmt.Printf("...........>>>>>>...............New Receiver Address:%s   Balance:%d\n" , to , newbalanceTo.Amount)
		}

	}

	//make a result
	results := make([]intertypes.ActionResult , 2)
	results = results[:0]
	results = append(results , intertypes.ActionResult{Key:fromAddress[:] , Val:bytesFrom})
	results = append(results , intertypes.ActionResult{Key:toAddress[:] , Val:bytesTo})

	return results , nil
}

func RewordBlockProducer(param map[string]interface{},sysparam *intertypes.SystemParams)([]intertypes.ActionResult , error){

	var producer types.Address
	if fromi,ok := param["producer"];ok{
		producerStr := fromi.(string)
		producer = types.HexToAddress(producerStr)
	}else{
		return nil ,errors.New(fmt.Sprintf("no producer"))
	}

	balance := new(BalanceValue)
	data := sdk.Sys_GetValue(sysparam.SdkHandler ,  BalanceTransferAddress , producer[:])
	if nil == data{
		//return nil , errors.New(fmt.Sprintf("TransferBalance:Do not find data:From:%x" , producer))
		balance.Amount = 5e+5
	} else {
		err := json.Unmarshal(data , balance)
		if err != nil {
			return nil , errors.New(fmt.Sprintf("TransferBalance:Unmarshal json:%s" , err.Error()))
		}
		balance.Amount += 5e+5
	}

	fmt.Println("RewordBlockProducer", balance.Amount)
	bytesJosn , err := json.Marshal(balance)
	if err != nil {
		return nil , errors.New(fmt.Sprintf("TransferBalance:Marshal json:%s" , err.Error()))
	}

	if err = sdk.Sys_SetValue(sysparam.SdkHandler ,  BalanceTransferAddress , producer[:] , bytesJosn);err != nil{
		return nil , errors.New(fmt.Sprintf("TransferBalance:Set From :%s" , err.Error()))
	}

	//make a result
	results := make([]intertypes.ActionResult , 1)
	results = results[:0]
	results = append(results , intertypes.ActionResult{Key:producer[:] , Val:bytesJosn})

	return results , nil

}






