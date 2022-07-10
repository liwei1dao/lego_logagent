package loganget

import (
	"fmt"
	"lego_logagent/modules/logagent/core"
	"lego_logagent/modules/logagent/metaer"
	"lego_logagent/modules/logagent/parser"
	"lego_logagent/modules/logagent/reader"
	"lego_logagent/modules/logagent/sender"
	"lego_logagent/modules/logagent/tranfroms"
	"sync/atomic"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/log"
)

type Runner struct {
	conf           *RunnerConfig
	log            log.ILog
	readerPope     chan core.ICollData
	parserPope     chan core.ICollData
	transformsPope []chan core.ICollData
	sendersPope    map[string]chan core.ICollData
	metaer         core.IMetaer
	reader         core.IReader
	parser         core.IParser
	transforms     []core.ITransforms
	senders        []core.ISender
	state          int32
	readerCnt      int64
}

func (this *Runner) Name() string {
	return this.conf.Name
}
func (this *Runner) MaxProcs() int {
	return this.conf.MaxProcs
}
func (this *Runner) MaxMessageSzie() uint64 {
	return this.conf.MaxMessageSzie
}
func (this *Runner) Metaer() core.IMetaer {
	return this.metaer
}
func (this *Runner) State() core.RunnerState {
	return core.RunnerState(atomic.LoadInt32(&this.state))
}
func (this *Runner) ReaderPipe() chan<- core.ICollData {
	return this.readerPope
}
func (this *Runner) Push_ParserPipe(data core.ICollData) {
	this.parserPope <- data
}
func (this *Runner) ParserPipe() <-chan core.ICollData {
	return this.parserPope
}
func (this *Runner) Push_TransformsPipe(bucket core.ICollData) {
	if len(this.transformsPope) > 0 {
		this.transformsPope[0] <- bucket
	} else {
		this.Push_SenderPipe(bucket)
	}
}
func (this *Runner) TransformsPipe(index int) (pipe <-chan core.ICollData, ok bool) {
	if len(this.transformsPope) > index {
		pipe = this.transformsPope[index]
		ok = true
	} else {
		ok = false
	}
	return
}
func (this *Runner) Push_NextTransformsPipe(index int, bucket core.ICollData) {
	if index >= len(this.transformsPope) {
		this.Push_SenderPipe(bucket)
	} else {
		this.transformsPope[index] <- bucket
	}
}
func (this *Runner) Push_SenderPipe(bucket core.ICollData) {
	for _, v := range this.sendersPope {
		v <- bucket
	}
}
func (this *Runner) SenderPipe(stype string) (pipe <-chan core.ICollData, ok bool) {
	pipe, ok = this.sendersPope[stype]
	return
}

func (this *Runner) Init() (err error) {
	atomic.StoreInt32(&this.state, int32(core.Runner_Initing))
	//初始化元数据管理
	if this.metaer, err = metaer.NewMetaer(this); err != nil {
		log.Errorf("NewRunner NewMetaer fail:%v", err)
		return
	}
	this.readerPope = make(chan core.ICollData)
	if this.reader, err = reader.NewReader(this, this.conf.ReaderConfig); err != nil {
		log.Errorf("NewRunner NewReader fail:%v", err)
		return
	}
	this.parserPope = make(chan core.ICollData)
	if this.parser, err = parser.NewParser(this, this.conf.ParserConf); err != nil {
		log.Errorf("NewRunner NewParser fail:%v", err)
		return
	}
	if this.transforms, err = tranfroms.NewTransforms(this, this.conf.TranfromssConfig); err != nil {
		log.Errorf("NewRunner NewTransforms fail:%v", err)
		return
	}
	for i, _ := range this.transforms {
		this.transformsPope[i] = make(chan core.ICollData)
	}
	if this.senders, err = sender.NewSender(this, this.conf.SendersConfig); err != nil {
		log.Errorf("NewRunner NewSender fail:%v", err)
		return
	}
	for _, v := range this.senders {
		this.sendersPope[v.GetType()] = make(chan core.ICollData)
	}
	return
}

