package core

const (
	KeyRaw         = "message"
	KeyTimestamp   = "collect_time"
	KeyIdsCollecip = "collect_ip"
	KeyEqptip      = "source_ip"
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
		ServiceId() string
		ServiceIP() string
		MaxProcs() (maxprocs int)
		MaxCollDataSzie() (msgmaxsize uint64)
		GetRunnerState() RunnerState
		Name() (name string)
		Init() (err error)
		Start() (err error)
		Drive() (err error) //驱动工作 提供外部调度器使用 个采集模块根据自己的需求才实现此接口
		Close(state RunnerState, closemsg string) (err error)
		Metaer() IMetaer
		Reader() IReader
		ReaderPipe() chan<- ICollData
		Push_ParserPipe(bucket ICollDataBucket)
		ParserPipe() <-chan ICollDataBucket
		Push_TransformsPipe(bucket ICollDataBucket)
		TransformsPipe(index int) (pipe <-chan ICollDataBucket, ok bool)
		Push_NextTransformsPipe(index int, bucket ICollDataBucket)
		Push_SenderPipe(bucket ICollDataBucket)
		SenderPipe(stype string) (pipe <-chan ICollDataBucket, ok bool)
		SyncRunnerInfo()
		StatisticRunner()
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
		Input() chan<- ICollData
	}
	//读取器
	IParser interface {
		GetRunner() IRunner
		GetType() string
		Start() (err error)
		Close() (err error)
		Parse(bucket ICollDataBucket)
	}
	//变换器
	ITransforms interface {
		GetRunner() IRunner
		Start() (err error)
		Close() (err error)
		Trans(bucket ICollDataBucket)
	}
	//读取器
	ISender interface {
		GetRunner() IRunner
		GetType() string
		Start() (err error)
		Run(pipeId int, pipe <-chan ICollDataBucket, params ...interface{})
		Close() (err error)
		Send(pipeId int, bucket ICollDataBucket, params ...interface{})
		ReadCnt() int64
		ReadAnResetCnt() int64
	}
)