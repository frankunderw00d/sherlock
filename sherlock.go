package sherlock

import (
	"fmt"

	//"log"
	"sherlock/client"
	"sherlock/gateway"
	"sherlock/log"
)

type (
	Service interface {
		Run(...gateway.Gateway)
		Destroy()
		Client() client.Client
	}

	service struct {
		name    string        // 服务名称
		version string        // 服务版本
		c       client.Client // nats 客户端
	}
)

const (
	Version = `
======================================
========== Sherlock v1.0.0 ===========
======================================
`
)

var ()

func init() {}

// 新建服务
func NewService(name, version string, c client.Client) Service {
	return &service{
		name:    name,
		version: version,
		c:       c,
	}
}

// 启动
func (s *service) Run(gateways ...gateway.Gateway) {
	log.DebugF(Version)
	log.InfoF(s.info())

	if gateways == nil || len(gateways) == 0 {
		log.DebugF("[%s] Didn't have any gateway", s.name)
		return
	}

	for _, gw := range gateways {
		go func(g gateway.Gateway) {
			log.DebugF("[%s] gateway initialize", g.Name())
			// initialize
			if err := g.Init(s.c); err != nil {
				log.ErrorF("[%s] gateway initialize error : %s", g.Name(), err.Error())
				return
			}

			log.DebugF("[%s] gateway gonna running on the %s", g.Name(), g.Address())
			// run
			if err := g.Run(); err != nil {
				log.ErrorF("[%s] gateway run error : %s", g.Name(), err.Error())
				return
			}

			log.DebugF("[%s] gateway destroying", g.Name())
			// destroy
			if err := g.Destroy(); err != nil {
				log.ErrorF("[%s] gateway destroy error : %s", g.Name(), err.Error())
				return
			}
		}(gw)
	}
}

// 关闭
func (s *service) Destroy() {
	s.c.Close()
}

// 获取 nats-server 连接
func (s *service) Client() client.Client {
	return s.c
}

// 构造服务信息打印
func (s *service) info() string {
	max := 38
	infoLen := len(s.name) + len(s.version) + 3
	if infoLen > 38 {
		max = infoLen
		if max%2 != 0 {
			max += 1
		}
	}

	f := ""
	t := ""
	for i := 0; i < max; i++ {
		f += "="
		t += "="
	}

	info := fmt.Sprintf(" %s %s ", s.name, s.version)
	left := (max - infoLen) / 2
	lSide := ""
	rSide := ""
	for i := 0; i < left; i++ {
		lSide += "="
		rSide += "="
	}
	if infoLen%2 != 0 {
		rSide += ""
	}

	return fmt.Sprintf("\n%s\n%s%s%s\n%s", f, lSide, info, rSide, t)
}
