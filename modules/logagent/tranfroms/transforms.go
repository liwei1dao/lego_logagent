package tranfroms

import (
	"fmt"
	"lego_logagent/modules/logagent/core"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

type Transforms struct {
	Runner     core.IRunner
	transforms core.ITransforms
	options    core.ITransformsOptions
	index      int
	wg         *sync.WaitGroup
	Procs      int
}

func (this *Transforms) GetRunner() core.IRunner {
	return this.Runner
}

func (this *Transforms) Init(index int, runner core.IRunner, transforms core.ITransforms, options core.ITransformsOptions) (err error) {
	defer log.Infof("Transforms Init options:%+v", options)
	this.Runner = runner
	this.transforms = transforms
	this.options = options
	this.index = index
	this.Procs = this.Runner.MaxProcs()
	this.wg = new(sync.WaitGroup)
	return
}

func (this *Transforms) Start() (err error) {
	if this.Procs < 1 {
		this.Procs = 1
	}

	for i := 0; i < this.Procs; i++ {
		if pipe, ok := this.Runner.TransformsPipe(this.index); !ok {
			err = fmt.Errorf("no found TransformsPipe:%d", this.index)
			return
		} else {
			for i := 0; i < this.Procs; i++ {
				this.wg.Add(1)
				go this.run(i, pipe)
			}
		}
	}
	return
}

func (this *Transforms) run(pipeId int, pipe <-chan core.ICollData) {
	defer this.wg.Done()
	for v := range pipe {
		this.transforms.Trans(v)
	}
}

func (this *Transforms) Trans(data core.ICollData) {
	this.Runner.Push_NextTransformsPipe(this.index+1, data)
}

//关闭
func (this *Transforms) Close() (err error) {
	this.wg.Wait()
	log.Debugf("Transforms Close Succ")
	return
}
