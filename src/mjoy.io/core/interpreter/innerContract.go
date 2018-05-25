/*
This file define a InnerContract interface
*/

package interpreter

import (
	"sync"
	"mjoy.io/common/types"
	"math/big"
)

//InnerContrancInterface
type InnerContract interface {
	DoFun(fName string , params interface{})(interface{} , error)
}

//InnerContranctMap is a innerContract controller,like check contract ,do a contract
type InnerContractMaper struct {
	mu sync.RWMutex
	Inners map[types.Address]InnerContract
}

//New A InnerContractMaper
func NewInnerContractMaper()*InnerContractMaper{
	maper := new(InnerContractMaper)
	maper.Inners = make(map[types.Address]InnerContract)
	maper.init()
	return maper
}

func (this *InnerContractMaper)init(){
	this.register()
}

var zeroAddress = types.BigToAddress(big.NewInt(0))

func (this *InnerContractMaper)register(){
	this.mu.Lock()
	defer this.mu.Unlock()

	if len(allInnerRegister) == 0 {
		logger.Error("InnerContractMaper register len == 0")
		return
	}

	for _ , obj := range allInnerRegister {
		if obj.address != zeroAddress{
			this.Inners[obj.address] = obj.inner
		}
	}
}


//check innerContract is exist
func (this *InnerContractMaper)Exist(address types.Address)bool{
	this.mu.RLock()
	defer this.mu.RUnlock()

	if _ , ok := this.Inners[address];ok{
		return true
	}
	return false
}

//call a innerContract.Please call Exist ensure a innerContract is exist or not before this
func (this *InnerContractMaper)DoFun(address types.Address , fname string , params interface{})(interface{} ,  error){
	inner := this.Inners[address]
	return inner.DoFun(fname , params)
}


