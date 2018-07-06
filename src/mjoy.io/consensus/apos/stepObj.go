package apos

import (
	"sync"
	"time"
	"mjoy.io/common/types"
	"math/big"
	"fmt"
	"container/heap"
	"mjoy.io/core/blockchain/block"
)


var (
	LessTimeDelayFlag bool = false  //let step spend less time to deal msg received
	LessTimeDelayCnt int = 5
)

func calcTTL(step int) time.Duration {
	delayTm := (  (step - 1)*2 + 1  )
	//fmt.Println("step:" , step , "    delayTm:  " , delayTm)
	return time.Duration(Config().blockDelay + delayTm* Config().verifyDelay) * time.Second
}

//step 1
type stepObj1 struct {
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
}

func makeStepObj1()*stepObj1{
	s := new(stepObj1)
	return s
}

func (this *stepObj1)setCtx(ctx *stepCtx){
	this.ctx = ctx
	if this.timeEnd == 0 {
		this.timeEnd = calcTTL(ctx.getStep())
	}
}

func (this *stepObj1)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj1)timerHandle(){
	defer func(){
		go this.ctx.stopStep()
	}()
	//new a M1 data
	m1 := newBlockProposal()


	m1.Credential = this.ctx.getCredential()
	fmt.Println("m1.step:",m1.Credential.Step)
	//m1.Block = this.ctx.makeEmptyBlockForTest(m1.Credential)
	fmt.Println("!!!!!!!!!!WILL MAKE BLOCK..............")

	bcd := &block.ConsensusData{}
	bcd.Id = ConsensusDataId
	bcd.Para = m1.Credential.Signature.toBytes()
	m1.Block = this.ctx.getProducerNewBlock(bcd)
	fmt.Println(m1.Block)

	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = 1
	m1.Esig.val = make([]byte,0)
	h := m1.Block.Hash()
	m1.Esig.val = append(m1.Esig.val , h[:]...)
	err := this.ctx.esig(m1.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}


	//fill struct members
	//todo: should using interface
	this.ctx.sendInner(m1)

	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***[A]Out M1 CreHash:" ,m1.Credential.Signature.Hash().String()," BlockHash", m1.Block.B_header.Hash().String(),COLOR_SHORT_RESET)
	//logger.Debug("\033[31;m [A]Out M1 \033[0m  Block Hash:" , m1.Block.B_header.Hash().String())

}

func (this *stepObj1)dataHandle(data interface{}){

}

func (this *stepObj1)stopHandle(){

}

//step 2
type stepObj2 struct {
	smallestLBr *BlockProposal //this node regard smallestLBr as the smallest credential's block info
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
}

func makeStepObj2()*stepObj2{
	s := new(stepObj2)
	return s
}

func (this * stepObj2)setCtx(ctx *stepCtx){
	this.ctx = ctx
	if this.timeEnd == 0 {
		this.timeEnd = calcTTL(ctx.getStep())
	}
}

