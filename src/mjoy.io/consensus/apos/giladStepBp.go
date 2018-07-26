package apos

import (
	"sync"
	"time"
	"sort"
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
)

type BpWithPriority struct {
	j int   //the priofity
	bp *BlockProposal
}

type BpWithPriorityHeap []*BpWithPriority

func (h BpWithPriorityHeap)Len() int            {return len(h)}
func (h BpWithPriorityHeap)Less(i , j int)bool  {return h[i].j > h[j].j}
func (h BpWithPriorityHeap)Swap(i , j int)      {h[i],h[j] = h[j],h[i]}

func (h *BpWithPriorityHeap)Push(x interface{}){
	*h = append(*h , x.(*BpWithPriority))
}

func (h *BpWithPriorityHeap)Pop()interface{}{
	old := *h
	n := len(old)
	x := old[n -1]
	*h = old[0:n-1]
	return x
}


type BpObj struct {
	lock    sync.RWMutex
	BpHeap  BpWithPriorityHeap
	existMap map[types.Hash]*BlockProposal	//todo:here should be change to *bp
	msgChan chan *BlockProposal
	exit    chan interface{}
	ctx     *stepCtx
	priorityBp *BlockProposal
	nothingTodo bool
}

func makeBpObj(ctx *stepCtx) *BpObj {
	s := new(BpObj)
	s.ctx = ctx
	s.BpHeap = make(BpWithPriorityHeap , 0)
	s.msgChan = make(chan *BlockProposal)
	s.existMap = make(map[types.Hash]*BlockProposal)
	s.exit = make(chan interface{},2)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX, "makeBpObj" , COLOR_SHORT_RESET)
	return s
}

func (this *BpObj)isExistBlock(blockHash types.Hash)bool{
	this.lock.RLock()
	defer this.lock.RUnlock()

	if _,ok := this.existMap[blockHash];ok{
		return true
	}
	return false
}

func (this *BpObj)addExistBlock(bp *BlockProposal){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.existMap[bp.Block.Hash()] = bp
}

func (this *BpObj)getExistBlock(blockHash types.Hash) *block.Block {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if bp, ok := this.existMap[blockHash];ok{
		return bp.Block
	}
	return nil
}
func (this *BpObj)CommitteeVote(data *VoteData){

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
		this.ctx.sendInner(msgBa)
	}
}

func (this *BpObj)makeBlock(){
	bp := newBlockProposal()
	bp.Credential = this.ctx.getCredentialByStep(StepBp)
	if nil == bp.Credential {
		logger.Warn("makeBlock getCredentialByStep--->nil")
		return
	}
	//bcd := &block.ConsensusData{}
	//bcd.Id = ConsensusDataId
	//bcd.Para = bp.Credential.Signature.toBytes()
	bcd := this.ctx.makeBlockConsensusData(bp)

	bp.Block = this.ctx.getProducerNewBlock(bcd)

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = StepBp
	bp.Esig.val = make([]byte , 0)
	h := bp.Block.Hash()
	bp.Esig.val = append(bp.Esig.val , h[:]...)
	err := this.ctx.esig(bp.Esig)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//delay here ?

	this.ctx.sendInner(bp)

	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX, "***[A]Out M1 CreHash:", bp.Credential.Signature.Hash().String(), " BlockHash", bp.Block.B_header.Hash().String(), COLOR_SHORT_RESET)
}
func (this *BpObj)stop(){
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "Call BpObj Exit....:"  , COLOR_SHORT_RESET)
	this.exit <- 1
}
func (this *BpObj) run() {
	rd := this.ctx.getRound()
	//make block and send out
	go this.makeBlock()
	tProposer := int(Config().tProposer)
	timer := time.Tick(20 * time.Second)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "#########In BpObj-",rd , COLOR_SHORT_RESET)
	defer func() {
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "#######Out BpObj",rd , COLOR_SHORT_RESET)
	}()
	for {
		select {
		case <-this.exit:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "BpObj Exit....:return"  , COLOR_SHORT_RESET)
			return
		case <-timer:
			if this.nothingTodo {
				continue
			}
			value := types.Hash{}
			if this.BpHeap.Len() == 0 {
				//output empty hash
				value = this.ctx.getEmptyBlockHash()
			}else{
				sort.Sort(&this.BpHeap)
				value = this.BpHeap[0].bp.Block.Hash()
				//make reduction input data
			}
			vd := new(VoteData)
			vd.Round = this.ctx.getRound()
			vd.Step = StepReduction1
			vd.Value = value
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "BpObj timeOut dataOutput hash:" , vd.Value.Hex() , COLOR_SHORT_RESET)
			this.CommitteeVote(vd)

			this.ctx.setBpResult(value)
			this.ctx.startVoteTimer(int(Config().delayStep))
			this.nothingTodo = true
		case bp := <-this.msgChan:
				if this.nothingTodo {
					continue
				}
				//logic do
				//verify the block
				if !this.ctx.verifyBlock(bp.Block){
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "!this.ctx.verifyBlock(bp.Block) Wrong:hash:" , bp.Block.Hash().Hex(), COLOR_SHORT_RESET)
					continue
				}
				//check is exist a same block
				if this.isExistBlock(bp.Block.Hash()) {
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "this.isExistBlock(bp.Block.Hash()) Wrong:hash:" , bp.Block.Hash().Hex(), COLOR_SHORT_RESET)
					continue
				}

				//check the node has the right to produce a block
				pri := bp.Credential.votes

				bpp := new(BpWithPriority)
				bpp.j = int(pri)
				bpp.bp = bp

				this.BpHeap.Push(bpp)
				this.addExistBlock(bp)
				if this.priorityBp == nil {
					this.ctx.propagateMsg(bp)
				} else if pri > this.priorityBp.Credential.votes {
					this.priorityBp = bp
					this.ctx.propagateMsg(bp)
				}

				if this.BpHeap.Len() > tProposer {
					sort.Sort(&this.BpHeap)
					//get the bigger one
					x := this.BpHeap[0]
					_ = x


					vd := new(VoteData)
					vd.Round = x.bp.Credential.Round
					vd.Step = StepReduction1
					vd.Value = x.bp.Block.Hash()
						logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX , "BpObj >tProposer dataOutput hash:" , vd.Value.Hex() , COLOR_SHORT_RESET)
					this.CommitteeVote(vd)
					this.ctx.startVoteTimer(int(Config().delayStep))
					this.ctx.setBpResult(x.bp.Block.Hash())
					//todo:inform the reduction

					this.nothingTodo = true
				}
		}
	}
}

func (this *BpObj) sendMsg(bp *BlockProposal) {
	this.msgChan <- bp
}
