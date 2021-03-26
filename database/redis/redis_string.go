package redis

import (
	"errors"
	redisGo "github.com/gomodule/redigo/redis"
)

type (
	// 过期时间类型
	TimeoutType string
	// 对于键存在的要求
	KeyExistType string
)

const (
	// 以秒为过期时间
	StringSetTimeoutEX TimeoutType = "ex"
	// 以毫秒为过期时间
	StringSetTimeoutPX TimeoutType = "px"
	// 键必须不存在
	StringSetNotExist KeyExistType = "nx"
	// 键必须存在
	StringSetExist KeyExistType = "xx"
)

var ()

// 单一键全面设置
func Set(key, value string, timeoutType TimeoutType, timeoutValue int64, existType KeyExistType) error {
	conn, err := GetRedisConn()
	if err != nil {
		return err
	}

	// 成功 "OK"
	// 失败 nil
	_, err = redisGo.String(conn.Do("set", key, value, timeoutType, timeoutValue, existType))
	if err != nil {
		return err
	}

	return conn.Close()
}

// 获取值
func Get(key string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	// 成功 "OK"
	// 失败 nil
	v, err := redisGo.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 批量设置，减少网络耗时
func MSet(v map[string]interface{}) error {
	if v == nil || len(v) == 0 {
		return errors.New("value can't be nil or empty")
	}
	conn, err := GetRedisConn()
	if err != nil {
		return err
	}

	args := make([]interface{}, 0)
	for key, value := range v {
		args = append(args, key, value)
	}

	// 成功 "OK"
	_, err = redisGo.String(conn.Do("mset", args...))
	if err != nil {
		return err
	}

	return conn.Close()
}

// 批量获取，减少网络耗时
func MGet(v []string) (map[string]string, error) {
	if v == nil || len(v) == 0 {
		return nil, errors.New("value can't be nil or empty")
	}
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range v {
		args = append(args, key)
	}

	data, err := redisGo.ByteSlices(conn.Do("mget", args...))
	if err != nil {
		return nil, err
	}

	if len(data) != len(v) {
		return nil, errors.New("values count not equal to keys count")
	}

	r := make(map[string]string)
	for index, d := range data {
		r[v[index]] = string(d)
	}

	return r, conn.Close()
}

// 对键做自增，键的值必须为整数
// 键不存在，默认为0，自增为1后返回
func Incr(key string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("incr", key))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 追加值，返回追加过后的字符串长度
func Append(key string, value string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("append", key, value))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回字符串长度，键不存在返回 0
func StrLen(key string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("strlen", key))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 设置并返回原值
func GetSet(key, value string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("getset", key, value))
	if err != nil {
		if err == redisGo.ErrNil {
			return "", nil
		}
		return "", err
	}

	return v, conn.Close()
}

// 设置指定位置的字符
func SetRange(key string, offset int, value string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("setrange", key, offset, value))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 设置指定位置的字符
func GetRange(key string, start, end int) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("getrange", key, start, end))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}
