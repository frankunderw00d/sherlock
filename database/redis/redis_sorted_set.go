package redis

import (
	redisGo "github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
)

type (
	// redis 有序集合成员
	SortedSetMember struct {
		Member string  `json:"member"`
		Score  float64 `json:"score"`
	}

	// redis 有序集合中带分数展示类型
	SortedSetScoresType string

	// redis 有序集合中展示个数限制类型
	SortedSetLimit string

	// redis 有序集合中最小最大值类型
	SortedSetValue string

	// redis 有序集合中 LEX 最小最大值类型
	SortedSetLexValue string
)

const (
	// redis 有序集合中带分数展示
	SortedSetWithScores SortedSetScoresType = "withscores"
	// redis 有序集合中不带分数展示
	SortedSetWithOutScores SortedSetScoresType = ""
	// redis 有序集合中不限制展示个数
	SortedSetLimitNone SortedSetLimit = "NOLIMIT"
	// redis 有序集合中最大值 +inf
	SortedSetValuePositiveInf SortedSetValue = "+inf"
	// redis 有序集合中最小值 -inf
	SortedSetValueNegativeInf SortedSetValue = "-inf"
	// redis 有序集合中 LEX 最大值 +
	SortedSetValueMax SortedSetLexValue = "+"
	// redis 有序集合中最小值 -inf
	SortedSetValueMin SortedSetLexValue = "-"
)

var ()

// 解析 redis 有序集合中限制展示个数命令
func (rssl SortedSetLimit) parse() []interface{} {
	if rssl == SortedSetLimitNone {
		return []interface{}{}
	}

	orders := strings.Split(string(rssl), " ")
	v := make([]interface{}, 0)
	for _, o := range orders {
		v = append(v, o)
	}
	return v
}

func Limit(s string) SortedSetLimit {
	return SortedSetLimit(s)
}

func (rssv SortedSetValue) parse() interface{} {
	if rssv == SortedSetValuePositiveInf {
		return string(SortedSetValuePositiveInf)
	}
	if rssv == SortedSetValueNegativeInf {
		return string(SortedSetValueNegativeInf)
	}

	v, err := strconv.ParseFloat(string(rssv), 64)
	if err != nil {
		return 0
	}
	return v
}

func Value(s string) SortedSetValue {
	return SortedSetValue(s)
}

func (rsslv SortedSetLexValue) parse() interface{} {
	return string(rsslv)
}

func LEXValue(s string) SortedSetLexValue {
	return SortedSetLexValue(s)
}

// 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
// score 值可以是整数值或双精度浮点数
// 如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作
// 当 key 不是有序集合类型时，返回一个错误
func ZAdd(key string, sm map[string]float64) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	args := []interface{}{key}

	for member, score := range sm {
		args = append(args, []interface{}{score, member}...)
	}

	v, err := redisGo.Int(conn.Do("zadd", args...))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 返回有序集 key 中，成员 member 的 score 值
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil
// 当 key 不是有序集合类型时，返回一个错误
func ZScore(key string, member string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("zscore", key, member))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 为有序集 key 的成员 member 的 score 值加上增量 increment
// 通过传递一个负数值 increment ，让 score 减去相应的值
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member
// 当 key 不是有序集合类型时，返回一个错误
func ZIncrBy(key string, increment float64, member string) (string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return "", err
	}

	v, err := redisGo.String(conn.Do("zincrby", key, increment, member))
	if err != nil {
		return "", err
	}

	return v, conn.Close()
}

