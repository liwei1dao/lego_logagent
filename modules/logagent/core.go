package loganget

type (
	RunnerConfig struct {
		Name             string                   `json:"name" bson:"_id"`   //采集器名称 唯一
		MaxProcs         int                      `json:"maxprocs"`          //采集器任务并发数
		RunIp            []string                 `json:"ip" `               //采集器运行ip 默认 0.0.0.0 自动分配
		IsStopped        bool                     `json:"isstopped"`         //是否停止采集
		MaxCollDataSzie  uint64                   `json:"max_colldata_size"` //最大采集数据大小
		ReaderConfig     map[string]interface{}   `json:"reader"`            //读取器配置
		ParserConf       map[string]interface{}   `json:"parser"`            //解析器配置
		TranfromssConfig []map[string]interface{} `json:"transforms"`        //转换器配置
		SendersConfig    []map[string]interface{} `json:"senders"`           //发送器配置
	}
)
