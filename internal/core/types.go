package core

import (
	"context"

	"github.com/ethereum/go-ethereum/metrics"
)

type MetricsCollector interface {
	Setup(ctx *Ctx, registry metrics.Registry) error
	Collect(ctx *Ctx)
}

type MetricsPublisher interface {
	PublishMetrics(ctx context.Context, registry metrics.Registry)
}
