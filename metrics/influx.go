package metrics

import (
	"context"
	"fmt"
	"time"

	"git.verystar.cn/GaomingQian/gorgeous/provider"
	client "github.com/influxdata/influxdb/client/v2"
)

type InfluxMetricsConfig struct {
	Ctx       context.Context
	Name      string
	Addr      string
	Database  string
	Username  string
	Password  string
	Precision string
	Size      int
	BatchSize int
	Logger    provider.ILogger
}

func (c *InfluxMetricsConfig) Valid() error {
	if c.Size <= 0 {
		return fmt.Errorf("size can't less than 0")
	}

	if c.BatchSize <= 0 {
		return fmt.Errorf("batch size can't less than 0")
	}

	if c.Logger == nil {
		return fmt.Errorf("can't set logger nil")
	}

	if c.Ctx == nil {
		return fmt.Errorf("can't set context nil")
	}

	return nil
}

type InfluxMetrics struct {
	ctx       context.Context
	client    client.Client
	logger    provider.ILogger
	name      string
	batchSize int
	size      int

	point chan *client.Point
	err   chan error

	addr      string
	db        string
	precision string
	username  string
	password  string
}

func NewInfluxMetrics(config *InfluxMetricsConfig) (*InfluxMetrics, error) {
	var err error
	m := new(InfluxMetrics)

	if err = config.Valid(); err != nil {
		return nil, err
	}

	m.ctx = config.Ctx
	m.addr = config.Addr
	m.db = config.Database
	m.precision = config.Precision
	m.username = config.Username
	m.password = config.Password
	m.size = config.Size
	m.batchSize = config.BatchSize
	m.logger = config.Logger
	m.name = config.Name

	m.point = make(chan *client.Point, config.Size)
	m.err = make(chan error, 1)

	if m.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     m.addr,
		Username: m.username,
		Password: m.password,
		Timeout:  5 * time.Second,
	}); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *InfluxMetrics) Measure(tags map[string]string, fields map[string]interface{}) error {
	pt, err := client.NewPoint(m.name, tags, fields, time.Now())
	if err != nil {
		return err
	}
	m.point <- pt
	return nil
}

func (m *InfluxMetrics) Work() chan<- []byte {
	return nil
}

func (m *InfluxMetrics) Next(provider.IWorker) {

}

func (m *InfluxMetrics) Start() {
	for i := 0; i < m.size; i++ {
		go m.run(i)
	}
}

func (m *InfluxMetrics) Stop() {
	m.client.Close()
}

func (m *InfluxMetrics) run(id int) {
	points := make([]*client.Point, 0, m.batchSize)
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-m.ctx.Done():
			if err := m.write(&points); err != nil {
				m.err <- err
			}
			return

		case point := <-m.point:
			m.handleData(&points, point)
		case err := <-m.err:
			m.handleError(err)
		case <-tick.C:
			if err := m.write(&points); err != nil {
				m.err <- err
			}
		}
	}
}

func (m *InfluxMetrics) handleData(points *[]*client.Point, point *client.Point) {
	*points = append(*points, point)
	if len(*points) < m.batchSize {
		return
	}

	if err := m.write(points); err != nil {
		m.err <- err
	}
}

func (m *InfluxMetrics) write(points *[]*client.Point) error {
	if len(*points) == 0 {
		return nil
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  m.db,
		Precision: m.precision,
	})
	if err != nil {
		return fmt.Errorf("can't create batch points: err = %v", err)
	}

	bp.AddPoints(*points)

	if err := m.client.Write(bp); err != nil {
		return fmt.Errorf("can't write batch points: err = %v", err)
	}

	*points = make([]*client.Point, 0, m.batchSize)
	return nil
}

func (m *InfluxMetrics) handleError(err error) {
	m.logger.Errorf("influx metrics got an error: %v", err)
}
