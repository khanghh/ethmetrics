package main

import "github.com/urfave/cli/v2"

var (
	rpcUrlFlag = &cli.StringFlag{
		Name:    "rpc-url",
		Usage:   "RPC url of geth node",
		Value:   "ws://localhost:8546",
		EnvVars: []string{"RPC_URL"},
	}
	influxDbUrlFlag = &cli.StringFlag{
		Name:    "influxdb.url",
		Usage:   "InfluxDB url",
		Value:   "http://localhost:8086",
		EnvVars: []string{"INFLUXDB_URL"},
	}
	influxDbTokenFlag = &cli.StringFlag{
		Name:    "influxdb.token",
		Usage:   "InfluxDB admin token",
		EnvVars: []string{"INFLUXDB_TOKEN"},
	}
	influxDbOrgFlag = &cli.StringFlag{
		Name:    "influxdb.org",
		Usage:   "InfluxDB organization",
		Value:   "blockchain",
		EnvVars: []string{"INFLUXDB_ORG"},
	}
	influxDbBucketFlag = &cli.StringFlag{
		Name:    "influxdb.bucket",
		Usage:   "InfluxDB bucket to push metrics data",
		Value:   "geth",
		EnvVars: []string{"INFLUXDB_BUCKET"},
	}
)
