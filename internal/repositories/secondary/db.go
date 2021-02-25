package secondary

import (
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"github.com/spf13/viper"
)

type Db struct {
}

func NewDb(v *viper.Viper) Repo {
	return &Db{}
}

func (d *Db) Get(r requests.Get) (responses.Get, error) {
	panic("implement me")
}
