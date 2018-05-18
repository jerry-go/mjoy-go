package interpreter

import (
	"errors"
	"sync"
	"mjoy.io/core/transaction"
)

const (
	FeeCut = iota
	BalanceTransfer
)


type Vms struct {

	handlers map[int]HandleFunc
	stateDb interface{}
	keystore interface{}
	lock sync.RWMutex
	WorkingChan <-chan Work
}



func NewVm()*Vms{
	vm := new(Vms)
	return vm
}

func (this *Vms)RegisterHandlerFunc(codeId int , f HandleFunc)error{
	this.lock.Lock()
	defer this.lock.Unlock()
	if _ , ok := this.handlers[codeId];ok{
		return errors.New("Exist HandleFunc")
	}

	this.handlers[codeId] = f
	return nil
}

//DealActons is a full work
func (this *Vms)DealActions(actions []transaction.Action)error{
	for _ , a := range actions{
		err := this.DealAction(a)
		if err != nil{
			return err
		}
	}
	return nil
}
//DealActions is a little part of full work
func (this *Vms)DealAction(action transaction.Action)error{

	return nil
}

//SendWork is called when applytransaction
func (this *Vms)SendWork(tx *transaction.Transaction)<-chan error{

	return nil
}


























