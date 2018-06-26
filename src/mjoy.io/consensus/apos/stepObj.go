package apos

import (
	"sync"
	"time"
	"mjoy.io/common/types"
	"math/big"
	"fmt"
	"container/heap"
)

//step 1
type stepObj1 struct {
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
	stopCh chan interface{}
}

func makeStepObj1(timeEnd time.Duration , stopCh chan interface{})*stepObj1{
	s := new(stepObj1)
	s.timeEnd = timeEnd
	s.stopCh = stopCh
	return s
}

func (this *stepObj1)setCtx(ctx *stepCtx){
	this.ctx = ctx
}

func (this *stepObj1)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj1)timerHandle(){
	//new a M1 data
	m1 := new(M1)

	m1.Credential = this.ctx.getCredential()
	m1.Block = this.ctx.makeEmptyBlockForTest(m1.Credential)
	m1.Esig = this.ctx.esig(m1.Block.Hash())
	//fill struct members
	//todo: should using interface
	this.ctx.sendInner(m1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Out M1",COLOR_SHORT_RESET)

	go func(this *stepObj1){
		this.stopCh<-1
	}(this)
}

func (this *stepObj1)dataHandle(data interface{}){

}

func (this *stepObj1)stopHandle(){

}



//step 2
type stepObj2 struct {
	smallestLBr *M1 //this node regard smallestLBr as the smallest credential's block info
	lock sync.RWMutex
	ctx *stepCtx
	timeEnd time.Duration
	stopCh chan interface{}
}

func makeStepObj2(timeEnd time.Duration , stopCh chan interface{})*stepObj2{
	s := new(stepObj2)
	s.timeEnd = timeEnd
	s.stopCh = stopCh
	return s
}

func (this * stepObj2)setCtx(ctx *stepCtx){
	this.ctx = ctx
}

func (this *stepObj2)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj2)timerHandle(){
	m2 := new(M23)
	m2.Credential = this.ctx.getCredential()

	if this.smallestLBr == nil {
		m2.Hash = types.Hash{}
		logger.Error(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"Out M2 By a EmptyHash",COLOR_SHORT_RESET)
	}else{
		m2.Hash = this.smallestLBr.Block.Hash()
	}

	sigBytes := this.ctx.esig(m2.Hash)
	m2.Esig = append(m2.Esig , sigBytes...)
	this.ctx.sendInner(m2)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"Out M2",COLOR_SHORT_RESET)
	//turn to stop
	go func(this *stepObj2){
		this.stopCh<-1
	}(this)
}

func (this *stepObj2)dataHandle(data interface{}){
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"In M1",COLOR_SHORT_RESET)
	m1 := new(M1)
	m1 = data.(*M1)

	if m1 == nil{
		return
	}
	if this.smallestLBr == nil {
		this.smallestLBr = m1
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
	stopCh chan interface{}
}

func makeStepObj3(timeEnd time.Duration , stopCh chan interface{})*stepObj3{
	s := new(stepObj3)
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
	s.timeEnd = timeEnd
	s.stopCh = stopCh
	return s
}

func (this *stepObj3)setCtx(ctx *stepCtx){
	this.ctx = ctx
}

func (this *stepObj3)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj3)timerHandle(){
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

	m3 := new(M23)
	m3.Credential = this.ctx.getCredential()

	m3.Hash = v
	sigBytes := this.ctx.esig(m3.Hash)
	m3.Esig = append(m3.Esig , sigBytes...)
	this.ctx.sendInner(m3)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"Out M3 ",v.String(),COLOR_SHORT_RESET)
	go func(this *stepObj3){
		this.stopCh<-1
	}(this)
}

func (this *stepObj3)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()

	m2 := new(M23)
	m2 = data.(*M23)
	if m2 == nil{
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"In M2",m2.Hash.String(),COLOR_SHORT_RESET)
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
	stopCh chan interface{}
}

func makeStepObj4(timeEnd time.Duration , stopCh chan interface{})*stepObj4{
	s := new(stepObj4)
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
	s.timeEnd = timeEnd
	s.stopCh = stopCh
	return s
}

