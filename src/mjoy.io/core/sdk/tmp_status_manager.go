package sdk

import (
	"mjoy.io/common/types"
	"sync"
	"mjoy.io/utils/database"
)

/*
sdk is a last system status manager and modification manager,if you want get some last value esaily,
you can call a system function in sdk
*/




//TmpStatusManager,we can get a last value from TmpStatusManager ,or store a modified value (k-v) in TmpStatusManager
//by simple Api (f(stateRoot , contractAddress , key))

type TmpStatusManager struct {
	mu sync.RWMutex
	db database.IDatabaseGetter
	LastRoot types.Hash
	TmpConTracts map[types.Address]*TmpStatusNode
}

func NewTmpStatusManager(root types.Hash , db database.IDatabaseGetter)*TmpStatusManager{
	t := new(TmpStatusManager)
	t.LastRoot = root
	t.db = db
	t.TmpConTracts = make(map[types.Address]*TmpStatusNode)

	return t
}

//SetValue always set into memery
func (this *TmpStatusManager)SetValue(contractAddress types.Address , key []byte , value []byte)error{
	this.mu.Lock()
	defer this.mu.Unlock()

	//step 1: get a statusNode from manager
	statusNode := this.ExistContract(contractAddress)
	if statusNode == nil {
		//if not exist create One
		statusNode = this.CreateStatusNode(contractAddress)
	}

	//step 2: make TmpKey
	tmpKey := TmpKey{contractAddress:contractAddress , key:types.BytesToAddress(key)}

	//step 3:set value
	statusNode.SetValue(tmpKey , value)
	return nil
}


func (this *TmpStatusManager)GetValue(contractAddress types.Address , key []byte)[]byte{
	this.mu.RLock()
	defer this.mu.RUnlock()

	tmpKey := TmpKey{contractAddress:contractAddress , key:types.BytesToAddress(key)}


	tmpNode := this.ExistContract(contractAddress)
	if tmpNode != nil {
		return tmpNode.ExistValue(tmpKey)
	}

	//get data from database, should add the last root
	tmpKey.stateRoot = this.LastRoot
	hashKey , err := tmpKey.MakeHashKey()
	if err != nil{
		return nil
	}



	//if not find in memery,check in the LDB
	data , err := this.db.Get(hashKey[:])
	if err != nil{
		logger.Error("db.Get(hashKey):" , err.Error())
		return nil
	}
	return data
}


//TmpStatusManager basic functions,should not control the mu(lock), the lock should hold by Upper caller


//check a Contract is exist in the tmpStatusManager
func (this *TmpStatusManager)ExistContract(contractAddress types.Address)*TmpStatusNode{
	if node , ok := this.TmpConTracts[contractAddress];ok{
		return node
	}
	return nil
}

func (this *TmpStatusManager)CreateStatusNode(contractAddress types.Address)*TmpStatusNode{
	node := NewStatusNode()
	this.TmpConTracts[contractAddress] = node
	return node
}




