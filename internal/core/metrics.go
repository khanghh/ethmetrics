package core

import (
	"context"
	"ethmetrics/internal/logger"

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
	LatestBlock  *types.Block
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
	e.LatestBlock = block
	e.CachedBlocks = append(e.CachedBlocks, block)
	if len(e.CachedBlocks) > e.MaxCachedBlock {
		e.CachedBlocks = e.CachedBlocks[1:]
	}
	for _, collector := range e.Collectors {
		collector.Collect(&e.Ctx)
	}
}

func (e *EthMetrics) publishMetrics() {
	for _, publisher := range e.Publishers {
		publisher.PublishMetrics(e.Context, e.Registry())
	}
}

func (e *EthMetrics) collectOnNewHead() error {
	client := ethclient.NewClient(e.Client)
	headCh := make(chan *types.Header)
	subs, err := client.SubscribeNewHead(e.Context, headCh)
	if err != nil {
		return err
	}
	for {
		select {
		case head := <-headCh:
			block, err := client.BlockByNumber(e.Context, head.Number)
			if err != nil {
				return err
			}
			e.collectMetrics(block)
			e.publishMetrics()
		case err := <-subs.Err():
			return err
		case <-e.Context.Done():
			return e.Context.Err()
		}
	}
}

func (e *EthMetrics) Start(ctx context.Context) error {
	logger.Println("Dialing RPC node", e.RpcUrl)
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
	return e.collectOnNewHead()
}

func NewEthMetrics(opts MetricsOptions) *EthMetrics {
	return &EthMetrics{
		MetricsOptions: opts,
	}
}
