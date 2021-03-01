package service

import (
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"go.opentelemetry.io/otel/metric"
)

type ServiceFunc func(r interface{}) (interface{}, error)

type Service struct {
	primaryRepo primary.Repo
	m           metric.Meter
}

func NewService(p primary.Repo, metrics func() metric.MeterProvider) *Service {
	return &Service{
		primaryRepo: p,
		m:           metrics().Meter(`service`),
	}
}

func (s *Service) Get(r requests.Get) (responses.Get, error) {
	return s.primaryRepo.Get(r)
}
