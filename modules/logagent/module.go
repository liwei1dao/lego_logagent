package loganget

import (
	"lego_logagent/comm"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/lib/modules/http"
)

type LogAnget struct {
	http.Http
}

func (this *LogAnget) GetType() core.M_Modules {
	return comm.SM_LogAngetModule
}
