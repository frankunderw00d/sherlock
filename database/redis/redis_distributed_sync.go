package redis

import (
	uRand "sherlock/util/rand"
)

type (
	// redis 分布式同步机制
	DistributedSync interface {
		// 加入
		Join() error

		// 清除
		Clean() error

		// 从 set 中拉取所有成员，构建成员 list 键名，右 push
		Publish(string) error

		// 从自己的 list 中 左 拉取元素
		Subscribe() (string, error)
	}

	// redis 分布式同步机制
	distributedSync struct {
		set   string
		local string
	}
)

const (
	// redis 分布式同步机制键
	DistributedSyncKey = "DSKey:"
	// redis 分布式同步列表键
	DistributedSyncOwnKey = "DSOKey:"
)

var ()

// 新建分布式同步机制
func NewDistributedSync(set string) DistributedSync {
	return &distributedSync{
		set:   set,
		local: DistributedSyncOwnKey + uRand.RandomString(8),
	}
}

// 加入
func (ds *distributedSync) Join() error {
	if _, err := SAdd(DistributedSyncKey+ds.set, ds.local); err != nil {
		return err
	}
	return nil
}

// 清除
func (ds *distributedSync) Clean() error {
	// 1.从 set 中删除自身
	if _, err := SRem(DistributedSyncKey+ds.set, ds.local); err != nil {
		return err
	}
	// 2.删除自身的 list
	if _, err := Del(ds.local); err != nil {
		return err
	}
	return nil
}

// 从 set 中拉取所有成员，构建成员 list 键名，右 push
func (ds *distributedSync) Publish(v string) error {
	// 1.拉取指定 set 成员
	members, err := SMembers(DistributedSyncKey + ds.set)
	if err != nil {
		return err
	}

	// 无成员不操作
	if len(members) <= 0 {
		return nil
	}

	// 2.向所有成员的 list 右 push 发布消息
	for _, member := range members {
		if _, err := RPush(member, v); err != nil {
			return err
		}
	}

	return nil
}

// 从自己的 list 中 左 拉取元素
func (ds *distributedSync) Subscribe() (string, error) {
	return LPop(ds.local)
}
