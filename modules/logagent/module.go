package loganget

import (
	"lego_logagent/comm"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
)

func NewModule() core.IModule {
	m := new(LogAnget)
	return m
}

type LogAnget struct {
	cbase.ModuleBase
	service base.IClusterService
}

func (this *LogAnget) GetType() core.M_Modules {
	return comm.SM_LogAngetModule
}
func (this *LogAnget) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.service = service.(base.IClusterService)
	this.ModuleBase.Init(service, module, options)
	return
}
