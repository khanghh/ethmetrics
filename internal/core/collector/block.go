package collector

import (
	"ethmetrics/internal/core"
	"ethmetrics/internal/logger"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/metrics"
)

const (
	maxCachedBlock = 200
)

type BlockMetrics struct {
}

func (m *BlockMetrics) Setup(ctx *core.Ctx) error {
	return nil
}

func (m *BlockMetrics) getBlock(ctx *core.Ctx, number *big.Int) (*core.Block, error) {
	var block core.Block
	err := ctx.Client.CallContext(ctx.Context, &block, "eth_getBlockByNumber", hexutil.EncodeBig(number), false)
	if err != nil {
		return nil, err
	}
	return &block, err
}

func (m *BlockMetrics) cacheBlock(ctx *core.Ctx, block *core.Block) {
	cachedBlocks, ok := ctx.Storage["block"].([]*core.Block)
	if !ok {
		cachedBlocks = []*core.Block{}
		ctx.Storage["block"] = cachedBlocks
	}
	cachedBlocks = append(cachedBlocks, block)
	if len(cachedBlocks) > maxCachedBlock {
		cachedBlocks = cachedBlocks[1:]
	}
	ctx.Storage["block"] = cachedBlocks
}

func (m *BlockMetrics) Collect(ctx *core.Ctx) {
	block, err := m.getBlock(ctx, ctx.Head.Number)
	if err != nil {
		logger.Errorln(err)
	}
	m.cacheBlock(ctx, block)

	blockHeadGauge := metrics.GetOrRegisterGauge("eth/block.number", ctx.Registry)
	blockTimestampGauge := metrics.GetOrRegisterGauge("eth/block.timestamp", ctx.Registry)
	blockSizeGauge := metrics.GetOrRegisterGaugeFloat64("eth/block.size", ctx.Registry)
	blockGasUsedGauge := metrics.GetOrRegisterGauge("eth/block.gasused", ctx.Registry)
	blockGasLimitGauge := metrics.GetOrRegisterGauge("eth/block.gaslimit", ctx.Registry)
	blockTxCountGauge := metrics.GetOrRegisterGauge("eth/block.txcount", ctx.Registry)

	blockHeadGauge.Update(int64(block.Number.Uint64()))
	blockTimestampGauge.Update(int64(block.Time))
	blockSizeGauge.Update(float64(block.Size))
	blockGasUsedGauge.Update(int64(block.GasUsed))
	blockGasLimitGauge.Update(int64(block.GasLimit))
	blockTxCountGauge.Update(int64(len(block.Transactions)))
	logger.Debugf("Receive new block #%d with %d transactions\n", block.Number, len(block.Transactions))
}
