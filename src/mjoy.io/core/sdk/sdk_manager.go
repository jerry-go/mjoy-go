package sdk

import (
	"sync"
	"mjoy.io/utils/database"
	"mjoy.io/common/types"
)

type SdkManager struct {
	mu sync.RWMutex
	db database.IDatabaseGetter
	lastroot types.Hash
	pStatusManager *TmpStatusManager
}

var PtrSdkManager *SdkManager

func NewSdkManager(db database.IDatabaseGetter){
	PtrSdkManager = new(SdkManager)
	PtrSdkManager.db = db
	PtrSdkManager.pStatusManager = nil
}

//call this before block producer
func (this *SdkManager)Prepare(lastRoot types.Hash){
	this.mu.Lock()
	defer this.mu.Unlock()

	this.lastroot = lastRoot
	//new a statusManager and run it
	this.pStatusManager = NewTmpStatusManager(this.lastroot , this.db)
}

//call this after block producer
func (this *SdkManager)Down(){
	this.mu.Lock()
	defer this.mu.Unlock()

	//gc should take back the memery
	this.pStatusManager = nil
}




