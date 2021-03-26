package redis

import (
	redisGo "github.com/gomodule/redigo/redis"
	"time"
)

type (
	DistributedLock interface {
		// 初始化
		Initialize() error

		// 加锁
		Lock(string, int) bool

		// 阻塞调用线程直到加锁为止，间隔尝试时间为 0.04s、0.2s、1s、1s、1s......
		UntilLock(string, int) bool

		// 解锁
		Unlock(string) bool

		// 关闭
		Close() error
	}

	redisLock struct {
		conn redisGo.Conn // redis 连接
	}
)

const (
	DefaultEx    = 5  // 默认 ex 过期时间
	DefaultMinEx = 1  // 最小 ex 过期时间
	DefaultMaxEx = 15 // 最大 ex 过期时间
)

var ()

func NewRedisLock() DistributedLock {
	return &redisLock{}
}

func (rl *redisLock) Initialize() error {
	c, err := GetRedisConn()
	if err != nil {
		return err
	}

	rl.conn = c
	return nil
}

// 加锁
// key ： 键
// ex : 过期时间，默认 5 s ，可选闭区间[1,15]
func (rl *redisLock) Lock(key string, ex int) bool {
	// 连接为 nil ， 加锁失败
	if rl.conn == nil {
		return false
	}
	// 键为空，加锁失败
	if key == "" {
		return false
	}
	// 过期时间优化
	if ex < DefaultMinEx || ex > DefaultMaxEx {
		ex = DefaultEx
	}

	_, err := redisGo.String(rl.conn.Do("set", key, 1, "ex", ex, "nx"))
	if err != nil {
		// 无论是 nil return 还是别的错误，都加锁失败
		return false
	}

	// “OK” 加锁成功
	return true
}

func (rl *redisLock) UntilLock(key string, ex int) bool {
	// 连接为 nil ， 加锁失败
	if rl.conn == nil {
		return false
	}
	// 键为空，加锁失败
	if key == "" {
		return false
	}

	wait := time.Millisecond * time.Duration(40)

	for !rl.Lock(key, ex) {
		time.Sleep(wait)

		if wait < time.Second {
			wait = wait * 5
		}
	}

	return true
}

func (rl *redisLock) Unlock(key string) bool {
	// 连接为 nil ， 解锁失败
	if rl.conn == nil {
		return false
	}
	// 键为空，解锁失败
	if key == "" {
		return false
	}

	v, err := redisGo.Int(rl.conn.Do("del", key))
	if err != nil {
		// 错误，都加锁失败
		return false
	}

	// v == 1 说明删除对应的 key
	return v == 1
}

func (rl *redisLock) Close() error {
	return rl.conn.Close()
}
