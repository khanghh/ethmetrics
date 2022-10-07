package collector

import (
	"ethmetrics/internal/core"
	"ethmetrics/internal/logger"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/metrics"
)

type BlockMetrics struct {
}

func (m *BlockMetrics) Setup(ctx *core.Ctx, registry metrics.Registry) error {
	return nil
}

func (m *BlockMetrics) avgBlockTime(blocks []*types.Block, numBlock int) {

}

func (m *BlockMetrics) Collect(ctx *core.Ctx) {
	latestBlock := ctx.CachedBlocks[len(ctx.CachedBlocks)-1]
	blockHeadGauge := metrics.GetOrRegisterGauge("eth/block.number", ctx.Registry)
	blockTimestampGauge := metrics.GetOrRegisterGauge("eth/block.timestamp", ctx.Registry)
	blockSizeGauge := metrics.GetOrRegisterGaugeFloat64("eth/block.size", ctx.Registry)
	blockGasUsedGauge := metrics.GetOrRegisterGauge("eth/block.gasused", ctx.Registry)
	blockGasLimitGauge := metrics.GetOrRegisterGauge("eth/block.gaslimit", ctx.Registry)
	blockTxCountGauge := metrics.GetOrRegisterGauge("eth/block.txcount", ctx.Registry)

	blockHeadGauge.Update(int64(latestBlock.NumberU64()))
	blockTimestampGauge.Update(int64(latestBlock.Time()))
	blockSizeGauge.Update(float64(latestBlock.Size()))
	blockGasUsedGauge.Update(int64(latestBlock.GasUsed()))
	blockGasLimitGauge.Update(int64(latestBlock.GasLimit()))
	blockTxCountGauge.Update(int64(latestBlock.Transactions().Len()))
	logger.Debugf("Receive new block #%d with %d transactions\n", latestBlock.NumberU64(), latestBlock.Transactions().Len())
}
