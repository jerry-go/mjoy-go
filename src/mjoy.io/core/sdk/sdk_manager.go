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

func (this *SdkManager)




