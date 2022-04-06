package document

import (
	"lego_logagent/modules/logagent/core"
	doc "lego_logagent/modules/logagent/metaer/document"
	"lego_logagent/modules/logagent/reader"
)

type Reader struct {
	reader.Reader
	options  IOptions                //以接口对象传递参数 方便后期继承扩展
	meta     doc.IDocumentMetaerData //愿数据
	schemas  map[string]string       //sql 类型处理
	sourceip string
}

func (this *Reader) Init(runner core.IRunner, reader core.IReader, meta core.IMetaerData, options core.IReaderOptions) (err error) {
	if err = this.Reader.Init(runner, reader, meta, options); err != nil {
		return
	}
	this.options = options.(IOptions)
	this.meta = meta.(doc.IDocumentMetaerData)
	return
}

///外部调度器 驱动执行  此接口 不可阻塞
func (this *Reader) Drive() (err error) {
	if err = this.Reader.Drive(); err != nil {
		return
	}
	return
}

func (this *Reader) Read(task core.IMetaerNodeData) (err error) {
	return
}
