package main

import (
	"ethmetrics/internal/core"
	"ethmetrics/internal/core/collector"
	"ethmetrics/internal/logger"
	"ethmetrics/internal/service"
	"fmt"
	"os"
	"strings"

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
		influxDbTagsFlag,
	}
)

func init() {
	app = cli.NewApp()
	app.Version = params.VersionWithCommit(gitCommit, gitDate)
	app.Usage = "Ethereum metrics collector"
	app.Flags = append(app.Flags, flags...)
	app.Action = run
}

func splitTagsFlag(tagsFlag string) map[string]string {
	tags := strings.Split(tagsFlag, ",")
	tagsMap := map[string]string{}

	for _, t := range tags {
		if t != "" {
			kv := strings.Split(t, "=")

			if len(kv) == 2 {
				tagsMap[kv[0]] = kv[1]
			}
		}
	}

	return tagsMap
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
	influxdbTags := splitTagsFlag(ctx.String(influxDbTagsFlag.Name))
	influxdbPublisher := service.NewInfluxDBPublisher(influxdbServerUrl, influxdbToken, influxdbOrg, influxdbBucket, influxdbTags)
	engine := core.NewEthMetrics(core.MetricsOptions{
		RpcUrl:         rpcUrl,
		MaxCachedBlock: 200,
		Collectors: []core.MetricsCollector{
			&collector.BlockMetrics{},
			&collector.TxsMetrics{},
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
