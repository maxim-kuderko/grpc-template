package service

import (
	"github.com/maxim-kuderko/service-template/internal/initializers"
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/internal/repositories/secondary"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
)

type ServiceFunc func(r interface{}) (interface{}, error)

type Service struct {
	primaryRepo   primary.Repo
	secondaryRepo secondary.Repo
	m             initializers.MetricsReporter
}

func NewService(p primary.Repo, s secondary.Repo, metrics initializers.MetricsReporter) *Service {
	return &Service{
		primaryRepo:   p,
		secondaryRepo: s,
		m:             metrics,
	}
}

func (s *Service) Get(r requests.Get) (responses.Get, error) {
	return s.primaryRepo.Get(r)
}
