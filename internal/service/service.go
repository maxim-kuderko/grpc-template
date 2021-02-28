package service

import (
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
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
	ctx, sp := trace.SpanFromContext(r.Context()).Tracer().Start(r.Context(), `get-service`)
	defer sp.End()
	r.WithContext(ctx)
	return s.primaryRepo.Get(r)
}
