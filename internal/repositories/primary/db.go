package primary

import (
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
)

type Db struct {
}

func NewDb() Repo {
	return &Db{}
}

func (d *Db) Get(r requests.Get) (responses.Get, error) {
	panic("implement me")
}
