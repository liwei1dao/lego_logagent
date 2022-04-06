package document

import (
	"lego_logagent/modules/logagent/core"
	doc "lego_logagent/modules/logagent/metaer/document"
	"lego_logagent/modules/logagent/reader"
)

func init() {
	reader.RegisterReader(ReaderType, NewReader)
}

const (
	ReaderType = "document"
)

func NewReader(runner core.IRunner, conf map[string]interface{}) (rder core.IReader, err error) {
	var (
		opt IOptions
		r   *Reader
	)
	if opt, err = newOptions(conf); err != nil {
		return
	}
	r = &Reader{}
	if err = r.Init(runner, r, doc.NewMeta("reader"), opt); err != nil {
		return
	}
	rder = r
	return
}
