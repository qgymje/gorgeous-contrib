package fetcher

import (
	"log"
	"os"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type NsqFetcher struct {
	name       string
	size       int
	decodeFunc DecodeFunc
	data       chan []byte
	consumers  []*nsqConsumer
}

func NewNsqFetcher(name string, size int, nsqlookupdAddrs []string, topic, channel string, decodeFunc DecodeFunc) (*NsqFetcher, error) {
	nf := &NsqFetcher{
		name:       name,
		size:       size,
		decodeFunc: decodeFunc,
		data:       make(chan []byte),
	}

	nf.consumers = make([]*nsqConsumer, 0, size)
	for i := 0; i < size; i++ {
		consumer, err := newNsqConsumer(topic, channel, nsqlookupdAddrs)
		if err != nil {
			return nil, err
		}
		nf.consumers = append(nf.consumers, consumer)
		go func(c *nsqConsumer) {
			nf.run(c)
		}(consumer)
	}

	return nf, nil
}

func (nf *NsqFetcher) run(consumer *nsqConsumer) {
	for v := range consumer.data {
		nf.data <- v
	}
}

func (nf *NsqFetcher) Name() string {
	return nf.name
}

func (nf *NsqFetcher) Size() int {
	return nf.size
}

func (nf *NsqFetcher) Action() (interface{}, error) {
	return nf.decodeFunc(<-nf.data)
}

func (nf *NsqFetcher) Interval() time.Duration {
	return 0
}

func (nf *NsqFetcher) Close() error {
	for _, c := range nf.consumers {
		c.StopConsume()
	}
	return nil
}

// nsqConsumer is a single consumer instance
type nsqConsumer struct {
	topic, channel string
	data           chan []byte
	consumer       *nsq.Consumer
}

func newNsqConsumer(topic, channel string, lookupdAddress []string) (*nsqConsumer, error) {
	var err error
	c := &nsqConsumer{
		topic:   topic,
		channel: channel,
		data:    make(chan []byte),
	}

	config := nsq.NewConfig()
	c.consumer, err = nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}
	c.consumer.SetLogger(log.New(os.Stderr, "", log.Flags()), nsq.LogLevelError)
	c.consumer.AddConcurrentHandlers(c, 1)

	err = c.consumer.ConnectToNSQLookupds(lookupdAddress)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *nsqConsumer) StopConsume() {
	c.consumer.Stop()
}

func (c *nsqConsumer) HandleMessage(msg *nsq.Message) error {
	c.data <- msg.Body
	return nil
}
