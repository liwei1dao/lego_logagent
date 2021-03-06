package sql

import (
	"lego_logagent/modules/logagent/core"
	"lego_logagent/modules/logagent/metaer/sql"
	"lego_logagent/modules/logagent/reader"
)

func init() {
	reader.RegisterReader(ReaderType, NewReader)
}

const (
	ReaderType = "sql"
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
	if err = r.Init(runner, r, sql.NewMeta("reader"), opt); err != nil {
		return
	}
	rder = r
	return
}
