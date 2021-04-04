package sherlock

import (
	"errors"
	"sherlock/client"
	"sherlock/log"
	"sync"
)

type (
	// 服务定义
	Service interface {
		// 服务信息
		Info() string
		// 初始化
		Init(client.Client) error
		// 运行 （必须阻塞该执行线程）
		Run() error
		// 销毁
		Destroy() error
	}

	// 启动器
	Sherlock interface {
		// 初始化
		Init(name, address, token string) error
		// 运行嵌套服务
		Run(services ...Service) error
		// 关闭
		Close() error
	}

	sherlock struct {
		client client.Client
		wg     sync.WaitGroup
	}
)

const (
	Version = `
======================================
========== Sherlock v1.0.0 ===========
======================================
`
)

var (
	ErrEmptyServices = errors.New("service list is nil or empty")
)

var (
	defaultSherlock Sherlock
)

func init() {
	defaultSherlock = &sherlock{wg: sync.WaitGroup{}}
}

func NewSherlock() Sherlock {
	return &sherlock{wg: sync.WaitGroup{}}
}

// 初始化
func (s *sherlock) Init(name, address, token string) error {
	log.DebugLn(Version)

	c, err := client.NewClient(name, address, token)
	if err != nil {
		return err
	}

	s.client = c

	return nil
}

// 运行嵌套服务
func (s *sherlock) Run(services ...Service) error {
	if services == nil || len(services) == 0 {
		return ErrEmptyServices
	}

	for _, ser := range services {
		s.wg.Add(1)
		go func(service Service) {
			defer s.wg.Done()
			log.InfoF("Initialize %s service", service.Info())
			if err := service.Init(s.client); err != nil {
				log.ErrorF("Initialize %s service error : %s", service.Info(), err.Error())
				return
			}

			log.InfoF("Running %s service", service.Info())
			if err := service.Run(); err != nil {
				log.ErrorF("Running %s service error : %s", service.Info(), err.Error())
				return
			}

			log.InfoF("Destroy %s service", service.Info())
			if err := service.Destroy(); err != nil {
				log.ErrorF("Destroy %s service error : %s", service.Info(), err.Error())
				return
			}

			log.InfoF("%s service done", service.Info())
		}(ser)
	}

	return nil
}

// 关闭
func (s *sherlock) Close() error {
	// 等待所有 Service 销毁完成
	s.wg.Wait()

	s.client.Close()

	return nil
}

// 初始化
func Init(name, address, token string) error {
	return defaultSherlock.Init(name, address, token)
}

// 运行嵌套服务
func Run(services ...Service) error {
	return defaultSherlock.Run(services...)
}

// 关闭
func Close() error {
	return defaultSherlock.Close()
}
