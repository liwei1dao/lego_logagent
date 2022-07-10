package loganget

import (
	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"

	"github.com/liwei1dao/lego/sys/gin"
	"github.com/liwei1dao/lego/sys/gin/engine"
	"github.com/liwei1dao/lego/sys/log"
)

type API_Comp struct {
	cbase.ModuleCompBase
	options *Options
	module  *LogAnget
	gin     gin.ISys
}

func (this *API_Comp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(*Options)
	this.module = module.(*LogAnget)
	this.gin, err = gin.NewSys(gin.SetListenPort(this.options.ListenPort))
	logkit := this.gin.Group("/datacollector")
	logkit.POST("/createrunner", this.CreateRunnerReq)
	return
}
func (this *API_Comp) CreateRunnerReq(c *engine.Context) {
	defer lego.Recover("CreateRunnerReq")
	req := &RunnerConfig{
		RunIp:          []string{core.AutoIp},
		MaxProcs:       8,
		MaxMessageSzie: 2 * 1024 * 1024,
	}
	c.ShouldBindJSON(req)
	log.Debugf("AddNeRunnerReq:%+v", req)
}