// 返回有序集 key 的基数
// 当 key 不存在时，返回 0
// 当 key 不是有序集合类型时，返回一个错误
func ZCard(key string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("zcard", key))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量
// 当 key 不存在时，返回 0
// 当 key 不是有序集合类型时，返回一个错误
func ZCount(key string, min, max float64) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	v, err := redisGo.Int(conn.Do("zcount", key, min, max))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 返回有序集 key 中，指定区间内的成员,其中成员的位置按 score 值递增(从小到大)来排序
// 具有相同 score 值的成员按字典序(lexicographical order)来排列
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推
// 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
// 当 key 不存在时，返回 0
// 当 key 不是有序集合类型时，返回一个错误
func ZRange(key string, start, stop int, scoresType SortedSetScoresType) ([]SortedSetMember, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := []interface{}{key, start, stop}
	if scoresType == SortedSetWithScores {
		args = append(args, scoresType)
	}

	data, err := redisGo.ByteSlices(conn.Do("zrange", args...))
	if err != nil {
		return nil, err
	}

	v := make([]SortedSetMember, 0)
	if scoresType == SortedSetWithScores {
		for i := 0; i < len(data); i += 2 {
			score, err := strconv.ParseFloat(string(data[i+1]), 64)
			if err != nil {
				return nil, err
			}

			member := SortedSetMember{
				Member: string(data[i]),
				Score:  score,
			}

			v = append(v, member)
		}
	} else {
		for i := 0; i < len(data); i++ {
			v = append(v, SortedSetMember{Member: string(data[i])})
		}
	}

	return v, conn.Close()
}

// 返回有序集 key 中，指定区间内的成员,其中成员的位置按 score 值递增(从大到小)来排序
// 具有相同 score 值的成员按字典序反序(reverse lexicographical order)来排列
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推
// 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
// 当 key 不存在时，返回 0
// 当 key 不是有序集合类型时，返回一个错误
func ZRevRange(key string, start, stop int, scoresType SortedSetScoresType) ([]SortedSetMember, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := []interface{}{key, start, stop}
	if scoresType == SortedSetWithScores {
		args = append(args, scoresType)
	}

	data, err := redisGo.ByteSlices(conn.Do("zrevrange", args...))
	if err != nil {
		return nil, err
	}

	v := make([]SortedSetMember, 0)
	if scoresType == SortedSetWithScores {
		for i := 0; i < len(data); i += 2 {
			score, err := strconv.ParseFloat(string(data[i+1]), 64)
			if err != nil {
				return nil, err
			}

			member := SortedSetMember{
				Member: string(data[i]),
				Score:  score,
			}

			v = append(v, member)
		}
	} else {
		for i := 0; i < len(data); i++ {
			v = append(v, SortedSetMember{Member: string(data[i])})
		}
	}

	return v, conn.Close()
}

// 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列
// 具有相同 score 值的成员按字典序(lexicographical order)来排列
// 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count)
// 可选的 WITHSCORES 参数决定结果集是单单返回有序集的成员，还是将有序集成员及其 score 值一起返回
// 当 key 不存在时，返回 0
// 当 key 不是有序集合类型时，返回一个错误
func ZRangeByScore(key string, min, max SortedSetValue, scoresType SortedSetScoresType, limit SortedSetLimit) ([]SortedSetMember, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := []interface{}{key, min.parse(), max.parse()}
	if scoresType == SortedSetWithScores {
		args = append(args, scoresType)
	}
	l := limit.parse()
	if len(l) > 0 {
		args = append(args, l...)
	}

	data, err := redisGo.ByteSlices(conn.Do("zrangebyscore", args...))
	if err != nil {
		return nil, err
	}

	v := make([]SortedSetMember, 0)
	if scoresType == SortedSetWithScores {
		for i := 0; i < len(data); i += 2 {
			score, err := strconv.ParseFloat(string(data[i+1]), 64)
			if err != nil {
				return nil, err
			}

			member := SortedSetMember{
				Member: string(data[i]),
				Score:  score,
			}

			v = append(v, member)
		}
	} else {
		for i := 0; i < len(data); i++ {
			v = append(v, SortedSetMember{Member: string(data[i])})
		}
	}

	return v, conn.Close()
}

