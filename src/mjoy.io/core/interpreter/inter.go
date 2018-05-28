package interpreter

import (
	"errors"
	"sync"
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
	"time"
	"fmt"
	"mjoy.io/core/interpreter/intertypes"
)

//Test addressd
var BalanceTransferAddress  = types.HexToAddress("0x0000000000000000000000000000000000000001")

func NewVm()*Vms{
	vm := new(Vms)
	vm.init()
	return vm
}

type Vms struct {
	pInnerContractMaper *InnerContractManager

	lock sync.RWMutex   //working  mux
	WorkingChan chan *Work    //work chan
	exit        chan struct{}

}


func (this *Vms)init(){
	//init workingChan
	this.WorkingChan = make(chan *Work , 1000)
	//exit chan
	this.exit = make(chan struct{} , 1)
	//init innerContractMapper
	this.pInnerContractMaper = NewInnerContractManager()

}


/********************************************************************/
//Deal Actions..........
/********************************************************************/
//DealActons is a full work,and return a workresult to caller,if get one err ,return
func (this *Vms)DealActions(pWork *Work)error{
	var workResult WorkResult
	workResult.Err = nil
	for _ , a := range pWork.actions{
		workResult.Results = make([]intertypes.ActionResult , 0 )
		workResult.Err = nil

		r , err := this.DealAction(pWork.contractAddress,a )
		if err != nil{
			//get a err,return
			workResult.Results = nil
			workResult.Err = err
			pWork.resultChan <- workResult
			return err
		}

		//get a result
		workResult.Results = append(workResult.Results , r...)
		pWork.resultChan <- workResult
	}

	return nil
}

//DealActions is a little part of full work
func (this *Vms)DealAction(contractAddress types.Address , action transaction.Action )([]intertypes.ActionResult , error){
	if this.pInnerContractMaper.Exist(contractAddress){
		results , err := this.pInnerContractMaper.DoFun(contractAddress , action.Params)
		if err != nil {
			return nil , err
		}
		return results , nil
	}
	return nil , errors.New("innerContract Not Exist....")
}

/********************************************************************/
//Deal Work..........
/********************************************************************/
//SendWork is called when applytransaction
func (this *Vms)SendWork(from types.Address , action transaction.Action)<-chan WorkResult{

	actions := []transaction.Action{}
	actions = append(actions , action)

	w := NewWork(from , actions)
	this.WorkingChan<-w
	return w.resultChan
}

//will uesed by txpool
func (this *Vms)GetPriority(from types.Address , actions []transaction.Action)int{
	//some calculation for priority
	return 10
}


/********************************************************************/
//cycle Dealing
/********************************************************************/
func (this *Vms)Run(){
	go this.TestRun()

	for{
		select  {
		case newWork := <-this.WorkingChan:
			go this.DealActions(newWork)
		case <-this.exit:
			return

		}
	}
}


func (this *Vms)TestRun(){
	for{
		time.Sleep(4*time.Second)
		fmt.Println("For Vm testing print.......")
	}
}




















