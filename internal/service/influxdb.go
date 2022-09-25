package service

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/metrics"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
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

func parseName(name string) (string, string) {
	dotIndex := strings.LastIndex(name, ".")
	return name[0:dotIndex], name[dotIndex+1:]
}

func (p *InfluxDBPublisher) PublishMetrics(ctx context.Context, reg metrics.Registry) {
	now := time.Now()
	reg.Each(func(name string, i interface{}) {
		measurementName, fieldName := parseName(name)
		if fieldName == "" {
			fieldName = "value"
		}
		switch metric := i.(type) {
		case metrics.Counter:
			fields := map[string]interface{}{
				fieldName: metric.Count(),
			}
			pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.Gauge:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				fieldName: ms.Value(),
			}
			pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.GaugeFloat64:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				fieldName: ms.Value(),
			}
			pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.Histogram:
			ms := metric.Snapshot()
			if ms.Count() > 0 {
				ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
				fields := map[string]interface{}{
					"count":    ms.Count(),
					"max":      ms.Max(),
					"mean":     ms.Mean(),
					"min":      ms.Min(),
					"stddev":   ms.StdDev(),
					"variance": ms.Variance(),
					"p50":      ps[0],
					"p75":      ps[1],
					"p95":      ps[2],
					"p99":      ps[3],
					"p999":     ps[4],
					"p9999":    ps[5],
				}
				pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
				p.writeAPI.WritePoint(ctx, pt)
			}
		case metrics.Meter:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				"count": float64(ms.Count()),
				"m1":    ms.Rate1(),
				"m5":    ms.Rate5(),
				"m15":   ms.Rate15(),
				"mean":  ms.RateMean(),
			}
			pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.Timer:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			fields := map[string]interface{}{
				"count":    ms.Count(),
				"max":      ms.Max(),
				"mean":     ms.Mean(),
				"min":      ms.Min(),
				"stddev":   ms.StdDev(),
				"variance": ms.Variance(),
				"p50":      ps[0],
				"p75":      ps[1],
				"p95":      ps[2],
				"p99":      ps[3],
				"p999":     ps[4],
				"p9999":    ps[5],
				"m1":       ms.Rate1(),
				"m5":       ms.Rate5(),
				"m15":      ms.Rate15(),
				"meanrate": ms.RateMean(),
			}
			pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
			p.writeAPI.WritePoint(ctx, pt)
		case metrics.ResettingTimer:
			t := metric.Snapshot()
			if len(t.Values()) > 0 {
				ps := t.Percentiles([]float64{50, 95, 99})
				val := t.Values()
				fields := map[string]interface{}{
					"count": len(val),
					"max":   val[len(val)-1],
					"mean":  t.Mean(),
					"min":   val[0],
					"p50":   ps[0],
					"p95":   ps[1],
					"p99":   ps[2],
				}
				pt := influxdb2.NewPoint(measurementName, p.Tags, fields, now)
				p.writeAPI.WritePoint(ctx, pt)
			}
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
