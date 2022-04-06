package reader

import (
	"lego_logagent/modules/logagent/core"
	"sync"
	"sync/atomic"

	"github.com/liwei1dao/lego/sys/log"
)

type Reader struct {
	Runner     core.IRunner
	reader     core.IReader
	meta       core.IMetaerData
	options    core.IReaderOptions
	Procs      int
	TaskPipe   chan core.IMetaerNodeData
	State      int32 //任务组状态 0 未开始 1 任务执行中 2任务组完执行完毕
	RuntaskNum int32 //当前运行任务数
	Wg         *sync.WaitGroup
}

func (this *Reader) GetRunner() core.IRunner {
	return this.Runner
}
func (this *Reader) GetType() string {
	return this.options.GetType()
}

func (this *Reader) GetEncoding() core.Encoding {
	return this.options.GetEncoding()
}

func (this *Reader) Init(runner core.IRunner, reader core.IReader, meta core.IMetaerData, options core.IReaderOptions) (err error) {
	defer log.Infof("NewReader options:%+v", options)
	this.Runner = runner
	this.meta = meta
	this.options = options
	if meta != nil { //有些采集器不需要保存记录
		err = this.Runner.Metaer().Read(meta)
	}
	return
}

func (this *Reader) Start() (err error) {
	if this.Procs < 1 {
		this.Procs = 1
	}
	this.Wg.Add(this.Procs)
	for i := 0; i < this.Procs; i++ {
		go this.run()
	}
	return
}

func (this *Reader) run() {
	defer this.Wg.Done()
	for v := range this.TaskPipe {
		if err := this.reader.Read(v); err != nil {
			log.Errorf("err:%v", err)
		}
		if atomic.CompareAndSwapInt32(&this.State, 1, 2) { //最后一个任务已经完成
			this.Wg.Wait()
			atomic.StoreInt32(&this.State, 0)
		}
	}
}

func (this *Reader) Read(task core.IMetaerData) (err error) {

	return
}

///外部调度器 驱动执行  此接口 不可阻塞
func (this *Reader) Drive() (err error) {
	if !atomic.CompareAndSwapInt32(&this.State, 0, 1) {
		log.Debugf("Reader is collectioning runtaskNum:%d", atomic.LoadInt32(&this.RuntaskNum))
		err = core.Error_RunnerTaskExecuting
		return
	}
	return
}

//关闭 只允许 Runner 对象调用
func (this *Reader) Close() (err error) {
	// err = this.SyncMeta()此处无需加同步数据代码 Runner 在关闭最后 会将通过 Metaer 同步所有的 MetaeData
	log.Debugf("Reader Start Close!")
	return
}

func (this *Reader) SyncMeta() (err error) {
	return this.Runner.Metaer().Write()
}

func (this *Reader) Input() chan<- core.ICollData {
	return this.Runner.ReaderPipe()
}
