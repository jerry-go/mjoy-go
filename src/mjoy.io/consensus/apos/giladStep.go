package apos

import (
	"mjoy.io/common/types"
	"sync"
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
	v.emptyHash = v.ctx.getGiladEmptyHash(uint64(ctx.getRound()))

	return v
}

func (this *VoteObj)stop(){
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
	v := new(VoteData)
	v.Round = r
	v.Step = s
	v.Value = hash

	this.msgChan <- v
}

func (this *VoteObj)run(){

	for{
		select {
		case data := <-this.msgChan:
			//data deal
			if this.isSendSameStepData(data.Step) == false{
				this.dataDeal(data)
			}
		case <-this.exit:
			logger.Debug("VoteObj exit")

		}
	}
}




func (this *VoteObj)CommitteeVote(data *VoteData){

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
	this.lock.Lock()
	defer this.lock.Unlock()

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

	}else if step == StepFinal{
		//when get StepFinal,return the hash
		this.ctx.writeRet(data)
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
				//return ,just write the Ret
				this.ctx.writeRet(data)


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

				this.ctx.writeRet(data)
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







