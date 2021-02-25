package service

import (
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/internal/repositories/secondary"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"go.opentelemetry.io/otel/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	primaryRepo   primary.Repo
	secondaryRepo secondary.Repo

	t trace.Tracer
	m metric.Meter
}

func NewService(p primary.Repo, s secondary.Repo, tracer func() *sdktrace.TracerProvider, metrics func() metric.MeterProvider) *Service {
	return &Service{
		primaryRepo:   p,
		secondaryRepo: s,
		t:             tracer().Tracer(`service`),
		m:             metrics().Meter(`service`),
	}
}

func (s *Service) Get(r requests.Get) (responses.Get, error) {
	_, span := s.t.Start(r, `get`)
	defer span.End()

	return responses.Get{}, nil
}
