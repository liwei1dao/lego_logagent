package sql

import (
	"lego_logagent/modules/logagent/core"
)

func NewMeta(name string) (meta ISqlMetaerData) {
	meta = &SqlMetaData{
		Name:  name,
		nodes: make(map[string]ISqlMetaNodeData),
	}
	return
}

type (
	ISqlMetaerData interface {
		core.IMetaerData
		GetNodeData(name string) (node ISqlMetaNodeData, ok bool)
		SetNodeData(name string, node ISqlMetaNodeData)
	}
	ISqlMetaNodeData interface {
		core.IMetaerNodeData
		Get_TableName() string
		Set_TableName(v string)
		Get_TableDataCount() uint64
		Set_TableDataCount(v uint64)
		Get_TableAlreadyReadOffset() uint64
		Set_TableAlreadyReadOffset(v uint64)
	}
	SqlMetaData struct {
		Name  string
		nodes map[string]ISqlMetaNodeData
	}
	SqlMetaNodeData struct {
		TableName              string //表名
		TableDataCount         uint64 //表的数据长度
		TableAlreadyReadOffset uint64 //已采集数据长度
	}
)

//map 结构对象序列化需要 map的指针
func (this *SqlMetaData) GetName() string {
	return this.Name
}

//map 结构对象序列化需要 map的指针
func (this *SqlMetaData) GetValue() core.IMetaerNodeData {
	return &this.nodes
}
func (this *SqlMetaData) GetNodeData(name string) (node ISqlMetaNodeData, ok bool) {
	node, ok = this.nodes[name]
	return
}
func (this *SqlMetaData) SetNodeData(name string, node ISqlMetaNodeData) {
	this.nodes[name] = node
	return
}

func (this *SqlMetaNodeData) Get_TableName() string {
	return this.TableName
}
func (this *SqlMetaNodeData) Set_TableName(v string) {
	this.TableName = v
}
func (this *SqlMetaNodeData) Get_TableDataCount() uint64 {
	return this.TableDataCount
}
func (this *SqlMetaNodeData) Set_TableDataCount(v uint64) {
	this.TableDataCount = v
}
func (this *SqlMetaNodeData) Get_TableAlreadyReadOffset() uint64 {
	return this.TableAlreadyReadOffset
}
func (this *SqlMetaNodeData) Set_TableAlreadyReadOffset(v uint64) {
	this.TableAlreadyReadOffset = v
}
