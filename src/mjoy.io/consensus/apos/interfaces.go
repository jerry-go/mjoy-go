package apos

import "sync"

/*
for inner interfaces
*/

type dataSetp interface {

}
//goroutine interfaces for algorand main structure
type stepInterface interface {
	sendMsg(dataPack,*Round) error
	stop()
	run(wg *sync.WaitGroup)
}




