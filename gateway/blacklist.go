package gateway

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"sherlock/log"
	"sync"
)

type (
	// 黑名单 , 单独对网关提供黑名单封锁功能，以订阅方式支持启动、关闭、热更新、封禁时间等
	BlackList interface {
		// 在线订阅更新
		OnlineUpdate(*nats.Msg)
		// 在线订阅开关
		OnlineSwitch(*nats.Msg)
		// 过滤
		Filter(string) bool
	}

	blackList struct {
		open  bool
		mutex sync.Mutex
		list  map[string]struct{}
	}
)

const (
	BlackListUpdateSubject = "Gateway.BlackList.Update"
	BlackListSwitchSubject = "Gateway.BlackList.Switch"
)

var ()

func init() {}

// 新建黑名单
func NewBlackList() BlackList {
	return &blackList{
		open:  true,
		mutex: sync.Mutex{},
		list:  map[string]struct{}{},
	}
}

func (bl *blackList) OnlineUpdate(order *nats.Msg) {
	log.InfoF("Black list before update : %+v", bl.list)

	newList := new(map[string]struct{})

	if err := json.Unmarshal(order.Data, newList); err != nil {
		log.ErrorF("Unmarshal data to []string error : %s", err.Error())
		return
	}

	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	bl.list = *newList

	log.InfoF("Black list after update : %+v", bl.list)
}

func (bl *blackList) OnlineSwitch(order *nats.Msg) {
	log.InfoF("Black list before open : %v", bl.open)

	state := new(int)
	if err := json.Unmarshal(order.Data, state); err != nil {
		log.ErrorF("Unmarshal data to int error : %s", err.Error())
		return
	}

	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	bl.open = *state == 1

	log.InfoF("Black list after open : %v", bl.open)
}

func (bl *blackList) Filter(ip string) bool {
	if !bl.open {
		return true
	}

	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	// 名单为空，直接放行
	if len(bl.list) == 0 {
		return true
	}

	if _, exist := bl.list[ip]; exist {
		return false
	}

	return true
}
