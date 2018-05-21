package interpreter

import (
	"errors"
	"sync"
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
	"time"
	"fmt"
)

//Test addressd
var FeeCutAddress = types.HexToAddress("0x0000000000000000000000000000000000000001")
var BalanceTransferAddress  = types.HexToAddress("0x0000000000000000000000000000000000000002")

func NewVm()*Vms{
	vm := new(Vms)
	vm.init()
	return vm
}


type Vms struct {

	handlers map[types.Address]HandleFunc   //all address -> handleFunc
	stateDb interface{} //database for account state
	keystore interface{}    //accouts interface
	lock sync.RWMutex   //working  mux
	WorkingChan chan *Work    //work chan
}

func (this *Vms)init(){
	this.handlers = make(map[types.Address]HandleFunc)
	this.stateDb = nil
	this.keystore = nil
	this.WorkingChan = make(chan *Work , 1000)
}



func (this *Vms)RegisterHandlerFunc(codeAddress types.Address , f HandleFunc)error{
	this.lock.Lock()
	defer this.lock.Unlock()
	if _ , ok := this.handlers[codeAddress];ok{
		return errors.New("Exist HandleFunc")
	}

	this.handlers[codeAddress] = f
	return nil
}




/********************************************************************/
//Deal Actions..........
/********************************************************************/
//DealActons is a full work
func (this *Vms)DealActions(pWork *Work)error{
	for _ , a := range pWork.actions{
		err := this.DealAction(pWork.from,a , pWork.resultChan)
		if err != nil{
			wkResult := WorkResult{nil,nil}
			pWork.resultChan <- wkResult  //return the err to the caller
			return err
		}
	}
	return nil
}
//DealActions is a little part of full work
func (this *Vms)DealAction(from types.Address , action transaction.Action ,c <-chan WorkResult)error{
	if f,ok:=this.handlers[*action.Address];ok{
		err := f(action.Params)
		return err
	}
	return errors.New("No Dealing Callbaak")
}

/********************************************************************/
//Deal Work..........
/********************************************************************/
//SendWork is called when applytransaction
func (this *Vms)SendWork(from types.Address , actions []transaction.Action)<-chan WorkResult{
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
		newWork := <-this.WorkingChan
		go this.DealActions(newWork)
	}
}


func (this *Vms)TestRun(){
	for{
		time.Sleep(4*time.Second)
		fmt.Println("For Vm testing print.......")
	}
}




















