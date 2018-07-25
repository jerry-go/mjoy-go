package apos

import (
	"mjoy.io/common/types"
	"sync"
	"time"
	"fmt"
)

const(

	StepBp          = STEP_BP
	StepReduction1  = STEP_REDUCTION_1
	StepReduction2  = STEP_REDUCTION_2
	StepFinal       = STEP_FINAL

)

var(
	TimeOut types.Hash = types.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)


type VoteData struct {
	Round uint64
	Step uint64
	Value types.Hash
}

type  VoteObj struct {
	lock sync.RWMutex
	ctx *stepCtx
	msgChan chan *VoteData
	SendStatus map[uint64]*VoteData

	emptyHash types.Hash    //H(Empty(round H(ctx.last_block)))
	bbaBlockHash types.Hash //the block hash set by the reduction last step

	isBbaIsOk bool
	bbaStayList []*VoteData
	listLock sync.RWMutex
	exit chan interface{}
}

func makeVoteObj(ctx *stepCtx)*VoteObj{
	v := new(VoteObj)
	v.ctx = ctx
	v.SendStatus = make(map[uint64]*VoteData)
	v.msgChan = make(chan *VoteData , 1000)
	v.exit = make(chan interface{}, 1)
	v.emptyHash = v.ctx.getGiladEmptyHash(uint64(ctx.getRound()))
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"This Round EmptyHash:" , v.emptyHash.Hex(),COLOR_SHORT_RESET)
	logger.Debug("***********Print StepBp:" , StepBp)
	logger.Debug("***********Print Reduction1:" , StepReduction1)
	logger.Debug("***********Print Reduction2:" , StepReduction2)
	logger.Debug("***********Print StepFinal :" , StepFinal)
	return v
}

func (this *VoteObj)stop(){
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "In VoteObj stop()....", COLOR_SHORT_RESET)
	this.exit<-1
}

func (this *VoteObj)isBbaEmpty()bool{
	return this.isBbaIsOk
}


func (this *VoteObj)setBbaBlockHash(bHash types.Hash){
	this.listLock.Lock()
	defer this.listLock.Unlock()


	if this.isBbaEmpty(){

		copy(this.bbaBlockHash[:] , bHash[:])
		//set isBbaIsOk
		this.isBbaIsOk = true
		//need clear bbastay list
		for _ , v := range this.bbaStayList{
			copy(v.Value[:] , this.bbaBlockHash[:])
			this.CommitteeVote(v)
		}
		//a new one ,old one need GC
		this.bbaStayList = []*VoteData{}

	}
}

func (this *VoteObj)addStayBbaData(data *VoteData){
	this.listLock.Lock()
	defer this.listLock.Unlock()

	this.bbaStayList = append(this.bbaStayList , data)
}

//return true:we have send a data with same step ,
//return false:we do not send a data with same step,can send a data
func (this *VoteObj)isSendSameStepData(step uint64)bool{
	this.lock.RLock()
	defer this.lock.RUnlock()

	if _ , ok := this.SendStatus[step];ok {
		return true
	}
	return false
}

//send the data ,and set the data to mark map
func (this *VoteObj)markSendData(data *VoteData){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.SendStatus[data.Step] = data
}

func (this *VoteObj)SendVoteData(r,s uint64 , hash types.Hash){

	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "SendVoteData Step:" , s , "   Hash:",hash.Hex(), COLOR_SHORT_RESET)
	v := new(VoteData)
	v.Round = r
	v.Step = s
	v.Value = hash

	this.msgChan <- v
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "SendVoteData Ok" , COLOR_SHORT_RESET)
}

func (this *VoteObj)run(){
	round := this.ctx.getRound()
	tick := time.Tick(3*time.Second)
	for{
		select {
		case data := <-this.msgChan:

			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "VoteObj RecvData Step:" , data.Step , "   VoteHash:" , data.Value.Hex() , COLOR_SHORT_RESET)
			//data deal

			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "isSendSameStepData == false , call dataDeal"  , COLOR_SHORT_RESET)
			this.dataDeal(data)

		case <-this.exit:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX ,"VoteObj exit" , COLOR_SHORT_RESET)
			return
		case <-tick:
			fmt.Println("Round:",round , "  VoteObj is running............" )

		}
	}
}




