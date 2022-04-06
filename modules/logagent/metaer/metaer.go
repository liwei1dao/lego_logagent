package metaer

import (
	"lego_logagent/modules/logagent/core"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/mgo"
	"github.com/liwei1dao/lego/sys/redis"
)

func NewMetaer(runner core.IRunner) (metaer core.IMetaer, err error) {
	metaer = &RedisMetaer{}
	err = metaer.Init(runner)
	return
}

type RedisMetaer struct {
	runner core.IRunner
	db     core.IDB
	lock   sync.RWMutex
	meta   core.IMetaerData
}

func (this *RedisMetaer) Init(runner core.IRunner) (err error) {
	this.runner = runner
	return
}

//注册加载元数据
func (this *RedisMetaer) Read(meta core.IMetaerData) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.meta = meta
	if err = this.db.ReadMetaData(this.runner.Name(), this.meta.GetName(), this.meta); err != nil {
		if err == redis.RedisNil || err == mgo.MongodbNil {
			err = nil
		} else {
			log.Errorf("RedisMetaer Read:%s err:%v", this.meta.GetName(), err)
		}
	}
	return
}

//注册加载元数据
func (this *RedisMetaer) Write() (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if err = this.db.WriteMetaData(this.runner.Name(), this.meta.GetName(), this.meta); err != nil {
		log.Errorf("RedisMetaer Write:%s err:%v", this.meta.GetName(), err)
	}
	return
}

//采集器需要关闭时元数据这边需要保存处理
func (this *RedisMetaer) Close() (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if err = this.Write(); err != nil {
		log.Errorf("MetaData:%s Write Fatal err:%v", this.meta.GetName(), err)
		return
	}
	return
}
