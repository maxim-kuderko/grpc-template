package initializers

import "github.com/spf13/viper"

func NewMetrics(v *viper.Viper) MetricsReporter {
	return nil
}

type MetricsReporter interface {
	Send(name string, value float64, tags ...string)
	Tracer(name string, tags ...string) Tracer
}

type Tracer interface {
	Start()
	Stop()
}
