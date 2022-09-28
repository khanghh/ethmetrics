package collector

import (
	"ethmetrics/internal/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/metrics"
)

type TxnsMetrics struct {
}

func (m *TxnsMetrics) Setup(ctx *core.Ctx, registry metrics.Registry) error {
	return nil
}

func (m *TxnsMetrics) calculateTps(blocks []*types.Block) float64 {
	if len(blocks) < 2 {
		return 0
	}
	totalTxns := 0
	for _, block := range blocks {
		totalTxns += block.Transactions().Len()
	}
	duration := blocks[len(blocks)-1].Time() - blocks[0].Time()
	return float64(totalTxns) / float64(duration)
}

func (m *TxnsMetrics) Collect(ctx *core.Ctx) {
	// TODO: add more block ranges tps measurement
	tpsAvg100Gauge := metrics.GetOrRegisterGaugeFloat64("eth/txns/tps.avg100", ctx.Registry)
	tpsAvg100Gauge.Update(m.calculateTps(ctx.CachedBlocks))
}
