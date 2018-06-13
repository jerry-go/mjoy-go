package apos

import "fmt"

//steps handle

//step 1:Block Proposal
type step1BlockProposal struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int            //which step the obj stay
	apos *Apos

}
func makeStep1Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step1BlockProposal{
	s := new(step1BlockProposal)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)

	s.step = step
	return s
}
func (this *step1BlockProposal)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step1BlockProposal)stop(){
	this.exit<-1
}

func (this *step1BlockProposal)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 2:First step of GC
type step2FirstStepGC struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep2Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step2FirstStepGC{
	s := new(step2FirstStepGC)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step2FirstStepGC)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step2FirstStepGC)stop(){
	this.exit<-1
}

func (this *step2FirstStepGC)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}


//step 3:Second Step of GC
type step3SecondStepGC struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep3Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step3SecondStepGC{
	s := new(step3SecondStepGC)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step3SecondStepGC)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step3SecondStepGC)stop(){
	this.exit<-1
}

func (this *step3SecondStepGC)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}


//step 4:First Step of BBA*
type step4FirstStepBBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep4Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step4FirstStepBBA{
	s := new(step4FirstStepBBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step4FirstStepBBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step4FirstStepBBA)stop(){
	this.exit<-1
}

func (this *step4FirstStepBBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}


//step 5<= s <= m+2 ,s-2 mod 3 == 0 mod 3:Coin-Fixed-To-0 step of BBA*
type step5CoinFixedTo0BBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep5Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step5CoinFixedTo0BBA{
	s := new(step5CoinFixedTo0BBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step5CoinFixedTo0BBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step5CoinFixedTo0BBA)stop(){
	this.exit<-1
}

func (this *step5CoinFixedTo0BBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 6<= s <= m+2 ,s-2 mod 3 == 1 mod 3:Coin-Fixed-To-1 step of BBA*
type step6CoinFixedTo1BBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep6Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step6CoinFixedTo1BBA{
	s := new(step6CoinFixedTo1BBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step6CoinFixedTo1BBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step6CoinFixedTo1BBA)stop(){
	this.exit<-1
}

func (this *step6CoinFixedTo1BBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 7<= s <= m+2 ,s-2 mod 3 == 2 mod 3:Coin-Fixed-To-1 step of BBA*
type step7CoinGenFlipBBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep7Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step7CoinGenFlipBBA{
	s := new(step7CoinGenFlipBBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step7CoinGenFlipBBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step7CoinGenFlipBBA)stop(){
	this.exit<-1
}

func (this *step7CoinGenFlipBBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}














