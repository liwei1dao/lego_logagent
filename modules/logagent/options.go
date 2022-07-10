package loganget

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	Options struct {
		ListenPort int
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		err = mapstructure.Decode(settings, this)
	}
	return
}
