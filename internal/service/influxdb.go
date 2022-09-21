package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/metrics"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

const (
	metricsNamespace = "ethmetrics"
)

type InfluxDBPublisher struct {
	client       influxdb2.Client
	writeAPI     api.WriteAPIBlocking
	ServerURL    string
	AuthToken    string
	Organization string
	Bucket       string
	Tags         map[string]string
}

func (p *InfluxDBPublisher) PublishMetrics(ctx context.Context, reg metrics.Registry) {
	now := time.Now()
	reg.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			measurement := fmt.Sprintf("%s/%s.count", metricsNamespace, name)
			fields := map[string]interface{}{
				"value": metric.Count(),
			}
			pt := influxdb2.NewPoint(measurement, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.Gauge:
			ms := metric.Snapshot()
			measurement := fmt.Sprintf("%s/%s.gauge", metricsNamespace, name)
			fields := map[string]interface{}{
				"value": ms.Value(),
			}
			pt := influxdb2.NewPoint(measurement, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.GaugeFloat64:
			ms := metric.Snapshot()
			measurement := fmt.Sprintf("%s/%s.gauge", metricsNamespace, name)
			fields := map[string]interface{}{
				"value": ms.Value(),
			}
			pt := influxdb2.NewPoint(measurement, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.Histogram:
		case metrics.Meter:
		case metrics.Timer:
		case metrics.ResettingTimer:
		}
	})
	p.writeAPI.Flush(ctx)
}

func NewInfluxDBPublisher(serverURL, authToken, org, bucket string, tags map[string]string) *InfluxDBPublisher {
	client := influxdb2.NewClient(serverURL, authToken)
	writeAPI := client.WriteAPIBlocking(org, bucket)
	return &InfluxDBPublisher{
		client:       client,
		writeAPI:     writeAPI,
		ServerURL:    serverURL,
		AuthToken:    authToken,
		Organization: org,
		Bucket:       bucket,
		Tags:         tags,
	}
}
