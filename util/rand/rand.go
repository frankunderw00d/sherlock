package rand

import (
	"crypto/rand"
	"math/big"
)

type (
	Seed string
)

const (
	// lower case letter
	SeedLCL Seed = "abcdefghijklmnopqrstuvwxyz"
	// upper case letter
	SeedUCL Seed = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// number
	SeedNum Seed = "0123456789"
	// special case char
	SeedSCC Seed = "!@#$%^&*(){}:,./"
	// default seed
	DefaultSeed = SeedLCL + SeedUCL + SeedNum
)

var ()

func init() {}

func baseRand(max int64) int64 {
	p, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0
	}
	return p.Int64()
}

func Int(max int) int {
	return int(baseRand(int64(max)))
}

func Int8(max int8) int8 {
	return int8(baseRand(int64(max)))
}

func Int16(max int16) int16 {
	return int16(baseRand(int64(max)))
}

func Int32(max int32) int32 {
	return int32(baseRand(int64(max)))
}

func Int64(max int64) int64 {
	return baseRand(max)
}

func Float32(max float32) float32 {
	return float32(baseRand(int64(max*100.0))) / 100.0
}

func Float64(max float64) float64 {
	return float64(baseRand(int64(max*100.0))) / 100.0
}

// 生成指定长度的随机字符串
// 可以指定多个seed组成，seed可以自己构造
// 如果不指定seed，默认使用 DefaultSeed ，即小写字母+大写字母+数字
func RandomString(length int, seeds ...Seed) string {
	if length <= 0 {
		return ""
	}

	seed := Seed("")
	for _, s := range seeds {
		seed += s
	}
	if seed == "" {
		seed = DefaultSeed
	}

	bytes := make([]byte, 0)
	for i := 0; i < length; i++ {
		bytes = append(bytes, seed[int(baseRand(int64(len(seed))))])
	}
	return string(bytes)
}
