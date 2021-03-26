package redis

import redisGo "github.com/gomodule/redigo/redis"

// 设置指定位置的字符
func Del(keys ...string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}

	v, err := redisGo.Int(conn.Do("del", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}
