package secondary

import (
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"github.com/spf13/viper"
)

type Cache struct {
	origin Repo
}

func NewCache(origin Repo, v *viper.Viper) Repo {
	return &Cache{origin: origin}
}

func NewCachedDB(v *viper.Viper) Repo {
	return NewCache(NewDb(v), v)
}

func (c *Cache) Get(r requests.Get) (responses.Get, error) {
	return responses.Get{}, nil
}