func (this *Runner) Start() (err error) {
	atomic.StoreInt32(&this.state, int32(core.Runner_Starting))
	for i, v := range this.senders {
		if err = v.Start(); err != nil {
			log.Errorf("启动Runner Senders:%v fail:%v", i, err)
			return
		}
	}
	for i, v := range this.transforms {
		if err = v.Start(); err != nil {
			log.Errorf("启动Runner Transforms:%v fail:%v", i, err)
			return
		}
	}
	if err = this.parser.Start(); err != nil {
		log.Errorf("启动Runner Parser fail:%v", err)
		return
	}
	if err = this.reader.Start(); err != nil {
		log.Errorf("启动Runner Reader fail:%v", err)
		return
	}
	log.Infof("启动Runner:%s", this.conf.Name)
	atomic.StoreInt32(&this.state, int32(core.Runner_Runing))
	if this.conf.CronSql != "" {
		cron.AddFunc(this.conf.CronSql, func() {
			if err := this.Drive(); err != nil {
				log.Errorf("err:%v", err)
			}
		})
	}
	return
}

func (this *Runner) Drive() (err error) {
	err = this.reader.Drive()
	return
}

func (this *Runner) Close(closemsg string) (err error) {
	state := atomic.LoadInt32(&this.state)
	if state == int32(core.Runner_Stoped) {
		err = core.Error_RunnerStoped
		return
	}
	if state == int32(core.Runner_Stoping) {
		err = core.Error_RunnerStoping
		return
	}
	atomic.StoreInt32(&this.state, int32(core.Runner_Stoping))
	defer lego.Recover(fmt.Sprintf("%s Close:%s", this.conf.Name, closemsg))
	log.Infof("Runner Close : %s", closemsg)
	if this.reader != nil { //初始化启动失败 也需要走一次close
		if err = this.reader.Close(); err != nil {
			log.Errorf("Runner reader Close err: %v", err)
			return
		}
	}
	close(this.readerPope)
	close(this.parserPope)
	if this.parser != nil {
		if err = this.parser.Close(); err != nil {
			log.Errorf("Runner parser Close err: %v", err)
			return
		}
	}
	for i, v := range this.transforms {
		if v != nil {
			close(this.transformsPope[i])
			if err = v.Close(); err != nil {
				log.Errorf("Runner transforms Close err: %v", err)
				return
			}
		}
	}
	for _, v := range this.senders {
		if v != nil {
			close(this.sendersPope[v.GetType()])
			if err = v.Close(); err != nil {
				log.Errorf("Runner sender Close err: %v", err)
				return
			}
		}
	}
	if this.metaer != nil {
		if err = this.metaer.Close(); err != nil {
			log.Errorf("Runner metaer Close err: %v", err)
			return
		}
	}
	atomic.StoreInt32(&this.state, int32(core.Runner_Stoped))
	return
}

///日志***********************************************************************

func (this *Runner) Debugf(format string, a ...interface{}) {
	this.log.Debugf(fmt.Sprintf("[Runner:%s] "+format, this.conf.Name), a...)
}
func (this *Runner) Infof(format string, a ...interface{}) {
	this.log.Infof(fmt.Sprintf("[Runner:%s] "+format, this.conf.Name), a...)
}
func (this *Runner) Warnf(format string, a ...interface{}) {
	this.log.Warnf(fmt.Sprintf("[Runner:%s] "+format, this.conf.Name), a...)
}
func (this *Runner) Errorf(format string, a ...interface{}) {
	this.log.Errorf(fmt.Sprintf("[Runner:%s] "+format, this.conf.Name), a...)
}
func (this *Runner) Panicf(format string, a ...interface{}) {
	this.log.Panicf(fmt.Sprintf("[Runner:%s] "+format, this.conf.Name), a...)
}
func (this *Runner) Fatalf(format string, a ...interface{}) {
	this.log.Fatalf(fmt.Sprintf("[Runner:%s] "+format, this.conf.Name), a...)
}
