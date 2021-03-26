package redis

import (
	redisGo "github.com/gomodule/redigo/redis"
)

type (
	// redis 插入类型
	ListInsertType string
)

const (
	// 在锚定值之前插入
	ListInsertBefore ListInsertType = "before"
	// 在锚定值之后插入
	ListInsertAfter ListInsertType = "after"
)

var ()

// 将一个或多个值 value 插入到列表 key 的表头
// 如果 key 存在，返回 key 列表的元素个数
// 如果 key 不存在，创建未列表类型并执行 LPush 操作
// 如果 key 存在且不为列表类型，返回错误
func LPush(key string, values ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	values = append([]interface{}{key}, values...)

	v, err := redisGo.Int(conn.Do("lpush", values...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 将一个或多个值 value 插入到列表 key 的表头，当 key 列表不在，不执行任何操作
// 如果 key 存在，返回 key 列表的元素个数
// 如果 key 存在且不为列表类型，返回错误
func LPushX(key string, values ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	values = append([]interface{}{key}, values...)

	v, err := redisGo.Int(conn.Do("lpushx", values...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 将一个或多个值 value 插入到列表 key 的表尾
// 如果 key 存在，返回 key 列表的元素个数
// 如果 key 不存在，创建未列表类型并执行 RPush 操作
// 如果 key 存在且不为列表类型，返回错误
func RPush(key string, values ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	values = append([]interface{}{key}, values...)

	v, err := redisGo.Int(conn.Do("rpush", values...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 将一个或多个值 value 插入到列表 key 的表尾，当 key 列表不在，不执行任何操作
// 如果 key 存在，返回 key 列表的元素个数
// 如果 key 存在且不为列表类型，返回错误
func RPushX(key string, values ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	values = append([]interface{}{key}, values...)

	v, err := redisGo.Int(conn.Do("rpushx", values...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 移除并返回列表 key 的头元素
// 如果 key 存在，返回 key 列表的头元素
// 如果 key 不存在，返回 nil 错误
// 如果 key 存在且不为列表类型，返回错误
func LPop(key string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("lpop", key))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 移除并返回列表 key 的尾元素
// 如果 key 存在，返回 key 列表的头元素
// 如果 key 不存在，返回 nil 错误
// 如果 key 存在且不为列表类型，返回错误
func RPop(key string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("rpop", key))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// rpop 列表 source 的尾元素，将其 lpush 到 destination 表头，并将该元素返回给客户端
// 如果 source 不存在，值 nil 被返回，并且不执行其他动作
// 如果 source 存在且不为列表类型，返回错误
func RPopLPush(source, destination string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("rpoplpush", source, destination))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 根据参数 count 的值，移除列表 key 中与参数 value 相等的元素，返回移除掉的个数
// count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count
// count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值
// count = 0 : 移除表中所有与 value 相等的值
// 如果 key 不存在，返回0
// 如果 key 存在且不为列表类型，返回错误
func LRem(key string, count int, value interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("lrem", key, count, value))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回列表 key 的长度
// 如果 key 不存在，则 key 被解释为一个空列表，返回 0
// 如果 key 不是列表类型，返回一个错误
func LLen(key string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("llen", key))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回列表 key 中，下标为 index 的元素,以 0 表示列表的第一个元素
// 以 1 表示列表的第二个元素，以此类推
// 以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
// 如果 key 不存在，则返回 nil 错误
// 如果 key 不是列表类型，返回一个错误
func LIndex(key string, index int) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("lindex", key, index))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 将值 value 插入到列表 key 当中，位于值 pivot 之前或之后
// 返回插入操作完成之后，列表的长度。 如果没有找到 pivot ，返回 -1 。 如果 key 不存在或为空列表，返回 0
// 当 pivot 不存在于列表 key 时，不执行任何操作
// 当 key 不存在时， key 被视为空列表，不执行任何操作
// 如果 key 不是列表类型，返回一个错误
func LInsert(key string, insertType ListInsertType, pivot, value interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("linsert", key, insertType, pivot, value))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 将列表 key 下标为 index 的元素的值设置为 value
// 当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误
// 如果 key 不是列表类型，返回一个错误
func LSet(key string, index int, value interface{}) error {
	conn, err := GetRedisConn()
	if err != nil {
		return err
	}

	_, err = redisGo.String(conn.Do("lset", key, index, value))
	if err != nil {
		return err
	}

	return conn.Close()
}

// 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定，闭区间
// start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推
// 也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
// 如果 key 不存在，返回空字符串组
// 如果 key 存在且不是列表类型，返回一个错误
func LRange(key string, start, end int) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("lrange", key, start, end))
	if err != nil {
		return nil, err
	}

	values := make([]string, 0)
	for _, d := range data {
		values = append(values, string(d))
	}

	return values, conn.Close()
}

// 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除,start 和 end 都会被保留
// start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推
// 也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
// 如果 key 不存在，也返回 OK
// 如果 key 存在且不是列表类型，返回一个错误
func LTrim(key string, start, end int) error {
	conn, err := GetRedisConn()
	if err != nil {
		return err
	}

	_, err = conn.Do("ltrim", key, start, end)
	if err != nil {
		return err
	}

	return conn.Close()
}

// 列表 lpop 的阻塞式(blocking)弹出，当 keys 中的所有键都没有值可供给弹出时，等待 timeout 超时退出
// timeout == 0 ,表示无限等待
// 如果 keys 中存在任何一个元素不是列表类型，返回一个错误
func BLPop(keys []string, timeout int) (map[string]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	args = append(args, timeout)

	data, err := redisGo.ByteSlices(conn.Do("blpop", args...))
	if err != nil {
		return nil, err
	}

	vMap := make(map[string]string)
	for i := 0; i < len(data); i += 2 {
		vMap[string(data[i])] = string(data[i+1])
	}

	return vMap, conn.Close()
}

// 列表 rpop 的阻塞式(blocking)弹出，当 keys 中的所有键都没有值可供给弹出时，等待 timeout 超时退出
// timeout == 0 ,表示无限等待
// 如果 keys 中存在任何一个元素不是列表类型，返回一个错误
func BRPop(keys []string, timeout int) (map[string]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	args = append(args, timeout)

	data, err := redisGo.ByteSlices(conn.Do("brpop", args...))
	if err != nil {
		return nil, err
	}

	vMap := make(map[string]string)
	for i := 0; i < len(data); i += 2 {
		vMap[string(data[i])] = string(data[i+1])
	}

	return vMap, conn.Close()
}

// 列表 RPOPLPUSH 的阻塞式(blocking)版本，当给定列表 source 不为空时， BRPOPLPUSH 的表现和 RPOPLPUSH source destination 一样
// timeout == 0 ,表示无限等待
// 如果 source 存在且不是列表类型，返回一个错误
func BRPopLPush(source, destination string, timeout int) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("brpoplpush", source, destination, timeout))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}
