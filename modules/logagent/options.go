package loganget

import (
	"github.com/liwei1dao/lego/lib/modules/http"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		http.IOptions
	}
	Options struct {
		http.Options
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {

	if err = this.Options.LoadConfig(settings); err == nil {
		if settings != nil {
			err = mapstructure.Decode(settings, this)
		}
	}
	return
}
