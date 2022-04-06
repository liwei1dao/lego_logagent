package document

import (
	"lego_logagent/modules/logagent/core"
)

func NewMeta(name string) (meta IDocumentMetaerData) {
	meta = &SqlMetaData{
		Name:  name,
		nodes: make(map[string]IDocumentMetaNodeData),
	}
	return
}

type (
	IDocumentMetaerData interface {
		core.IMetaerData
		GetNodeData(name string) (node IDocumentMetaNodeData, ok bool)
		SetNodeData(name string, node IDocumentMetaNodeData)
	}
	IDocumentMetaNodeData interface {
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
		nodes map[string]IDocumentMetaNodeData
	}
	DocumentMetaNodeData struct {
		Name       string //表名
		Size       uint64 //表的数据长度
		ReadOffset uint64 //已采集数据长度
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
func (this *SqlMetaData) GetNodeData(name string) (node IDocumentMetaNodeData, ok bool) {
	node, ok = this.nodes[name]
	return
}
func (this *SqlMetaData) SetNodeData(name string, node IDocumentMetaNodeData) {
	this.nodes[name] = node
	return
}

func (this *DocumentMetaNodeData) Get_TableName() string {
	return this.Name
}
func (this *DocumentMetaNodeData) Set_TableName(v string) {
	this.Name = v
}
func (this *DocumentMetaNodeData) Get_TableDataCount() uint64 {
	return this.Size
}
func (this *DocumentMetaNodeData) Set_TableDataCount(v uint64) {
	this.Size = v
}
func (this *DocumentMetaNodeData) Get_TableAlreadyReadOffset() uint64 {
	return this.ReadOffset
}
func (this *DocumentMetaNodeData) Set_TableAlreadyReadOffset(v uint64) {
	this.ReadOffset = v
}
