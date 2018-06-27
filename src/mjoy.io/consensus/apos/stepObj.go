package apos

import (
	"sync"
	"time"
	"mjoy.io/common/types"
	"math/big"
	"fmt"
	"container/heap"
)


var (
	LessTimeDelayFlag bool = false  //let step spend less time to deal msg received
	LessTimeDelayCnt int = 5
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
	m1 := newBlockProposal()


	m1.Credential = this.ctx.getCredential()
	m1.Block = this.ctx.makeEmptyBlockForTest(m1.Credential)

	m1.Esig.round = m1.Credential.Round
	m1.Esig.step = 1
	m1.Esig.val = make([]byte,0)
	h := m1.Block.Hash()
	m1.Esig.val = append(m1.Esig.val , h[:]...)
	R,S,V := this.ctx.esig(m1.Block.Hash())
	m1.Esig.Signature.R = new(types.BigInt)
	m1.Esig.Signature.R.IntVal = *R

	m1.Esig.Signature.S = new(types.BigInt)
	m1.Esig.Signature.S.IntVal = *S

	m1.Esig.Signature.V = new(types.BigInt)
	m1.Esig.Signature.V.IntVal = *V

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
	smallestLBr *BlockProposal //this node regard smallestLBr as the smallest credential's block info
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
	m2 := newGradedConsensus()
	m2.Credential = this.ctx.getCredential()
	if m2.Credential == nil {
		panic("timeHandle Credential Wrong.........")
	}
	fmt.Println("m2.Credential:" , *m2.Credential)
	if this.smallestLBr == nil {
		m2.Hash = types.Hash{}
		logger.Error(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"Out M2 By a EmptyHash",COLOR_SHORT_RESET)
	}else{
		m2.Hash = this.smallestLBr.Block.Hash()
	}

	m2.Esig.round = m2.Credential.Round
	m2.Esig.step = m2.Credential.Step
	m2.Esig.val = make([]byte , 0)
	m2.Esig.val = append(m2.Esig.val , m2.Hash[:]...)

	R,S,V := this.ctx.esig(m2.Hash)

	m2.Esig.Signature.R = new(types.BigInt)
	m2.Esig.Signature.R.IntVal = *R

	m2.Esig.Signature.S = new(types.BigInt)
	m2.Esig.Signature.S.IntVal = *S

	m2.Esig.Signature.V = new(types.BigInt)
	m2.Esig.Signature.V.IntVal = *V

	this.ctx.sendInner(m2)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"Out M2",COLOR_SHORT_RESET)
	//turn to stop
	go func(this *stepObj2){
		this.stopCh<-1
	}(this)
}

func (this *stepObj2)dataHandle(data interface{}){
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"In M1",COLOR_SHORT_RESET)
	m1 := new(BlockProposal)
	m1 = data.(*BlockProposal)

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

	m3 := newGradedConsensus()
	m3.Credential = this.ctx.getCredential()

	m3.Hash = v
	m3.Esig.round = m3.Credential.Round
	m3.Esig.step = m3.Credential.Step
	m3.Esig.val = make([]byte , 0)
	m3.Esig.val = append(m3.Esig.val , m3.Hash[:]...)

	R,S,V := this.ctx.esig(m3.Hash)

	m3.Esig.R = new(types.BigInt)
	m3.Esig.R.IntVal = *R

	m3.Esig.S = new(types.BigInt)
	m3.Esig.S.IntVal = *S

	m3.Esig.V = new(types.BigInt)
	m3.Esig.V.IntVal = *V

	this.ctx.sendInner(m3)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"Out M3 ",v.String(),COLOR_SHORT_RESET)
	go func(this *stepObj3){
		this.stopCh<-1
	}(this)
}