func (this *stepObj4)setCtx(ctx *stepCtx){
	this.ctx = ctx
}

func (this *stepObj4)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObj4)timerHandle(){
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
	v := types.Hash{}
	g := 0
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
		v = types.Hash{}
		g = 0
	}

	b := 0

	if g == 2 {
		b = 0
	}else{
		b = 1
	}

	//pack m4 Data
	m4 := new(MCommon)
	m4.Hash = v
	m4.B = uint(b)
	m4.Credential = this.ctx.getCredential()

	//b big.Int
	h := types.BytesToHash(big.NewInt(int64(m4.B)).Bytes())
	sigBytes := this.ctx.esig(h)
	m4.EsigB = append(m4.EsigB,sigBytes...)
	//v
	sigBytes = this.ctx.esig(m4.Hash)
	m4.EsigV = append(m4.EsigV , sigBytes...)
	this.ctx.sendInner(m4)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"Out M4",m4.Hash.String(),m4.B,COLOR_SHORT_RESET)

	go func(this *stepObj4){
		this.stopCh<-1
	}(this)

}

func (this *stepObj4)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()

	m3 := new(M23)
	m3 = data.(*M23)
	if m3 == nil{
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"In M3",m3.Hash.String(),COLOR_SHORT_RESET)
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
	stopCh chan interface{}
}


func makeStepObj567(timeEnd time.Duration , stopCh chan interface{})*stepObj567{
	s := new(stepObj567)
	s.allMxIndex = make(map[types.Hash]*binaryStatus)
	s.timeEnd = timeEnd
	s.stopCh = stopCh
	return s
}

func (this *stepObj567)setCtx(ctx *stepCtx){
	this.ctx = ctx
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




	mx := new(MCommon)
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
	h := types.BytesToHash(big.NewInt(int64(mx.B)).Bytes())
	sigBytes := this.ctx.esig(h)
	mx.EsigB = append(mx.EsigB , sigBytes...)
	//v
	sigBytes = this.ctx.esig(mx.Hash)
	mx.EsigV = append(mx.EsigV , sigBytes...)

	this.ctx.sendInner(mx)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"Out M",mx.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)

	go func(this *stepObj567){
		this.stopCh<-1
	}(this)
}


func (this *stepObj567)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()

	m6 := new(MCommon)
	m6 = data.(*MCommon)
	if m6 == nil{
		return
	}
	//add to IndexMap
	//logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M",m6.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step.IntVal.Int64(),"In M3 Hash:",m6.Hash.String(),"B:",m6.B,COLOR_SHORT_RESET)
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
	stopCh chan interface{}
}

func makeStepObjm3(timeEnd time.Duration , stopCh chan interface{})*stepObjm3{
	s := new(stepObjm3)
	s.timeEnd = timeEnd
	s.stopCh = stopCh
	return s
}

func (this *stepObjm3)setCtx(ctx *stepCtx){
	this.ctx = ctx
}

func (this *stepObjm3)getTTL()time.Duration{
	return this.timeEnd
}

func (this *stepObjm3)timerHandle(){
	this.lock.Lock()
	defer this.lock.Unlock()


	m3 := new(MCommon)
	//todo:should be H(Be)
	m3.Hash = types.Hash{}
	m3.B = 1
	m3.Credential = this.ctx.getCredential()
	h := types.BytesToHash(big.NewInt(int64(m3.B)).Bytes())
	sigBytes := this.ctx.esig(h)
	m3.EsigB = append(m3.EsigB , sigBytes...)
	//v
	sigBytes = this.ctx.esig(m3.Hash)
	m3.EsigV = append(m3.EsigV , sigBytes...)
	this.ctx.sendInner(m3)

	go func(this *stepObjm3){
		this.stopCh<-1
	}(this)
}


func (this *stepObjm3)dataHandle(data interface{}){

}

func (this *stepObjm3)stopHandle(){

}































