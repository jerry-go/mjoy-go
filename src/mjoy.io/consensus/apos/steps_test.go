package apos

import "testing"

//instructions:
//this test file just test the steps running right or not,we just send the data which the steps need
//That is , we just need focus on the result of the test


/*
each step obj like below:
type stepInterface interface {
	sendMsg(dataPack,*Round) error
	stop()
	run(wg *sync.WaitGroup)
}





*/


func TestStep2(t *testing.T){
	msger := newMsgManager()

	actualNode := NewApos(msger , newOutCommonTools())
	actualNode.SetOutMsger(msger)

	//make data


}