func (this *stepObj2)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj2)timerHandle(){
	defer func(){
		go this.ctx.stopStep()
	}()
	m2 := newGradedConsensus()
	m2.Credential = this.ctx.getCredential()
	if m2.Credential == nil {
		panic("timeHandle Credential Wrong.........")
	}
	fmt.Println("m2.Credential:" , *m2.Credential)
	if this.smallestLBr == nil {
		m2.Hash = types.Hash{}
		logger.Error(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***Obj2[A]Step:",this.ctx.getCredential().Step,"Out M2 By a EmptyHash",COLOR_SHORT_RESET)
	}else{
		m2.Hash = this.smallestLBr.Block.Hash()
	}

	m2.Esig.round = m2.Credential.Round
	m2.Esig.step = m2.Credential.Step
	m2.Esig.val = make([]byte , 0)
	m2.Esig.val = append(m2.Esig.val , m2.Hash[:]...)

	err := this.ctx.esig(m2.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	this.ctx.sendInner(m2)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***Obj2[A]Step:",this.ctx.getCredential().Step,"Out M2-Hash:",m2.Hash.String(),COLOR_SHORT_RESET)

	//turn to stop

}

func (this *stepObj2)dataHandle(data interface{}){
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"Obj2[A]Step:",this.ctx.getCredential().Step,"In M1",COLOR_SHORT_RESET)
	m1 := new(BlockProposal)
	m1 = data.(*BlockProposal)

	if m1 == nil{
		return
	}
	if this.smallestLBr == nil {
		this.smallestLBr = m1
		this.ctx.propagateMsg(m1)
		return
	}
	//compare M1 before and M1 current
	r := this.smallestLBr.Credential.Cmp(m1.Credential)
	if r > 0 {
		//exchange smallestLBr and m1
		this.smallestLBr = m1
		this.ctx.propagateMsg(m1)
	}
}

func (this *stepObj2)stopHandle(){
	//todo:what last operations should do
}

//step 3
type stepObj3 struct {

	allM2Index map[types.Hash]map[CredentialSigForKey]bool
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
}

func makeStepObj3()*stepObj3{
	s := new(stepObj3)
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
	return s
}

func (this *stepObj3)setCtx(ctx *stepCtx){
	this.ctx = ctx
	if this.timeEnd == 0 {
		this.timeEnd = calcTTL(ctx.getStep())
	}
}

func (this *stepObj3)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj3)timerHandle(){
	defer func(){
		go this.ctx.stopStep()
	}()

	this.lock.Lock()
	defer this.lock.Unlock()



	//time to work now,send all
	total:=0
	maxLen := 0
	maxHash := types.Hash{}


	for hash,supporter := range this.allM2Index{

		currentLen := len(supporter)
		if currentLen > maxLen{
			maxLen = currentLen
			maxHash = hash
		}
		total += currentLen

	}
	v := types.Hash{}

	if maxLen * 3 > 2*total{
		v = maxHash
	}
	//pack m3 Data

	m3 := newGradedConsensus()
	m3.Credential = this.ctx.getCredential()

	m3.Hash = v

	m3.Esig.round = m3.Credential.Round
	m3.Esig.step = m3.Credential.Step
	m3.Esig.val = make([]byte , 0)
	m3.Esig.val = append(m3.Esig.val , m3.Hash[:]...)

	err := this.ctx.esig(m3.Esig)
	if err != nil{
		logger.Error(err.Error())
		return
	}


	this.ctx.sendInner(m3)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***Obj3[A]Step:",this.ctx.getCredential().Step,"Out M3-Hash:",m3.Hash.String(),COLOR_SHORT_RESET)

}

func (this *stepObj3)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()

	m2 := newGradedConsensus()
	m2 = data.(*GradedConsensus)
	if m2 == nil{
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"Obj3[A]Step:",this.ctx.getCredential().Step,"In M2",m2.Hash.String(),COLOR_SHORT_RESET)
	//add to IndexMap
	var subIndex map[CredentialSigForKey]bool
	subIndex = this.allM2Index[m2.Hash]
	if subIndex == nil {
		this.allM2Index[m2.Hash] = make(map[CredentialSigForKey]bool)
		subIndex = this.allM2Index[m2.Hash]
	}

	sigKey := *m2.Credential.ToCredentialSigKey()
	if _ , ok := subIndex[sigKey];!ok{
		subIndex[sigKey] = true
	}
}

func (this *stepObj3)stopHandle(){

}


//step 4
type stepObj4 struct {

	allM2Index map[types.Hash]map[CredentialSigForKey]bool
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
}

func makeStepObj4()*stepObj4{
	s := new(stepObj4)
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
	return s
}

func (this *stepObj4)setCtx(ctx *stepCtx){
	this.ctx = ctx
	//if test
	if LessTimeDelayFlag{
		this.timeEnd = time.Duration(LessTimeDelayCnt)
	}

	if this.timeEnd == 0 {
		this.timeEnd = calcTTL(ctx.getStep())
	}
}

func (this *stepObj4)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj4)timerHandle(){

	defer func(){
		go this.ctx.stopStep()
	}()
	this.lock.Lock()
	defer this.lock.Unlock()


	//time to work now,send all
	total:=0
	maxLen := 0
	maxHash := types.Hash{}


	for hash,supporter := range this.allM2Index{

		currentLen := len(supporter)
		if currentLen > maxLen{
			maxLen = currentLen
			maxHash = hash
		}
		total += currentLen

	}
	//todo:should be H(Be)
	v := this.ctx.getEmptyBlockHash()

	g := 0
	fmt.Println("MaxLen:" , maxLen)
	fmt.Println("total:" , total)
	if maxLen * 3 > 2 * total{
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"Step4  Do :maxLen * 3 > 2 * total,g=2",COLOR_SHORT_RESET)
		v = maxHash
		g = 2
	}else if maxLen * 3 > total{
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"Step4  Do :maxLen * 3 > total,g=1",COLOR_SHORT_RESET)
		v = maxHash
		g = 1
	}else{
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"Step4  Do :Else,g=0",COLOR_SHORT_RESET)
		g = 0
	}

	b := 0

	if g == 2 {
		b = 0
	}else{
		b = 1
	}

	//pack m4 Data
	m4 := newBinaryByzantineAgreement()

	m4.Hash = v
	m4.B = uint(b)
	m4.Credential = this.ctx.getCredential()

	//sig b
	m4.EsigB.round = m4.Credential.Round
	m4.EsigB.step = m4.Credential.Step
	m4.EsigB.val = make([]byte , 0)
	m4.EsigB.val = append(m4.EsigB.val , big.NewInt(int64(m4.B)).Bytes()...)

	err := this.ctx.esig(m4.EsigB)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	//v

	m4.EsigV.round = m4.Credential.Round
	m4.EsigV.step = m4.Credential.Step
	m4.EsigV.val = make([]byte , 0)
	m4.EsigV.val = append(m4.EsigV.val , m4.Hash[:]...)

	err = this.ctx.esig(m4.EsigV)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	this.ctx.sendInner(m4)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***Obj4[A]Step:",this.ctx.getCredential().Step,"Out M4-Hash:",m4.Hash.String(),"   B:",m4.B,COLOR_SHORT_RESET)



}

