package redis

import (
	"errors"
	"fmt"
	redisGo "github.com/gomodule/redigo/redis"
	"time"
)

type (
	// Redis 定义
	Redis interface {
		// 初始化
		Initialize(idleTimeout time.Duration, maxIdle, maxActive int, host string, port int, password string)

		// 获取连接
		Get() (redisGo.Conn, error)

		// 关闭连接池
		Close() error
	}

	// Redis 定义实现
	redis struct {
		pool *redisGo.Pool
	}
)

const (
	// 默认网络协议
	DefaultNetwork = "tcp"
	// 默认连接空闲超时时间 2分钟
	DefaultRedisIdleTimeout = time.Duration(120) * time.Second
	// 默认最大空闲连接数量
	DefaultRedisMaxIdleConn = 10
	// 默认最大连接数
	DefaultRedisMaxActiveConn = 30
	// 默认 host
	DefaultRedisHost = "localhost"
	// 默认 port
	DefaultRedisPort = 6379
)

const (
	ErrNilPoolText = "redis pool is nil"
)

var (
	// redis pool 为 nil 错误
	ErrNilPool = errors.New(ErrNilPoolText)
)

// 新建 Redis
func NewRedis() Redis {
	return &redis{}
}

// 初始化
func (r *redis) Initialize(idleTimeout time.Duration, maxIdle, maxActive int, host string, port int, password string) {
	if idleTimeout == 0 {
		idleTimeout = DefaultRedisIdleTimeout
	}
	if maxIdle == 0 {
		maxIdle = DefaultRedisMaxIdleConn
	}
	if maxActive == 0 {
		maxActive = DefaultRedisMaxActiveConn
	}
	if host == "" {
		host = DefaultRedisHost
	}
	if port == 0 {
		port = DefaultRedisPort
	}

	r.pool = &redisGo.Pool{
		Dial: func() (redisGo.Conn, error) {
			c, err := redisGo.Dial(DefaultNetwork, fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					return nil, c.Close()
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redisGo.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
}

// 获取连接
func (r *redis) Get() (redisGo.Conn, error) {
	if r.pool == nil {
		return nil, ErrNilPool
	}

	return r.pool.Get(), nil
}

// 关闭连接池
func (r *redis) Close() error {
	return r.pool.Close()
}

// --------------------------------------------------- Redis Public Methods --------------------------------------------
var (
	// 默认的 Redis 实例
	defaultRedis = NewRedis()
)

// 1.初始化 Redis
func InitializeRedis(idleTimeout time.Duration, maxIdle, maxActive int, host string, port int, password string) {
	defaultRedis.Initialize(idleTimeout, maxIdle, maxActive, host, port, password)
}

// 2.Redis 获取连接
func GetRedisConn() (redisGo.Conn, error) {
	return defaultRedis.Get()
}

// 3.关闭 Redis
func CloseRedis() error {
	return defaultRedis.Close()
}
