package service

import "github.com/ethereum/go-ethereum/metrics"

type InfluxDBPublisher struct {
	DatabaseUrl  string
	AuthToken    string
	Organization string
	Bucket       string
}

func (p *InfluxDBPublisher) PublishMetrics(registry metrics.Registry) {

}