func (this *stepObj4)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()
	m3 := newGradedConsensus()
	m3 = data.(*GradedConsensus)
	if m3 == nil{
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"Obj4[A]Step:",this.ctx.getCredential().Step,"In M3",m3.Hash.String(),COLOR_SHORT_RESET)
	//add to IndexMap
	var subIndex map[CredentialSigForKey]bool
	subIndex = this.allM2Index[m3.Hash]
	if subIndex == nil {
		this.allM2Index[m3.Hash] = make(map[CredentialSigForKey]bool)
		subIndex = this.allM2Index[m3.Hash]
	}
	sigKey := *m3.Credential.ToCredentialSigKey()
	if _ , ok := subIndex[sigKey];!ok{
		subIndex[sigKey] = true
	}

}
func (this *stepObj4)stopHandle(){

}


//step 567
type stepObj567 struct {

	stepIndex int   //will calculate in setCtx

	allMxIndex map[types.Hash]*binaryStatus
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
}


func makeStepObj567()*stepObj567{
	s := new(stepObj567)
	s.allMxIndex = make(map[types.Hash]*binaryStatus)
	return s
}

func (this *stepObj567)setCtx(ctx *stepCtx){
	this.ctx = ctx

	if this.timeEnd == 0 {
		this.timeEnd = calcTTL(ctx.getStep())
	}
	//init step and stepIndex
	step := this.ctx.getStep()

	if ((step - 2)%3 == 0%3) {
		this.stepIndex = 5
	}else if  ((step - 2)%3 == 1%3) {
		this.stepIndex = 6
	}else if  ((step - 2)%3 == 2%3) {
		this.stepIndex = 7
	}

}

