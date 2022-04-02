package services

import (
	"github.com/liwei1dao/lego/base/cluster"
)

type ServiceBase struct {
	cluster.ClusterService
}

func (this *ServiceBase) InitSys() {
	this.ClusterService.InitSys()

}
