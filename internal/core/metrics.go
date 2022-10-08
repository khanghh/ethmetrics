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
	RpcUrl     string
	Collectors []MetricsCollector
	Publishers []MetricsPublisher
}

type Ctx struct {
	Client   *rpc.Client
	Context  context.Context
	Head     *types.Header
	Registry metrics.Registry
	Storage  map[string]any
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

func (e *EthMetrics) setup() error {
	for _, collector := range e.Collectors {
		err := collector.Setup(&e.Ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EthMetrics) collectMetrics() {
	for _, collector := range e.Collectors {
		collector.Collect(&e.Ctx)
	}
}

func (e *EthMetrics) publishMetrics() {
	for _, publisher := range e.Publishers {
		publisher.PublishMetrics(e.Context, e.Registry())
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
			e.Head = head
			e.collectMetrics()
			e.publishMetrics()
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
		Storage:  make(map[string]any),
	}
	if err := e.setup(); err != nil {
		return err
	}
	return e.collectOnNewHead(ctx)
}

func NewEthMetrics(opts MetricsOptions) *EthMetrics {
	return &EthMetrics{
		MetricsOptions: opts,
	}
}
