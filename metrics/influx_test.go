package metrics

import (
	"context"
	"log"
	"testing"
	"time"

	"git.verystar.cn/golib/factory"
	"git.verystar.cn/golib/logger"
	"git.verystar.cn/golib/utility"
)

func TestNewInfluxMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m, err := factory.CreateInfluxMetrics(
		ctx,
		"http://localhost:8086",
		"wemediareply",
		"wemediareply",
		"wemediareply",
		"us",
		2,
		10,
		logger.NewStdLogger(),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100000; i++ {
		m.Measure("test_lib", map[string]string{"service": "redis_fetcher"}, map[string]interface{}{
			"count":    utility.RandInt(1, 10),
			"duration": utility.RandFloat64(20, 100),
		})
		time.Sleep(time.Duration(utility.RandInt(1, 100)) * time.Microsecond)
	}
	time.Sleep(2 * time.Second)
	cancel()
}
