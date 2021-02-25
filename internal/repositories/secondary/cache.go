package secondary

import (
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
)

type Cache struct {
	origin Repo
}

func NewCache(origin Repo) Repo {
	return &Cache{origin: origin}
}

func (c *Cache) Get(r requests.Get) (responses.Get, error) {
	return responses.Get{}, nil
}
