package document

import (
	"lego_logagent/modules/logagent/core"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

const (
	LoaclCollection CollectionType = iota //本地文件采集
	FTPCollection                         //ftp文件采集
	SFTPCollection                        //sftp文件采集
)

type (
	CollectionType uint8
	IOptions       interface {
		core.IReaderOptions //继承基础配置
		GetDoc_collectiontype() CollectionType
		GetDoc_server_addr() string
		GetDoc_server_port() int
		GetDoc_server_user() string
		GetDoc_server_password() string
		GetDoc_directory() string
		GetDoc_regularrules() string
	}
	Options struct {
		core.ReaderOptions                 //继承基础配置
		Doc_collectiontype  CollectionType //采集类型
		Doc_server_addr     string         //远程服务地址
		Doc_server_port     int            //远程服务 端口
		Doc_server_user     string         //远程服务账号
		Doc_server_password string         //远程服务密码
		Doc_directory       string         //采集路径
		Doc_regularrules    string         //文件过滤正则规则
	}
)

func (this *Options) GetDoc_collectiontype() CollectionType {
	return this.Doc_collectiontype
}
func (this *Options) GetDoc_server_addr() string {
	return this.Doc_server_addr
}
func (this *Options) GetDoc_server_port() int {
	return this.Doc_server_port
}
func (this *Options) GetDoc_server_user() string {
	return this.Doc_server_user
}
func (this *Options) GetDoc_server_password() string {
	return this.Doc_server_password
}
func (this *Options) GetDoc_directory() string {
	return this.Doc_directory
}
func (this *Options) GetDoc_regularrules() string {
	return this.Doc_regularrules
}

func newOptions(config map[string]interface{}) (opt IOptions, err error) {
	options := &Options{}
	if config != nil {
		if err = mapstructure.Decode(config, options); err == nil {
			err = mapstructure.Decode(config, &options.ReaderOptions)
		}
	}
	opt = options
	return
}
