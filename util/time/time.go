package time

import "time"

type (
	Ticker interface {
		// 运行
		Run()

		// 停止
		Stop()
	}

	ticker struct {
		t *time.Ticker         // 定时器
		f func(time.Time) bool // 执行函数，形参为执行时间，返回值 false 断开定时工作，true 继续
	}
)

func NewTicker(duration time.Duration, function func(time.Time) bool) Ticker {
	return &ticker{
		t: time.NewTicker(duration),
		f: function,
	}
}

// 运行
func (t *ticker) Run() {
	if t.t == nil {
		return
	}

	go func() {
		// 开始就进行一次
		if !t.f(time.Now()) {
			return
		}

		for {
			done := false
			select {
			case now, ok := <-t.t.C:
				{
					if !ok || !t.f(now) {
						t.t.Stop()
						done = true
						break
					}
				}
			}
			if done {
				break
			}
		}
	}()
}

// 停止
func (t *ticker) Stop() {
	t.t.Stop()
}
