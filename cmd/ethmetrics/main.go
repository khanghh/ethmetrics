package main

import (
	"ethmetrics/internal/core"
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
	rpcUrl := ctx.String(rpcUrlFlag.Name)
	engine := core.NewEthMetrics(core.MetricsOptions{
		RpcUrl:     rpcUrl,
		Publishers: []core.MetricsPublisher{},
	})
	return engine.Start(ctx.Context)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
