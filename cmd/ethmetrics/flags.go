package main

import "github.com/urfave/cli/v2"

var (
	debugFlag = &cli.BoolFlag{
		Name:    "debug",
		Usage:   "Enable debug log",
		EnvVars: []string{"DEBUG"},
	}
	rpcUrlFlag = &cli.StringFlag{
		Name:    "rpc-url",
		Usage:   "RPC url of geth node",
		Value:   "ws://localhost:8546",
		EnvVars: []string{"RPC_URL"},
	}
	influxDbUrlFlag = &cli.StringFlag{
		Name:    "influxdb.url",
		Usage:   "InfluxDB API url to push metrics to",
		Value:   "http://localhost:8086",
		EnvVars: []string{"INFLUXDB_URL"},
	}
	influxDbTokenFlag = &cli.StringFlag{
		Name:    "influxdb.token",
		Usage:   "Authentication token to access InfluxDB database",
		EnvVars: []string{"INFLUXDB_TOKEN"},
	}
	influxDbOrgFlag = &cli.StringFlag{
		Name:    "influxdb.org",
		Usage:   "InfluxDB organization name",
		Value:   "blockchain",
		EnvVars: []string{"INFLUXDB_ORG"},
	}
	influxDbBucketFlag = &cli.StringFlag{
		Name:    "influxdb.bucket",
		Usage:   "InfluxDB bucket to push metrics to",
		Value:   "geth",
		EnvVars: []string{"INFLUXDB_BUCKET"},
	}
	influxDbTagsFlag = &cli.StringFlag{
		Name:    "influxdb.tags",
		Usage:   "Comma-separated InfluxDB tags (key/values) attached to all measurements",
		Value:   "job=ethmetrics",
		EnvVars: []string{"INFLUXDB_TAGS"},
	}
)
