package primary

import (
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
)

type Repo interface {
	Get(r requests.Get) (responses.Get, error)
}
