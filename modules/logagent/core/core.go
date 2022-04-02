package core

import "errors"

const (
	KeyRaw         = "message"
	KeyTimestamp   = "collect_time"
	KeyIdsCollecip = "collect_ip"
	KeyEqptip      = "source_ip"
)

var (
	Error_NoRunner            = errors.New("runner is found")
	Error_RunnerStoping       = errors.New("runner stoping")
	Error_RunnerStoped        = errors.New("runner stoped")
	Error_RunnerTaskExecuting = errors.New("task is executing") //任务正在执行中
)

const (
	Runner_Stoped   RunnerState = iota //已停止
	Runner_Initing                     //初始化中
	Runner_Starting                    //启动中
	Runner_Runing                      //运行中
	Runner_Stoping                     //关闭中
)

type (
	RunnerState     int32
	ICollDataBucket interface {
		AddCollData(data ICollData) (full bool)
		IsEmplty() bool
		IsFull() bool
		CurrCap() int
		Items() []ICollData
		SuccItems() []ICollData
		SuccItemsCount() (reslut int)
		ErrItems() []ICollData
		ErrItemsCount() (reslut int)
		Reset()
		SetError(err error)
		Error() (err error)
		ToString(onlysucc bool) (value string, err error) //序列化字符串接口
		MarshalJSON() ([]byte, error)
	}
	//采集数据结构
	ICollData interface {
		SetError(err error)                  //错误信息
		GetError() error                     //错误信息
		GetSource() string                   //数据来源
		GetTime() int64                      //采集时间
		GetData() map[string]interface{}     //整理数据
		GetValue() interface{}               //采集原数据
		GetSize() uint64                     //数据大小
		ToString() (value string, err error) //序列化
		MarshalJSON() ([]byte, error)        //重构json序列化接口
	}
	//采集器结构
	IRunner interface {
		Name() string
		MaxProcs() int
		MaxMessageSzie() uint64
		State() RunnerState
		Metaer() IMetaer
		Init() (err error)
		Start() (err error)
		Close(closemsg string) (err error)
		ReaderPipe() chan<- ICollData
		Push_ParserPipe(data ICollData)
		ParserPipe() <-chan ICollData
		Push_TransformsPipe(data ICollData)
		TransformsPipe(index int) (pipe <-chan ICollData, ok bool)
		Push_NextTransformsPipe(index int, data ICollData)
		Push_SenderPipe(data ICollData)
		SenderPipe(stype string) (pipe <-chan ICollData, ok bool)
	}
	IDB interface {
		WriteMetaData(rname, name string, metae interface{}) error
		ReadMetaData(rname, name string, metae interface{}) error
	}
	IMetaerData interface {
		GetName() string
		GetMetae() interface{} //注意这处返回 指针对象 map对象许返回&map
	}
	//元数据
	IMetaer interface {
		Init(runner IRunner) (err error)
		Close() (err error)
		Read(meta IMetaerData) (err error)
		Write(meta IMetaerData) (err error)
	}
	//读取器
	IReader interface {
		GetRunner() IRunner
		GetType() string
		Start() (err error)
		Drive() (err error) //驱动工作 外部程序驱动采集器工作
		Close() (err error)
		Read(task IMetaerData) (err error)
		Input() chan<- ICollData
	}
	//读取器
	IParser interface {
		GetRunner() IRunner
		GetType() string
		Start() (err error)
		Close() (err error)
		Parse(data ICollData)
	}
	//变换器
	ITransforms interface {
		GetRunner() IRunner
		Start() (err error)
		Close() (err error)
		Trans(data ICollData)
	}
	//读取器
	ISender interface {
		GetRunner() IRunner
		GetType() string
		Start() (err error)
		Run(pipeId int, pipe <-chan ICollData, params ...interface{})
		Close() (err error)
		Send(pipeId int, data ICollData, params ...interface{})
		ReadCnt() int64
		ReadAnResetCnt() int64
	}
)
