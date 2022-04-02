package sql

import (
	"lego_logagent/modules/logagent/core"

	lgsql "github.com/liwei1dao/lego/sys/sql"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	MyCollectionType uint8
	IOptions         interface {
		core.IReaderOptions //继承基础配置
		GetSql_type() lgsql.SqlType
		GetSql_sqlurl() string
		GetSql_database() string
		GetSql_querysql() string
		GetSql_tables() []string
		GetSql_limit() int
		GetSql_schema() []string
	}
	Options struct {
		core.ReaderOptions               //继承基础配置
		Sql_type           lgsql.SqlType //sql数据库类型
		Sql_sqlurl         string        //连接字段
		Sql_database       string        //数据库
		Sql_tables         []string      //静态表列表
		Sql_querysql       string        //查询语句
		Sql_limit          int           //查询数量
		Sql_schema         []string      //类型转换
	}
)

func (this *Options) GetSql_type() lgsql.SqlType {
	return this.Sql_type
}

func (this *Options) GetSql_sqlurl() string {
	return this.Sql_sqlurl
}
func (this *Options) GetSql_database() string {
	return this.Sql_database
}
func (this *Options) GetSql_querysql() string {
	return this.Sql_querysql
}
func (this *Options) GetSql_tables() []string {
	return this.Sql_tables
}

func (this *Options) GetSql_limit() int {
	return this.Sql_limit
}
func (this *Options) GetSql_schema() []string {
	return this.Sql_schema
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
