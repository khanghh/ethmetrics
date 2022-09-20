package collector

import (
	"ethmetrics/internal/core"

	"github.com/ethereum/go-ethereum/metrics"
)

type BlockMetrics struct {
}

func (c *BlockMetrics) Register(ctx *core.Ctx, registry metrics.Registry) {

}
