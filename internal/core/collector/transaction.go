package collector

import (
	"ethmetrics/internal/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/metrics"
)

type TxsMetrics struct {
}

func (m *TxsMetrics) Setup(ctx *core.Ctx, registry metrics.Registry) error {
	return nil
}

func (m *TxsMetrics) calculateTps(blocks []*types.Block, numBlock int) float64 {
	if len(blocks) < 2 {
		return 0
	}
	calcBlocks := blocks
	if len(blocks) > numBlock {
		calcBlocks = blocks[len(blocks)-numBlock-1:]
	}
	totalTxns := 0
	for _, block := range calcBlocks[1:] {
		totalTxns += block.Transactions().Len()
	}
	duration := calcBlocks[len(calcBlocks)-1].Time() - calcBlocks[0].Time()
	tps := float64(totalTxns) / float64(duration)
	return tps
}

func (m *TxsMetrics) Collect(ctx *core.Ctx) {
	// TODO: add more block ranges tps measurement
	tpsAvg100Gauge := metrics.GetOrRegisterGaugeFloat64("eth/txs/tps.avg100", ctx.Registry)
	tpsAvg100Gauge.Update(m.calculateTps(ctx.CachedBlocks, 100))
}
