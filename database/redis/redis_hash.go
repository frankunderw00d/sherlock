package redis

import (
	"errors"
	redisGo "github.com/gomodule/redigo/redis"
)

// 将哈希表 key 中域 field 的值设置为 value
// 如果 key 并不存在， 那么一个新的哈希表将被创建并执行 HSET 操作
// 如果域 field 已经存在于哈希表中， 那么它的旧值将被新值 value 覆盖
// 当 HSET 命令在哈希表中新创建 field 域并成功为它设置值时， 命令返回 1
// 如果域 field 已经存在于哈希表， 并且 HSET 命令成功使用新值覆盖了它的旧值， 那么命令返回 0
func HSet(key string, field string, value interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("hset", key, field, value))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 当且仅当域 field 尚未存在于哈希表 key 的情况下将哈希表 key 中域 field 的值设置为 value
// 如果 key 并不存在， 那么一个新的哈希表将被创建并执行 HSETNX 操作
// 当 HSET 命令在哈希表中新创建 field 域并成功为它设置值时， 命令返回 1
// 如果域 field 已经存在于哈希表， 并且 HSET 命令成功使用新值覆盖了它的旧值， 那么命令返回 0
func HSetNX(key string, field string, value interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("hsetnx", key, field, value))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回哈希表 key 中给定域 field 的值
// 如果 key 或 field 并不存在， 那么返回 nil 错误
func HGet(key string, field string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("hget", key, field))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 检查给定域 field 是否存在于哈希表 key 当中
// 如果 key 或 field 并不存在， 那么返回 0
func HExists(key string, field string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("hexists", key, field))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略
// 如果 key 或 fields 并不存在， 那么返回 0
// 返回删除的 field 个数
func HDel(key string, fields ...string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{key}
	for _, field := range fields {
		args = append(args, field)
	}

	v, err := redisGo.Int(conn.Do("hdel", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回哈希表 key 中域的数量
// 如果 key 或 fields 并不存在， 那么返回 0
func HLen(key string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("hlen", key))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回哈希表 key 中， 与给定域 field 相关联的值的字符串长度
// 如果 key 或 fields 并不存在， 那么返回 0
func HStrLen(key, field string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("hstrlen", key, field))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 为哈希表 key 中的域 field 的值加上增量 increment , increment 必须为整数
// 如果 key 或 field 不存在，一个新的哈希表或域被创建并执行 HINCRBY 命令
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误
// 返回 key 中该 field 加上增量后的值
func HIncrBy(key, field string, increment int) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("hincrby", key, field, increment))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 为哈希表 key 中的域 field 的值加上增量 increment , increment 必须为浮点数
// 如果 key 或 field 不存在，一个新的哈希表或域被创建并执行 HINCRBY 命令
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误
// 返回 key 中该 field 加上增量后的值
// (注意 : 原本的整数可以通过 hincrbyfloat 变成浮点数，但是此过程不可逆)
func HIncrByFloat(key, field string, increment float64) (float64, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Float64(conn.Do("hincrbyfloat", key, field, increment))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 将哈希表 key 中域 field 的值设置为 value
// 如果 key 并不存在， 那么一个新的哈希表将被创建并执行 HSET 操作
// 如果域 field 已经存在于哈希表中， 那么它的旧值将被新值 value 覆盖
// 如果 key 不是哈希表类型，返回错误
func HMSet(key string, fv map[string]interface{}) error {
	conn, err := GetRedisConn()
	if err != nil {
		return err
	}

	args := []interface{}{key}
	for field, value := range fv {
		args = append(args, []interface{}{field, value}...)
	}

	_, err = conn.Do("hmset", args...)
	if err != nil {
		return err
	}

	return conn.Close()
}

// 返回哈希表 key 中，一个或多个给定域的值
// 如果给定的域不存在于哈希表，那么返回一个 nil 值
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表
func HMGet(key string, fields ...string) (map[string]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := []interface{}{key}
	for _, field := range fields {
		args = append(args, field)
	}

	data, err := redisGo.ByteSlices(conn.Do("hmget", args...))
	if err != nil {
		return nil, err
	}

	if len(data) != len(fields) {
		return nil, errors.New("values count not equal fields count")
	}

	vMap := make(map[string]string)
	for index, d := range data {
		vMap[fields[index]] = string(d)
	}

	return vMap, conn.Close()
}

// 返回哈希表 key 中的所有域
// 当 key 不存在时，返回一个空表
// 当 key 存在且不为哈希表类型时，返回错误
func HKeys(key string) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("hkeys", key))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回哈希表 key 中所有域的值
// 当 key 不存在时，返回一个空表
// 当 key 存在且不为哈希表类型时，返回错误
func HVals(key string) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("hvals", key))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回哈希表 key 中所有域的值
// 当 key 不存在时，返回一个空表
// 当 key 存在且不为哈希表类型时，返回错误
func HGetAll(key string) (map[string]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("hgetall", key))
	if err != nil {
		return nil, err
	}

	v := make(map[string]string, 0)
	for i := 0; i < len(data); i += 2 {
		v[string(data[i])] = string(data[i+1])
	}

	return v, conn.Close()
}