func (this *VoteObj)CommitteeVote(data *VoteData){

	if this.isSendSameStepData(data.Step){
		return
	}
	cret := this.ctx.getCredentialByStep(uint64(data.Step))
	if cret == nil {
		return
	}

	//todo :need pack ba msg

	msgBa := newByzantineAgreementStar()
	//hash
	msgBa.Hash = data.Value
	//Credential
	msgBa.Credential = cret

	//Esig
	msgBa.Esig.round = msgBa.Credential.Round
	msgBa.Esig.step = msgBa.Credential.Step
	msgBa.Esig.val = make([]byte , 0)
	msgBa.Esig.val = append(msgBa.Esig.val , msgBa.Hash.Bytes()...)

	err := this.ctx.esig(msgBa.Esig)
	if err != nil {
		logger.Error("CommitteeVote Esig Err:" , err.Error())
		return
	}


	if cret.votes > 0{
		this.markSendData(data)
		this.ctx.sendInner(msgBa)
	}
}
//this function just using in
func (this *VoteObj)safetyBbaCommitteeVote(data *VoteData){
	if this.isBbaEmpty(){
		this.addStayBbaData(data)
	}else{
		this.CommitteeVote(data)
	}
}

func (this *VoteObj)dataDeal(data *VoteData){

	step := data.Step
	if step == StepReduction1{
		//check the hblock1 is Timeout
		if timeout := data.Value.Equal(&TimeOut);timeout {
			copy(data.Value[:], this.emptyHash[:])
		}
		data.Step = StepReduction2
		this.CommitteeVote(data)

	}else if step == StepReduction2 {
		//check the hblock2 is TimeOut
		if timeout := data.Value.Equal(&TimeOut);timeout{
			copy(data.Value[:] , this.emptyHash[:])
		}
		this.setBbaBlockHash(data.Value)
		//set the bba block hash
		//copy(this.bbaBlockHash[:] , data.Value[:])
		//send the bba first step data
		data.Step = 1
		//this is the bba first step
		this.CommitteeVote(data)
		//equal return the reduction result
		this.ctx.setReductionResult(data.Value)
	}else if step == StepFinal{
		//when get StepFinal,return the hash
		this.ctx.setFinalResult(data.Value)
	}else{
		//common step:BBA Step
		index := step % 3
		if index == 0 {
			index = 3
		}


		switch index {
		case 1:
			if timeout := data.Value.Equal(&TimeOut);timeout{
				//set data.value to the bbaBlockHash
				copy(data.Value[:] , this.bbaBlockHash[:])
				data.Step += 1
				this.safetyBbaCommitteeVote(data)

			}else if empty := data.Value.Equal(&this.emptyHash);!empty {
				for i:= step + 1;i <= step + 3; i++{
					dataNew := new(VoteData)
					dataNew.Round = data.Round
					dataNew.Step = i
					copy(dataNew.Value[:] , data.Value[:])
					this.safetyBbaCommitteeVote(dataNew)
				}

				if step == 1 {
					data.Step = StepFinal
					this.safetyBbaCommitteeVote(data)
				}
				//equal return bba
				this.ctx.setBbaResult(data.Value)


			}
		case 2:
			if timeout := data.Value.Equal(&TimeOut);timeout{
				//ste data.value to the emptyhash
				copy(data.Value[:] , this.emptyHash[:])
				data.Step += 1
				this.safetyBbaCommitteeVote(data)

			}else if empty := data.Value.Equal(&this.emptyHash);empty{
				for i := step + 1; i <= step + 3; i++ {
					dataNew := new(VoteData)
					dataNew.Round = data.Round
					dataNew.Step = i
					copy(dataNew.Value[:] , data.Value[:])
					this.safetyBbaCommitteeVote(dataNew)
				}
				//equal return bba
				this.ctx.setBbaResult(data.Value)
			}
		case 3:
			if timeout := data.Value.Equal(&TimeOut);timeout{

				if this.ctx.commonCoin(data.Round , data.Step , 100) == 0 {
					copy(data.Value[:] , this.bbaBlockHash[:])
				}else{
					copy(data.Value[:] , this.emptyHash[:])
				}
			}
			data.Step += 1
			this.safetyBbaCommitteeVote(data)
		}
	}

}







