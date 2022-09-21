package collector

import (
	"ethmetrics/internal/core"
	"ethmetrics/internal/logger"

	"github.com/ethereum/go-ethereum/metrics"
)

type BlockMetrics struct {
}

func (c *BlockMetrics) Setup(ctx *core.Ctx, registry metrics.Registry) error {
	return nil
}

func (c *BlockMetrics) Collect(ctx *core.Ctx) {
	blockHeadGauge := metrics.GetOrRegisterGauge("eth/block/head", ctx.Registry)
	blockSizeGauge := metrics.GetOrRegisterGaugeFloat64("eth/block/size", ctx.Registry)
	blockGasUsedGauge := metrics.GetOrRegisterGauge("eth/block/gas_used", ctx.Registry)
	blockGasLimitGauge := metrics.GetOrRegisterGauge("eth/block/gas_limit", ctx.Registry)

	blockHeadGauge.Update(int64(ctx.LatestBlock.NumberU64()))
	blockSizeGauge.Update(float64(ctx.LatestBlock.Size()))
	blockGasUsedGauge.Update(int64(ctx.LatestBlock.GasUsed()))
	blockGasLimitGauge.Update(int64(ctx.LatestBlock.GasLimit()))
	logger.Debugf("Receive new block #%d\n", ctx.LatestBlock.NumberU64())
}
