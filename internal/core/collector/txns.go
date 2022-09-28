package collector

import (
	"ethmetrics/internal/core"
	"time"

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
	startTime := blocks[0].ReceivedAt
	duration := blocks[len(blocks)-1].ReceivedAt.Sub(startTime)
	return float64(totalTxns) / float64(duration/time.Second)
}

func (m *TxnsMetrics) Collect(ctx *core.Ctx) {
	// TODO: add more tps measurement
	tpsAvg100Gauge := metrics.GetOrRegisterGaugeFloat64("eth/txns/tps.avg100", ctx.Registry)
	tpsAvg100Gauge.Update(m.calculateTps(ctx.CachedBlocks))
}
