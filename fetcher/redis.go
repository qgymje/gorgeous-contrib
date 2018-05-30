package fetcher

import (
	"time"

	"git.verystar.cn/golib/utility"
	"github.com/mediocregopher/radix.v2/pool"
)

type DecodeFunc func([]byte) (interface{}, error)

type RedisFetcher struct {
	rdp        *pool.Pool
	name       string
	size       int
	key        string
	cmd        string
	decodeFunc DecodeFunc
}

func NewRedisFetcher(name string, size int, addr, auth, key, cmd string, decodeFunc DecodeFunc) (*RedisFetcher, error) {
	rdp, err := utility.CreateRedisPool(addr, auth, size)
	if err != nil {
		return nil, err
	}

	return &RedisFetcher{
		rdp:        rdp,
		name:       name,
		size:       size,
		key:        key,
		cmd:        cmd,
		decodeFunc: decodeFunc,
	}, nil
}

func (rf *RedisFetcher) Name() string {
	return rf.name
}

func (rf *RedisFetcher) Size() int {
	return rf.size
}

func (rf *RedisFetcher) Action() (interface{}, error) {
	data, err := rf.rdp.Cmd(rf.cmd, rf.key).Bytes()
	if err != nil || len(data) == 0 {
		time.Sleep(2e9)
		return rf.Action()
	}
	return rf.decodeFunc(data)
}

func (rf *RedisFetcher) Interval() time.Duration {
	return 0
}

func (rf *RedisFetcher) Close() error {
	rf.rdp.Empty()
	return nil
}
