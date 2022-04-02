package metaer

type MetaNodeBase struct {
	Name string
}

func (this *MetaNodeBase) GetName() string {
	return this.Name
}
