package core

import "github.com/liwei1dao/lego/utils/mapstructure"

const (
	Encoding_utf8    Encoding = "UTF-8" //utf-8
	Encoding_unicode Encoding = "Unicode"
	Encoding_gbk     Encoding = "GBK"
)

//读取器 基础配置 所有读取器都必须基础此配置
type (
	Encoding       string
	IReaderOptions interface {
		GetType() string
		GetEncoding() Encoding
	}
	ReaderOptions struct {
		Type     string   //读取器类型
		Encoding Encoding //解析编码
	}
)

func NewReaderOptions(config map[string]interface{}) (opt IReaderOptions, err error) {
	options := &ReaderOptions{}
	if config != nil {
		err = mapstructure.Decode(config, options)
	}
	opt = options
	return
}
func (this *ReaderOptions) GetType() string {
	return this.Type
}

func (this *ReaderOptions) GetEncoding() Encoding {
	return this.Encoding
}

//解析起 基础配置 所有读取器都必须基础此配置
type (
	IParserOptions interface {
		GetType() string
		GetLabels() []string
	}
	ParserOptions struct {
		Type   string //读取器类型
		Labels []string
	}
)

func NewParserOptions(config map[string]interface{}) (opt IParserOptions, err error) {
	options := &ParserOptions{}
	if config != nil {
		err = mapstructure.Decode(config, options)
	}
	opt = options
	return
}
func (this *ParserOptions) GetType() string {
	return this.Type
}
func (this *ParserOptions) GetLabels() []string {
	return this.Labels
}

//转换器 基础配置 所有读取器都必须基础此配置
type (
	ITransformsOptions interface {
		GetType() string
	}
	TransformsOptions struct {
		Type string //读取器类型
	}
)

func NewTransformsOptions(config map[string]interface{}) (opt ITransformsOptions, err error) {
	options := &TransformsOptions{}
	if config != nil {
		err = mapstructure.Decode(config, options)
	}
	opt = options
	return
}
func (this *TransformsOptions) GetType() string {
	return this.Type
}

//解析起 基础配置 所有读取器都必须基础此配置
type (
	ISenderOptions interface {
		GetType() string
		GetCacheMaxSzie() int64
	}
	SenderOptions struct {
		Type         string //读取器类型
		CacheMaxSzie int64  //缓存区大小
	}
)

func NewSenderOptions(config map[string]interface{}) (opt ISenderOptions, err error) {
	options := &SenderOptions{
		CacheMaxSzie: 100 * 1024 * 1024,
	}
	if config != nil {
		err = mapstructure.Decode(config, options)
	}
	opt = options
	return
}
func (this *SenderOptions) GetType() string {
	return this.Type
}
func (this *SenderOptions) GetCacheMaxSzie() int64 {
	return this.CacheMaxSzie
}
