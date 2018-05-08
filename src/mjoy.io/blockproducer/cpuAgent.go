
package blockproducer


import (
	"sync"

	"sync/atomic"
	"mjoy.io/consensus"

	"time"
	"math/rand"
)

type CpuAgent struct {
	mu sync.Mutex

	workCh        chan *Work
	stop          chan struct{}
	quitCurrentOp chan struct{}
	returnCh      chan<- *Result

	chain  consensus.ChainReader
	engine consensus.Engine

	isProducing int32 // isProducing indicates whether the agent is currently producing
}

func NewCpuAgent(chain consensus.ChainReader, engine consensus.Engine) *CpuAgent {
	cpuBlockproducer := &CpuAgent{
		chain:  chain,
		engine: engine,
		stop:   make(chan struct{}, 1),
		workCh: make(chan *Work, 1),
	}
	return cpuBlockproducer
}

func (self *CpuAgent) Work() chan<- *Work            { return self.workCh }

func (self *CpuAgent) SetReturnCh(ch chan<- *Result) { self.returnCh = ch }

func (self *CpuAgent) Stop() {
	if !atomic.CompareAndSwapInt32(&self.isProducing, 1, 0) {
		return // agent already stopped
	}
	self.stop <- struct{}{}
done:
	// Empty work channel
	for {
		select {
		case <-self.workCh:
		default:
			break done
		}
	}
}

func (self *CpuAgent) Start() {

	if !atomic.CompareAndSwapInt32(&self.isProducing, 0, 1) {
		return // agent already started
	}
	go self.update()
}

func (self *CpuAgent) update() {
out:
	for {
		select {
		case work := <-self.workCh:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
			}
			self.quitCurrentOp = make(chan struct{})
			go self.blockproducer(work, self.quitCurrentOp)
			self.mu.Unlock()
		case <-self.stop:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
				self.quitCurrentOp = nil
			}
			self.mu.Unlock()
			break out
		}
	}
}

func (self *CpuAgent) blockproducer(work *Work, stop <-chan struct{}) {

	if result, err := self.engine.Seal(self.chain, work.Block, stop); result != nil {

		logger.Infof("Successfully sealed new block number: %d  hash:0x%x\n" , result.Number() , result.Hash())
		//fmt.Println("ProduceBlock: num:" , result.Number().String(),"  Hash:",result.Hash().String())
		time.Sleep(time.Duration(rand.Intn(20))*time.Second)
		self.returnCh <- &Result{work, result}
		//fmt.Printf("!!!!!Return Produce work......")

	} else {
		if err != nil {
			logger.Warn("Block sealing failed", "err", err)
		}
		self.returnCh <- nil
	}
}


func (self *CpuAgent) GetHashRate() int64 {
	return 0
}

