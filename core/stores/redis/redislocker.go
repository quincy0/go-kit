package redis

import (
	"context"
	"math/rand"
	"time"

	"github.com/quincy0/go-kit/core/logx"
)

type Locker interface {
	Lock(key string) (LockerResource, bool)
	SpinLock(key string) (LockerResource, bool)
}

type LockerResource interface {
	Unlock()
}

type locker1 struct {
	redis *Redis
}

type lockerResource struct {
	redis *Redis
	key   string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewLocker(rd *Redis) Locker {
	return &locker1{
		redis: rd,
	}
}

func (l *locker1) Lock(key string) (LockerResource, bool) {
	ok, err := l.redis.SetnxExCtx(context.Background(), key, "1", 30)
	if !ok || (err != Nil && err != nil) {
		return nil, false
	}

	return &lockerResource{
		redis: l.redis,
		key:   key,
	}, true
}

func (l *locker1) SpinLock(key string) (LockerResource, bool) {
	sNum := 3
	var rs LockerResource
	var ok bool
	for {
		rs, ok = l.Lock(key)
		if !ok && sNum >= 0 {
			sNum--
			time.Sleep(time.Millisecond * 30)
		} else {
			break
		}
	}
	if !ok {
		return nil, false
	}

	return rs, sNum >= 0
}

func (l *lockerResource) Unlock() {
	_, err := l.redis.Del(l.key)
	if err != nil {
		logx.Errorw("l.redis.Del Err", logx.Field("err", err))
	}
}
