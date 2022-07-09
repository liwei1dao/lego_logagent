package sender

import (
	"fmt"
	"lego_logagent/modules/logagent/core"
	"sync"
	"sync/atomic"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/blockcache"
)

type Sender struct {
	Runner  core.IRunner
	sender  core.ISender
	options core.ISenderOptions
	Procs   int
	Cnt     int64
	Wg      *sync.WaitGroup
	cache blockcache.ISys
}

func (this *Sender) GetRunner() core.IRunner {
	return this.Runner
}
func (this *Sender) GetType() string {
	return this.options.GetType()
}

func (this *Sender) GetMetaerData() (meta core.IMetaerData) {
	return
}
func (this *Sender) Init(rer core.IRunner, ser core.ISender, options core.ISenderOptions) (err error) {
	defer log.Infof("NewSender options:%+v", options)
	this.Runner = rer
	this.sender = ser
	this.options = options
	this.Wg = new(sync.WaitGroup)
	this.Procs = this.Runner.MaxProcs()
	this.Cnt = 0
	return
}
func (this *Sender) Start() (err error) {
	if this.Procs < 1 {
		this.Procs = 1
	}

	if pipe, ok := this.Runner.SenderPipe(this.GetType()); !ok {
		err = fmt.Errorf("no found SenderPope:%s", this.GetType())
		return
	} else {
		for i := 0; i < this.Procs; i++ {
			this.Wg.Add(1)
			go this.sender.Run(i, pipe)
		}
	}
	return
}

func (this *Sender) Run(pipeId int, pipe <-chan core.ICollData, params ...interface{}) {
	defer this.Wg.Done()
	for v := range pipe {
		this.sender.Send(pipeId, v, params...)
	}
}

func (this *Sender) Send(pipeId int, data core.ICollData, params ...interface{}) {
	atomic.AddInt64(&this.Cnt, 1)
}

//关闭
func (this *Sender) Close() (err error) {
	this.Wg.Wait()
	log.Debugf("Sender Close Succ")
	return
}
func (this *Sender) ReadCnt() int64 {
	return atomic.LoadInt64(&this.Cnt)
}
func (this *Sender) ReadAnResetCnt() int64 {
	return atomic.SwapInt64(&this.Cnt, 0)
}
