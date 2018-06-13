package algorand

/*
for inner interfaces
*/

//goroutine interfaces for algorand main structure
type stepInterface interface {
	sendMsg([]byte)error
	stop()
	run()
}




