package core

import (
	"context"
	"ethmetrics/internal/logger"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
)

type MetricsOptions struct {
	RpcUrl         string
	MaxCachedBlock int
	Collectors     []MetricsCollector
	Publishers     []MetricsPublisher
}

type Ctx struct {
	Client       *rpc.Client
	Context      context.Context
	CachedBlocks []*types.Block
	Head         *types.Header
	Registry     metrics.Registry
}

type EthMetrics struct {
	MetricsOptions
	Ctx
}

func init() {
	metrics.Enabled = true
}

func (e *EthMetrics) Registry() metrics.Registry {
	return e.Ctx.Registry
}

func (e *EthMetrics) collectMetrics(block *types.Block) {
	for _, collector := range e.Collectors {
		collector.Collect(&e.Ctx)
	}
}

func (e *EthMetrics) publishMetrics() {
	for _, publisher := range e.Publishers {
		publisher.PublishMetrics(e.Context, e.Registry())
	}
}

func (e *EthMetrics) cacheBlock(block *types.Block) {
	e.CachedBlocks = append(e.CachedBlocks, block)
	if len(e.CachedBlocks) > e.MaxCachedBlock {
		e.CachedBlocks = e.CachedBlocks[1:]
	}
}

func (e *EthMetrics) collectOnNewHead(ctx context.Context) error {
	client := ethclient.NewClient(e.Client)
	newHeadCh := make(chan *types.Header)
	newHeadSub, err := client.SubscribeNewHead(e.Context, newHeadCh)
	if err != nil {
		return err
	}
	for {
		select {
		case head := <-newHeadCh:
			startTime := time.Now()
			block, err := client.BlockByNumber(ctx, head.Number)
			if err != nil {
				return err
			}
			e.Head = head
			e.cacheBlock(block)
			e.collectMetrics(block)
			e.publishMetrics()
			elapsed := time.Since(startTime)
			if elapsed > 1*time.Second {
				logger.Warnf("Process block #%d took %v")
			}
		case err := <-newHeadSub.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *EthMetrics) Start(ctx context.Context) error {
	logger.Println("Dialing RPC node ", e.RpcUrl)
	client, err := rpc.DialContext(ctx, e.RpcUrl)
	if err != nil {
		logger.Errorln("Failed to dial rpc node", e.RpcUrl, err)
		return err
	}
	e.Ctx = Ctx{
		Context:  ctx,
		Client:   client,
		Registry: metrics.NewRegistry(),
	}
	return e.collectOnNewHead(ctx)
}

func NewEthMetrics(opts MetricsOptions) *EthMetrics {
	return &EthMetrics{
		MetricsOptions: opts,
	}
}