func (this *stepObj3)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()

	m2 := newGradedConsensus()
	m2 = data.(*GradedConsensus)
	if m2 == nil{
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"In M2",m2.Hash.String(),COLOR_SHORT_RESET)
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
	m4 := newBinaryByzantineAgreement()
	m4.Hash = v
	m4.B = uint(b)
	m4.Credential = this.ctx.getCredential()

	//b big.Int
	R,S,V := this.ctx.esig(types.BytesToHash(big.NewInt(int64(m4.B)).Bytes()))
	m4.EsigB.R = new(types.BigInt)
	m4.EsigB.R.IntVal = *R

	m4.EsigB.S = new(types.BigInt)
	m4.EsigB.R.IntVal = *S

	m4.EsigB.V = new(types.BigInt)
	m4.EsigB.V.IntVal = *V


	//v
	R,S,V = this.ctx.esig(m4.Hash)
	m4.EsigV.R = new(types.BigInt)
	m4.EsigV.R.IntVal = *R

	m4.EsigV.S = new(types.BigInt)
	m4.EsigV.S.IntVal = *S

	m4.EsigV.V = new(types.BigInt)
	m4.EsigV.V.IntVal = *V


	this.ctx.sendInner(m4)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"Out M4",m4.Hash.String(),m4.B,COLOR_SHORT_RESET)

	go func(this *stepObj4){
		this.stopCh<-1
	}(this)

}

func (this *stepObj4)dataHandle(data interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()
	m3 := newGradedConsensus()
	m3 = data.(*GradedConsensus)
	if m3 == nil{
		return
	}
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"In M3",m3.Hash.String(),COLOR_SHORT_RESET)
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
	R,S,V := this.ctx.esig(types.BytesToHash(big.NewInt(int64(mx.B)).Bytes()))
	mx.EsigB.R = new(types.BigInt)
	mx.EsigB.R.IntVal = *R

	mx.EsigB.S = new(types.BigInt)
	mx.EsigB.R.IntVal = *S

	mx.EsigB.V = new(types.BigInt)
	mx.EsigB.V.IntVal = *V


	//v
	R,S,V = this.ctx.esig(mx.Hash)
	mx.EsigV.R = new(types.BigInt)
	mx.EsigV.R.IntVal = *R

	mx.EsigV.S = new(types.BigInt)
	mx.EsigV.S.IntVal = *S

	mx.EsigV.V = new(types.BigInt)
	mx.EsigV.V.IntVal = *V


	this.ctx.sendInner(mx)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"Out M",mx.Credential.Step,COLOR_SHORT_RESET)

	go func(this *stepObj567){
		this.stopCh<-1
	}(this)
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
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.ctx.getCredential().Step,"In M3 Hash:",m6.Hash.String(),"B:",m6.B,COLOR_SHORT_RESET)
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

	m3 := newBinaryByzantineAgreement()
	//todo:should be H(Be)
	m3.Hash = types.Hash{}
	m3.B = 1
	m3.Credential = this.ctx.getCredential()

	//b big.Int
	R,S,V := this.ctx.esig(types.BytesToHash(big.NewInt(int64(m3.B)).Bytes()))
	m3.EsigB.R = new(types.BigInt)
	m3.EsigB.R.IntVal = *R

	m3.EsigB.S = new(types.BigInt)
	m3.EsigB.R.IntVal = *S

	m3.EsigB.V = new(types.BigInt)
	m3.EsigB.V.IntVal = *V


	//v
	R,S,V = this.ctx.esig(m3.Hash)
	m3.EsigV.R = new(types.BigInt)
	m3.EsigV.R.IntVal = *R

	m3.EsigV.S = new(types.BigInt)
	m3.EsigV.S.IntVal = *S

	m3.EsigV.V = new(types.BigInt)
	m3.EsigV.V.IntVal = *V





	this.ctx.sendInner(m3)

	go func(this *stepObjm3){
		this.stopCh<-1
	}(this)
}


func (this *stepObjm3)dataHandle(data interface{}){

}

func (this *stepObjm3)stopHandle(){

}































