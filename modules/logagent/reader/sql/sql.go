package sql

import (
	"database/sql"
	"fmt"
	"lego_logagent/modules/logagent/core"
	msql "lego_logagent/modules/logagent/metaer/sql"
	"lego_logagent/modules/logagent/reader"
	"strings"
	"sync/atomic"

	"github.com/liwei1dao/lego/sys/log"
	lgsql "github.com/liwei1dao/lego/sys/sql"
	"github.com/liwei1dao/lego/sys/sql/convert"
)

type Reader struct {
	reader.Reader
	options  IOptions            //以接口对象传递参数 方便后期继承扩展
	meta     msql.ITableMetaData //愿数据
	schemas  map[string]string   //sql 类型处理
	sql      lgsql.ISys
	sourceip string
}

func (this *Reader) Init(runner core.IRunner, reader core.IReader, meta core.IMetaerData, options core.IReaderOptions) (err error) {
	if err = this.Reader.Init(runner, reader, meta, options); err != nil {
		return
	}
	if this.schemas, err = convert.SchemaCheck(this.options.GetSql_schema()); err != nil {
		return
	}
	this.sql, err = lgsql.NewSys(lgsql.SetSqlType(this.options.GetSql_type()), lgsql.SetSqlUrl(this.options.GetSql_sqlurl()))
	return
}

func (this *Reader) Start() (err error) {
	err = this.Reader.Start()
	return
}

///外部调度器 驱动执行  此接口 不可阻塞
func (this *Reader) Drive() (err error) {
	if err = this.Reader.Drive(); err != nil {
		return
	}
	var (
		collectiontables []core.IMetaerData
		tables           map[string]struct{}
		tablecount       uint64
		table            *msql.TableMeta
		ok               bool
	)
	if tables, err = this.scanddatabase(); err == nil && len(tables) > 0 {
		for _, k := range this.options.GetSql_tables() {
			if _, ok = tables[k]; ok {
				tablecount = this.gettablecount(k)
				if table, ok = this.meta.GetTableMeta(k); !ok {
					table = &msql.TableMeta{
						TableName:              k,
						TableDataCount:         tablecount,
						TableAlreadyReadOffset: 0,
					}
					this.meta.SetTableMeta(k, table)
				}
				if table.TableAlreadyReadOffset < tablecount { //有新的数据
					table.TableDataCount = tablecount
					collectiontables = append(collectiontables, table)
				}
			} else {
				log.Errorf("Reader DM not found table:%s", k)
			}
		}
		log.Debugf("collectiontables:%v", collectiontables)
		if len(collectiontables) > 0 {
			this.TaskPipe = make(chan core.IMetaerData, len(collectiontables))
			for _, v := range collectiontables {
				this.TaskPipe <- v
			}
			return
		} else {
			log.Debugf("Reader scan no found new data!")
		}
	}
	atomic.StoreInt32(&this.State, 0)
	return
}
func (this *Reader) Read(task core.ITask) (err error) {
	table := task.(*msql.TableMeta)
clocp:
	for {
		if ok, err := this.collection_table(table); ok || err != nil {
			log.Debugf("Reader collection_table ok:%v,err:%v", ok, err)
			break clocp
		}
		if this.Runner.State() == core.Runner_Stoping { //采集器进入停止过程中
			log.Debugf("Reader asyncollection exit")
			break clocp
		}
	}
	return
}

//-------------------------------------------------------------------------------------------------------------------------------
//扫描数据库 扫描数据库下所有的表 以及表中的数据条数
func (this *Reader) scanddatabase() (tables map[string]struct{}, err error) {
	var (
		data *sql.Rows
		sql  string
	)
	tables = make(map[string]struct{})
	sql = "select table_name from all_tables"
	switch this.options.GetSql_type() {
	case lgsql.DM:
		sql = "select table_name from all_tables"
		break
	default:
		sql = "select table_name from all_tables"
	}
	if data, err = this.sql.Query(sql); err == nil {
		tablename := ""
		for data.Next() {
			if e := data.Scan(&tablename); e == nil {
				tables[tablename] = struct{}{}
			}
		}
	}
	return
}

//获取表长度
func (this *Reader) gettablecount(tablename string) (count uint64) {
	var (
		sqlstr string
		err    error
		data   *sql.Rows
	)
	count = 0
	sqlstr = fmt.Sprintf(`select count(*) from %s.%s`, this.options.GetSql_database(), tablename)
	if data, err = this.sql.Query(sqlstr); err != nil {
		log.Errorf("gettablecount %s sql:%s err:%v", tablename, sqlstr, err)
	} else {
		for data.Next() {
			if err := data.Scan(&count); err != nil {
				log.Errorf("gettablecount %s sql:%s err:%v", tablename, sqlstr, err)
			}
		}
	}
	return
}

//查询语句
func (this *Reader) getsqlstr(table *msql.TableMeta) (sqlstr string) {
	sqlstr = strings.Replace(this.options.GetSql_querysql(), "$TABLE$", fmt.Sprintf("%s.%s", this.options.GetSql_database(), table.TableName), -1)
	sqlstr = fmt.Sprintf("%s limit %d,%d;", sqlstr, table.TableAlreadyReadOffset, table.TableAlreadyReadOffset+uint64(this.options.GetSql_limit()))
	return
}

//采集数据表
func (this *Reader) collection_table(table *msql.TableMeta) (isend bool, err error) {
	var (
		sqlstr   string
		data     *sql.Rows
		columns  []string
		scanArgs []interface{}
		nochiced []bool
	)
	sqlstr = this.getsqlstr(table)
	log.Debugf("sql:%s", sqlstr)
	if data, err = this.sql.Query(sqlstr); err == nil {
		if columns, err = data.Columns(); err != nil {
			log.Errorf("collection_table err%v", err)
			return
		}
		schemas := make(map[string]string)
		for k, v := range this.schemas {
			schemas[k] = v
		}
		scanArgs, nochiced = convert.GetInitScans(len(columns), data, schemas, table.TableName)
		isend, err = this.getAllDatas(table, data, scanArgs, columns, nochiced, schemas)
	}
	return
}

//读取数据
func (this *Reader) getAllDatas(table *msql.TableMeta, rows *sql.Rows, scanArgs []interface{}, columns []string, nochiced []bool, schemas map[string]string) (isend bool, err error) {
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Errorf("getAllDatas scan rows err:%v", err)
			continue
		}
		// this.Runner.Debugf("getAllDatas data:%v", scanArgs)
		var (
			data = make(map[string]interface{}, len(scanArgs))
		)
		for i := 0; i < len(scanArgs); i++ {
			_, err := convert.ConvertScanArgs(data, scanArgs[i], columns[i], table.TableName, nochiced[i], schemas)
			if err != nil {
				log.Errorf("getAllDatas ConvertScanArgs err:%v", err)
			}
		}
		err = nil
		if len(data) <= 0 {
			continue
		}
		isend = this.writeDataChan(table, data)
	}
	return
}

//发送数据 数据量打的时候会有堵塞
func (this *Reader) writeDataChan(table *msql.TableMeta, data map[string]interface{}) (isend bool) {
	this.Input() <- core.NewCollData(this.sourceip, data)
	table.TableAlreadyReadOffset += 1
	if table.TableAlreadyReadOffset >= table.TableDataCount {
		isend = true
	}
	return
}
