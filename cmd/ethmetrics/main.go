package main

import (
	"ethmetrics/internal/core"
	"ethmetrics/internal/core/collector"
	"ethmetrics/internal/logger"
	"ethmetrics/internal/service"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

var (
	gitCommit = ""
	gitDate   = ""

	app   *cli.App
	flags = []cli.Flag{
		debugFlag,
		rpcUrlFlag,
		influxDbUrlFlag,
		influxDbTokenFlag,
		influxDbOrgFlag,
		influxDbBucketFlag,
	}
)

func init() {
	app = cli.NewApp()
	app.Version = params.VersionWithCommit(gitCommit, gitDate)
	app.Usage = "Ethereum metrics collector"
	app.Flags = append(app.Flags, flags...)
	app.Action = run
}

func run(ctx *cli.Context) error {
	debugEnabled := ctx.Bool(debugFlag.Name)
	if debugEnabled {
		logger.SetLogLevel(logger.LevelDebug)
	}
	rpcUrl := ctx.String(rpcUrlFlag.Name)
	influxdbServerUrl := ctx.String(influxDbUrlFlag.Name)
	influxdbToken := ctx.String(influxDbTokenFlag.Name)
	influxdbOrg := ctx.String(influxDbOrgFlag.Name)
	influxdbBucket := ctx.String(influxDbBucketFlag.Name)
	influxdbPublisher := service.NewInfluxDBPublisher(influxdbServerUrl, influxdbToken, influxdbOrg, influxdbBucket)
	engine := core.NewEthMetrics(core.MetricsOptions{
		RpcUrl: rpcUrl,
		Collectors: []core.MetricsCollector{
			&collector.BlockMetrics{},
		},
		Publishers: []core.MetricsPublisher{
			influxdbPublisher,
		},
	})
	return engine.Start(ctx.Context)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
