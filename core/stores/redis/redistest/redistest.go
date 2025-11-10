package redistest

import (
	"time"

	"github.com/alicebob/miniredis/v2"
	"go-kit/core/lang"
	"go-kit/core/stores/redis"
)

// CreateRedis returns an in process redis.Redis.
func CreateRedis() (r *redis.Redis, clean func(), err error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	return redis.New(mr.Addr()), func() {
		ch := make(chan lang.PlaceholderType)

		go func() {
			mr.Close()
			close(ch)
		}()

		select {
		case <-ch:
		case <-time.After(time.Second):
		}
	}, nil
}