func (this *stepObj567)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj567)timerHandle(){

	defer func(){
		go this.ctx.stopStep()
	}()

	this.lock.Lock()
	defer this.lock.Unlock()


	//time to work now,send all
	total:=0
	maxLen := 0
	maxHash := types.Hash{}

	max1Len := 0
	max0Len := 0


	for hash,bStatus := range this.allMxIndex{
		//add total cnt
		currentTotalLen := bStatus.getTotalCnt()

		if currentTotalLen > maxLen{
			maxLen = currentTotalLen
			maxHash = hash
		}

		total += currentTotalLen
	}

	maxBStatus := this.allMxIndex[maxHash]
	if maxBStatus == nil {
		return
	}
	max0Len = maxBStatus.getCnt(0)
	max1Len = maxBStatus.getCnt(1)




	mx := newBinaryByzantineAgreement()

	//check 2/3 0 and 2/3 1
	if max0Len * 3 > 2 * total {
		mx.Hash = maxHash
		mx.Credential = this.ctx.getCredential()
		mx.B = 0
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
			"StepCommon  Do :max0Len * 3 > 2 * total,B=0",
			COLOR_SHORT_RESET)
	}else if max1Len * 3 > 2 * total {
		mx.Hash = maxHash
		mx.Credential = this.ctx.getCredential()
		mx.B = 1
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
			"StepCommon  Do :max0Len * 3 > 2 * total,B=1",
			COLOR_SHORT_RESET)
	}else {
		mx.Hash = maxHash
		mx.Credential = this.ctx.getCredential()
		switch this.stepIndex {
		case 5:
			mx.B = 0
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
				"StepCommon  Do :else %5,B=0",
				COLOR_SHORT_RESET)
		case 6:
			mx.B = 1
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
				"StepCommon  Do :else %6,B=1",
				COLOR_SHORT_RESET)
		case 7:
			{
				cHeap := new(CredentialSigStatusHeap)
				allCnt := 0
				for _,bStatus := range this.allMxIndex{
					allCredential := bStatus.export0Credential()
					// 0 credential
					for _,c := range allCredential{
						allCnt++
						*cHeap = append(*cHeap , &CredentialSigStatus{c:*c.ToCredentialSig() , v:0})
					}
					allCredential = bStatus.export1Credential()
					// 1 credential
					for _,c := range allCredential{
						allCnt++
						*cHeap = append(*cHeap , &CredentialSigStatus{c:*c.ToCredentialSig() , v:1})
					}

				}
				fmt.Println("...................All Mx Index:" , allCnt)
				heap.Init(cHeap)
				little := heap.Pop(cHeap).(*CredentialSigStatus)

				mx.B = uint(little.v)

				logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
					"StepCommon  Do :else %7,B=",mx.B,
					COLOR_SHORT_RESET)
			}
		}
	}


	//b big.Int
	mx.EsigB.round = mx.Credential.Round
	mx.EsigB.step = mx.Credential.Step
	mx.EsigB.val = make([]byte , 0)
	mx.EsigB.val = append(mx.EsigB.val , big.NewInt(int64(mx.B)).Bytes()...)

	err := this.ctx.esig(mx.EsigB)
	if err != nil{
		logger.Error(err.Error())
		return
	}


	//v
	mx.EsigV.round = mx.Credential.Round
	mx.EsigV.step = mx.Credential.Step
	mx.EsigV.val = make([]byte , 0)
	mx.EsigV.val = append(mx.EsigV.val , mx.Hash[:]...)

	err = this.ctx.esig(mx.EsigV)
	if err != nil{
		logger.Error(err.Error())
		return
	}

	this.ctx.sendInner(mx)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***Obj567[A]Step:",this.ctx.getCredential().Step,"Out M",mx.Credential.Step,"-Hash:",mx.Hash.String() , "  B:",mx.B,COLOR_SHORT_RESET)


}


func (this *stepObj567)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()
	m6 := newBinaryByzantineAgreement()
	m6 = data.(*BinaryByzantineAgreement)
	if m6 == nil{
		return
	}
	//add to IndexMap
	//logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M",m6.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"Obj567[A]Step:",this.ctx.getCredential().Step,"In M3 Hash:",m6.Hash.String(),"B:",m6.B,COLOR_SHORT_RESET)
	var subIndex *binaryStatus
	subIndex = this.allMxIndex[m6.Hash]
	if subIndex == nil {
		subIndex = makeBinaryStatus()
		this.allMxIndex[m6.Hash] = subIndex
	}
	//check sig
	//set status
	subIndex.setToStatus(*m6.Credential , int(m6.B))

}

func (this *stepObj567)stopHandle(){

}




//step m3
type stepObjm3 struct {

	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
}

func makeStepObjm3()*stepObjm3{
	s := new(stepObjm3)
	return s
}

func (this *stepObjm3)setCtx(ctx *stepCtx){
	if this.timeEnd == 0 {
		this.timeEnd = calcTTL(ctx.getStep())
	}
	this.ctx = ctx
}

func (this *stepObjm3)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObjm3)timerHandle(){

	defer func(){
		go this.ctx.stopStep()
	}()

	this.lock.Lock()
	defer this.lock.Unlock()

	m3 := newBinaryByzantineAgreement()
	//todo:should be H(Be)
	m3.Hash = this.ctx.getEmptyBlockHash()
	m3.B = 1
	m3.Credential = this.ctx.getCredential()

	//b big.Int
	m3.EsigB.round = m3.Credential.Round
	m3.EsigB.step = m3.Credential.Step
	m3.EsigB.val = make([]byte , 0)
	m3.EsigB.val = append(m3.EsigB.val , big.NewInt(int64(m3.B)).Bytes()...)

	err := this.ctx.esig(m3.EsigB)
	if err != nil{
		logger.Error(err.Error())
		return
	}


	//v
	m3.EsigV.round = m3.Credential.Round
	m3.EsigV.step = m3.Credential.Step
	m3.EsigV.val = make([]byte , 0)
	m3.EsigV.val = append(m3.EsigV.val , m3.Hash[:]...)

	err = this.ctx.esig(m3.EsigV)
	if err != nil{
		logger.Error(err.Error())
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"***ObjM3[A]Step:","Out-Hash:",m3.Hash.String(),"B:",m3.B,COLOR_SHORT_RESET)
	this.ctx.sendInner(m3)


}


func (this *stepObjm3)dataHandle(data interface{}){

}

func (this *stepObjm3)stopHandle(){

}































