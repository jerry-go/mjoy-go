package interpreter

import (
	"testing"
	"mjoy.io/core/interpreter/balancetransfer"
	"encoding/json"
	"mjoy.io/utils/database"
	"mjoy.io/core/sdk"
	"mjoy.io/common/types"
	"mjoy.io/core/transaction"
	"fmt"
	"reflect"
)

func checkResultsData(){
	contractAddr := types.Address{}
	contractAddr[0] = 1

	fromAddr := types.Address{}
	fromAddr[2] = 1

	refind := sdk.Sys_GetValue(balancetransfer.BalanceTransferAddress , fromAddr[:])
	if refind == nil{
		fmt.Println("not find data just store before")
		return
	}
	a := new(balancetransfer.BalanceValue)
	err := json.Unmarshal(refind , a)
	if err != nil {
		fmt.Println("err:" , err)
	}

	fmt.Println("account 1 balance:" , a.Amount)

	toAddr := types.Address{}
	toAddr[3] = 1
	refind = sdk.Sys_GetValue(balancetransfer.BalanceTransferAddress , toAddr[:])
	if refind == nil{
		fmt.Println("not find data just store before")
		return
	}

	err = json.Unmarshal(refind , a)
	if err != nil {
		fmt.Println("err:" , err)
	}

	fmt.Println("account 2 balance:" , a.Amount)
}


func makeTestData(){

	//init database
	db,err := database.OpenMemDB()
	if err != nil {
		panic(err)
	}
	sdk.NewSdkManager(db)

	a := new(balancetransfer.BalanceValue)
	a.Amount = 1000

	lastAccountInfoData , err := json.Marshal(a)
	if err != nil {
		return
	}
	//store the data
	sdk.PtrSdkManager.Prepare(types.Hash{})
	contractAddr := types.Address{}
	contractAddr[0] = 1

	accountAddr := types.Address{}
	accountAddr[2] = 1
	sdk.Sys_SetValue(balancetransfer.BalanceTransferAddress , accountAddr[:] , lastAccountInfoData)
	//sdk.PtrSdkManager.Down()
	refind := sdk.Sys_GetValue(balancetransfer.BalanceTransferAddress , accountAddr[:])
	if refind == nil{
		fmt.Println("not find data just store before")
	}
	fmt.Printf("get store data before :%x\n" , refind)
}

func makeActionParams()[]byte{
	a := make(map[string]interface{})
	a["funcId"] = 1

	fromAddr := types.Address{}
	fromAddr[2] = 1
	a["from"] = fromAddr.Hex()

	toAddr := types.Address{}
	toAddr[3] = 1
	a["to"] = toAddr.Hex()

	a["amount"] = int64(10)

	fmt.Println("type amount:" , reflect.TypeOf(a["amount"]))
	r , err :=json.Marshal(a)
	if err != nil {
		return nil
	}
	return r
}


func TestInterDbNoDataBefore(t *testing.T){
	//make test data
	makeTestData()
	checkResultsData()
	action:= transaction.Action{}
	contranctAddr := balancetransfer.BalanceTransferAddress
	action.Address = &contranctAddr
	action.Params = makeActionParams()

	pNewVm := NewVm()
	//sdk.PtrSdkManager.Prepare(types.Hash{})
	fmt.Println("Start Testing....")
	rChan := pNewVm.SendWork(types.Address{} , action)
	rw := <-rChan
	fmt.Println("get A result")
	fmt.Println("resultsLen :" , len(rw.Results))
	fmt.Println("err:" , rw.Err)
	_ = rw
	checkResultsData()

}