package apos

/*
for inner interfaces
*/

type dataSetp interface {

}
//goroutine interfaces for algorand main structure
type stepInterface interface {
	sendMsg(dataPack) error
	stop()
	run()
}