// 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列
// 具有相同 score 值的成员按字典序(lexicographical order)来排列
// 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count)
// 可选的 WITHSCORES 参数决定结果集是单单返回有序集的成员，还是将有序集成员及其 score 值一起返回
// 当 key 不存在时，返回 0
// 当 key 不是有序集合类型时，返回一个错误
func ZRevRangeByScore(key string, max, min SortedSetValue, scoresType SortedSetScoresType, limit SortedSetLimit) ([]SortedSetMember, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := []interface{}{key, max.parse(), min.parse()}
	if scoresType == SortedSetWithScores {
		args = append(args, scoresType)
	}
	l := limit.parse()
	if len(l) > 0 {
		args = append(args, l...)
	}

	data, err := redisGo.ByteSlices(conn.Do("zrevrangebyscore", args...))
	if err != nil {
		return nil, err
	}

	v := make([]SortedSetMember, 0)
	if scoresType == SortedSetWithScores {
		for i := 0; i < len(data); i += 2 {
			score, err := strconv.ParseFloat(string(data[i+1]), 64)
			if err != nil {
				return nil, err
			}

			member := SortedSetMember{
				Member: string(data[i]),
				Score:  score,
			}

			v = append(v, member)
		}
	} else {
		for i := 0; i < len(data); i++ {
			v = append(v, SortedSetMember{Member: string(data[i])})
		}
	}

	return v, conn.Close()
}

// 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0
// 使用 ZREVRANK key member 命令可以获得成员按 score 值递减(从大到小)排列的排名
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRank(key string, member interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	v, err := redisGo.Int(conn.Do("zrank", key, member))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0
// 使用 ZRANK key member 命令可以获得成员按 score 值递增(从小到大)排列的排名
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRevRank(key string, member interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	v, err := redisGo.Int(conn.Do("zrevrank", key, member))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRem(key string, members ...interface{}) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	args := []interface{}{key}
	for _, member := range members {
		args = append(args, member)
	}

	v, err := redisGo.Int(conn.Do("zrem", args...))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 移除有序集 key 中，指定排名(rank)区间内的所有成员
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推
// 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRemRangeByRank(key string, start, stop int) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	v, err := redisGo.Int(conn.Do("zremrangebyrank", key, start, stop))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max)的成员
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推
// 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRemRangeByScore(key string, min, max SortedSetValue) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	v, err := redisGo.Int(conn.Do("zremrangebyscore", key, min.parse(), max.parse()))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 当有序集合的所有成员都具有相同的分值时， 有序集合的元素会根据成员的字典序（lexicographical ordering）来进行排序
// 而这个命令则可以返回给定的有序集合键 key 中， 值介于 min 和 max 之间的成员
// 合法的 min 和 max 参数必须包含 ( 或者 [ ， 其中 ( 表示开区间（指定的值不会被包含在范围之内）， 而 [ 则表示闭区间（指定的值会被包含在范围之内）
// 特殊值 + 和 - 在 min 参数以及 max 参数中具有特殊的意义， 其中 + 表示正无限， 而 - 表示负无限
// // 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count)
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRangeByLex(key string, min, max SortedSetLexValue, limit SortedSetLimit) ([]string, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}

	args := []interface{}{key, min.parse(), max.parse()}
	l := limit.parse()
	if len(l) > 0 {
		args = append(args, l...)
	}

	data, err := redisGo.ByteSlices(conn.Do("zrangebylex", args...))
	if err != nil {
		return nil, err
	}

	v := make([]string, 0)
	for _, d := range data {
		v = append(v, string(d))
	}

	return v, conn.Close()
}

// 对于一个所有成员的分值都相同的有序集合键 key 来说， 这个命令会返回该集合中， 成员介于 min 和 max 范围内的元素数量
// 合法的 min 和 max 参数必须包含 ( 或者 [ ， 其中 ( 表示开区间（指定的值不会被包含在范围之内）， 而 [ 则表示闭区间（指定的值会被包含在范围之内）
// 特殊值 + 和 - 在 min 参数以及 max 参数中具有特殊的意义， 其中 + 表示正无限， 而 - 表示负无限
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZLexCount(key string, min, max SortedSetLexValue) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	v, err := redisGo.Int(conn.Do("zlexcount", key, min.parse(), max.parse()))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 对于一个所有成员的分值都相同的有序集合键 key 来说， 这个命令会移除该集合中， 成员介于 min 和 max 范围内的所有元素
// 合法的 min 和 max 参数必须包含 ( 或者 [ ， 其中 ( 表示开区间（指定的值不会被包含在范围之内）， 而 [ 则表示闭区间（指定的值会被包含在范围之内）
// 特殊值 + 和 - 在 min 参数以及 max 参数中具有特殊的意义， 其中 + 表示正无限， 而 - 表示负无限
// 当 key 不存在时，返回 nil 错误
// 当 key 不是有序集合类型时，返回一个错误
func ZRemRangeByLex(key string, min, max SortedSetLexValue) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return -1, err
	}

	v, err := redisGo.Int(conn.Do("zremrangebylex", key, min.parse(), max.parse()))
	if err != nil {
		return -1, err
	}

	return v, conn.Close()
}

