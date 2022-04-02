package loganget

import (
	"lego_logagent/modules/logagent/core"
)

type Runner struct {
	conf           *RunnerConfig
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

func (this *Runner) Init() (err error) {
	return
}
