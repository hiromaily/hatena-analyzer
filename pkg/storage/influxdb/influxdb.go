package influxdb

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func Ping(client influxdb2.Client, org, bucket string) error {
	// Create a write API
	writeAPI := client.WriteAPIBlocking(org, bucket)

	// Create a point
	p := influxdb2.NewPointWithMeasurement("test").
		AddField("value", 1).
		SetTime(time.Now())

	// Write the point
	return writeAPI.WritePoint(context.Background(), p)
}
