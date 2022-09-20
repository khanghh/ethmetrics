package core

import (
	"github.com/ethereum/go-ethereum/metrics"
)

type MetricsCollector interface {
	Setup(ctx *Ctx, registry metrics.Registry) error
	Collect(ctx *Ctx)
}

type MetricsPublisher interface {
	PublishMetrics(registry metrics.Registry)
}
