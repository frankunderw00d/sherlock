package redis

import (
	redisGo "github.com/gomodule/redigo/redis"
)

// 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略
// 假如 key 不存在，则创建一个只包含 member 元素作成员的集合
// 当 key 不是集合类型时，返回一个错误
func SAdd(key string, members ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{key}
	for _, member := range members {
		args = append(args, member)
	}

	v, err := redisGo.Int(conn.Do("sadd", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 判断 member 元素是否集合 key 的成员
// 如果 member 元素是集合的成员，返回 1
// 如果 member 元素不是集合的成员，或 key 不存在，返回 0
// 当 key 不是集合类型时，返回一个错误
func SIsMember(key string, member interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("sismember", key, member))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 移除并返回集合 key 中的 count 个随机元素
// 返回被移除的随机元素 , 如果 count 大于集合长度，则将集合内所有元素返回，以实际长度为准
// 当 key 不存在或 key 是空集时，返回 nil
// 当 key 不是集合类型时，返回一个错误
func SPop(key string, count int) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("spop", key, count))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回集合 key 中的 count 个随机元素，不移除
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同
// 如果 count 大于等于集合基数，那么返回整个集合
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值
// 如果 count == 0，那么返回空集合
// 当 key 不是集合类型时，返回一个错误
func SRandMember(key string, count int) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("srandmember", key, count))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同
// 当 key 不存在，返回 0
// 当 key 存在且不是集合类型时，返回一个错误
func SRem(key string, members ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{key}
	args = append(args, members...)

	v, err := redisGo.Int(conn.Do("srem", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 将 member 元素从 source 集合移动到 destination 集合
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，返回 0
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除
// 当 source 或 destination 不是集合类型时，返回一个错误
func SMove(source, destination string, member interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("smove", source, destination, member))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回集合 key 的基数(集合中元素的数量)
// 当 key 不存在时，返回 0
// 当 key 存在且不是集合类型时，返回一个错误
func SCard(key string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("scard", key))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回集合 key 的基数(集合中元素的数量)
// 当 key 不存在时，返回空集合
// 当 key 存在且不是集合类型时，返回一个错误
func SMembers(key string) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	data, err := redisGo.ByteSlices(conn.Do("smembers", key))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回一个集合的全部成员，该集合是所有给定集合 keys 的交集
// 不存在的 key 被视为空集,当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)
// 当 keys 中存在不是集合类型时，返回一个错误
func SInter(keys ...string) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}

	data, err := redisGo.ByteSlices(conn.Do("sinter", args...))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回一个集合的全部成员，该集合是所有给定集合 keys 的并集
// 不存在的 key 被视为空集,当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)
// 当 keys 中存在不是集合类型时，返回一个错误
func SUnion(keys ...string) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}

	data, err := redisGo.ByteSlices(conn.Do("sunion", args...))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回一个集合的全部成员，该集合是所有给定集合 keys 的差集，以第一个 key 为基准
// 不存在的 key 被视为空集,当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)
// 当 keys 中存在不是集合类型时，返回一个错误
func SDiff(keys ...string) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}

	data, err := redisGo.ByteSlices(conn.Do("sdiff", args...))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 返回一个集合的全部成员，该集合是所有给定集合 keys 的交集，但它将结果保存到 destination 集合
// 不存在的 key 被视为空集,当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)，空集不会保存
// 当 keys 中存在不是集合类型时，返回一个错误
func SInterStore(destination string, keys ...string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{destination}
	for _, key := range keys {
		args = append(args, key)
	}

	v, err := redisGo.Int(conn.Do("sinterstore", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回一个集合的全部成员，该集合是所有给定集合 keys 的并集，但它将结果保存到 destination 集合
// 不存在的 key 被视为空集,当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)，空集不会保存
// 当 keys 中存在不是集合类型时，返回一个错误
func SUnionStore(destination string, keys ...string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{destination}
	for _, key := range keys {
		args = append(args, key)
	}

	v, err := redisGo.Int(conn.Do("sunionstore", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回一个集合的全部成员，该集合是所有给定集合 keys 的差集，但它将结果保存到 destination 集合
// 不存在的 key 被视为空集,当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)，空集不会保存
// 当 keys 中存在不是集合类型时，返回一个错误
func SDiffStore(destination string, keys ...string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{destination}
	for _, key := range keys {
		args = append(args, key)
	}

	v, err := redisGo.Int(conn.Do("sdiffstore", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}