// 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination
// 使用 WEIGHTS 选项，你可以为 每个 给定有序集 分别 指定一个乘法因子(multiplication factor)
// 每个给定有序集的所有成员的 score 值在传递给聚合函数(aggregation function)之前都要先乘以该有序集的因子
// 使用 AGGREGATE 选项，你可以指定并集的结果集的聚合方式:
// 参数 SUM ，可以将所有集合中某个成员的 score 值之 和 作为结果集中该成员的 score 值 ， 默认参数
// 参数 MIN ，可以将所有集合中某个成员的 最小 score 值作为结果集中该成员的 score 值
// 参数 MAX 则是将所有集合中某个成员的 最大 score 值作为结果集中该成员的 score 值
// 当 key 不是有序集合类型时，返回一个错误
func ZUnionStore(destination string, keys []string, weights []float64, aggregate string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{destination, len(keys)}
	for _, key := range keys {
		args = append(args, key)
	}

	if len(weights) > 0 {
		args = append(args, "weights")
		for _, weight := range weights {
			args = append(args, weight)
		}
	}

	if aggregate != "" {
		args = append(args, "aggregate")
		args = append(args, aggregate)
	}

	v, err := redisGo.Int(conn.Do("zunionstore", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}

// 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination
// 使用 WEIGHTS 选项，你可以为 每个 给定有序集 分别 指定一个乘法因子(multiplication factor)
// 每个给定有序集的所有成员的 score 值在传递给聚合函数(aggregation function)之前都要先乘以该有序集的因子
// 使用 AGGREGATE 选项，你可以指定并集的结果集的聚合方式:
// 参数 SUM ，可以将所有集合中某个成员的 score 值之 和 作为结果集中该成员的 score 值 ， 默认参数
// 参数 MIN ，可以将所有集合中某个成员的 最小 score 值作为结果集中该成员的 score 值
// 参数 MAX 则是将所有集合中某个成员的 最大 score 值作为结果集中该成员的 score 值
// 当 key 不是有序集合类型时，返回一个错误
func ZInterStore(destination string, keys []string, weights []float64, aggregate string) (int, error) {
	conn, err := GetRedisConn()
	if err != nil {
		return 0, err
	}

	args := []interface{}{destination, len(keys)}
	for _, key := range keys {
		args = append(args, key)
	}

	if len(weights) > 0 {
		args = append(args, "weights")
		for _, weight := range weights {
			args = append(args, weight)
		}
	}

	if aggregate != "" {
		args = append(args, "aggregate")
		args = append(args, aggregate)
	}

	v, err := redisGo.Int(conn.Do("zinterstore", args...))
	if err != nil {
		return 0, err
	}

	return v, conn.Close()
}
