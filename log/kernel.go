package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type (
	// 日志核心定义
	Kernel interface {
		// 基础打印
		Print(Level, string)

		// 设置过滤层级
		SetFilterLevel(Level)

		// 设置构建标识
		SetFlag(Flag)

		// 设置跳过调用深度
		SetSkipDepth(int)

		// 设置格式化者
		SetFormatter(Formatter)

		// 设置输出钩子，最终的打印输出方
		SetHook(io.WriteCloser)

		// 使用以下提供的打印函数，请先调用 Kernel.SetSkipDepth(DefaultSkipCallDepth+1)
		// Trace 打印
		Trace(v ...interface{})
		TraceF(format string, v ...interface{})
		TraceLn(v ...interface{})
		// Debug 打印
		Debug(v ...interface{})
		DebugF(format string, v ...interface{})
		DebugLn(v ...interface{})
		// Info 打印
		Info(v ...interface{})
		InfoF(format string, v ...interface{})
		InfoLn(v ...interface{})
		// Warn 打印
		Warn(v ...interface{})
		WarnF(format string, v ...interface{})
		WarnLn(v ...interface{})
		// Error 打印
		Error(v ...interface{})
		ErrorF(format string, v ...interface{})
		ErrorLn(v ...interface{})
		// Fatal 打印，打印后程序会以 1 为状态码退出
		Fatal(v ...interface{})
		FatalF(format string, v ...interface{})
		FatalLn(v ...interface{})
	}

	// 日志核心实现
	kernel struct {
		filter    Level          // 过滤级别，打印级别必须大于等于 filter 才会给予打印，默认 LevelDebug 级别
		flag      Flag           // Header 构建标识，默认 FlagDefault
		depth     int            // 跳过调用深度，默认 3 层
		formatter Formatter      // 格式化者，将打印的信息格式化，默认为 StringFormatter
		hook      io.WriteCloser // 输出钩子，最终的打印输出方，默认为 os.Stdout
		mutex     sync.Mutex     // 并发锁
	}
)

const (
	// 默认跳过调用深度
	DefaultSkipCallDepth = 3
)

var (
	// 默认输出钩子
	StandardOutputHook = os.Stdout
)

// 实例化日志核心
func NewKernel() Kernel {
	return &kernel{
		filter:    LevelDebug,
		flag:      FlagDefault,
		depth:     DefaultSkipCallDepth,
		formatter: NewStringFormatter(nil),
		hook:      StandardOutputHook,
		mutex:     sync.Mutex{},
	}
}

// 基础打印
func (k *kernel) Print(level Level, str string) {
	if level < k.filter {
		return
	}

	// 记录打印时间
	t := time.Now()

	// 加锁，此锁为了多线程串流用
	k.mutex.Lock()
	defer k.mutex.Unlock()

	// 解锁后再加锁，防止调用堆栈信息占用时间
	var frames []runtime.Frame
	if k.flag&(FlagLPath|FlagSPath|FlagLine) != 0 {
		k.mutex.Unlock()
		frames = getCallers(k.depth)
		k.mutex.Lock()
	}

	// 格式化打印
	var final string
	if k.formatter != nil {
		final = k.formatter.Format(level, k.flag, t, frames, str)
	}

	// 输出
	if k.hook != nil {
		_, _ = k.hook.Write([]byte(final))
	}
}

// 设置过滤层级
func (k *kernel) SetFilterLevel(level Level) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.filter = level
}

// 设置构建标识
func (k *kernel) SetFlag(flag Flag) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.flag = flag
}

// 设置跳过调用深度
func (k *kernel) SetSkipDepth(depth int) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.depth = depth
}

// 设置格式化者
func (k *kernel) SetFormatter(formatter Formatter) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.formatter = formatter
}

// 设置输出钩子，最终的打印输出方
func (k *kernel) SetHook(hook io.WriteCloser) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if k.hook != nil {
		_ = k.hook.Close()
	}

	k.hook = hook
}

func (k *kernel) Trace(v ...interface{}) { k.Print(LevelTrace, fmt.Sprint(v...)) }
func (k *kernel) TraceF(format string, v ...interface{}) {
	k.Print(LevelTrace, fmt.Sprintf(format, v...))
}
func (k *kernel) TraceLn(v ...interface{}) { k.Print(LevelTrace, fmt.Sprintln(v...)) }

func (k *kernel) Debug(v ...interface{}) { k.Print(LevelDebug, fmt.Sprint(v...)) }
func (k *kernel) DebugF(format string, v ...interface{}) {
	k.Print(LevelDebug, fmt.Sprintf(format, v...))
}
func (k *kernel) DebugLn(v ...interface{}) { k.Print(LevelDebug, fmt.Sprint(v...)) }

func (k *kernel) Info(v ...interface{}) { k.Print(LevelInfo, fmt.Sprint(v...)) }
func (k *kernel) InfoF(format string, v ...interface{}) {
	k.Print(LevelInfo, fmt.Sprintf(format, v...))
}
func (k *kernel) InfoLn(v ...interface{}) { k.Print(LevelInfo, fmt.Sprint(v...)) }

func (k *kernel) Warn(v ...interface{}) { k.Print(LevelWarn, fmt.Sprint(v...)) }
func (k *kernel) WarnF(format string, v ...interface{}) {
	k.Print(LevelWarn, fmt.Sprintf(format, v...))
}
func (k *kernel) WarnLn(v ...interface{}) { k.Print(LevelWarn, fmt.Sprint(v...)) }

func (k *kernel) Error(v ...interface{}) { k.Print(LevelError, fmt.Sprint(v...)) }
func (k *kernel) ErrorF(format string, v ...interface{}) {
	k.Print(LevelError, fmt.Sprintf(format, v...))
}
func (k *kernel) ErrorLn(v ...interface{}) { k.Print(LevelError, fmt.Sprint(v...)) }

func (k *kernel) Fatal(v ...interface{}) {
	k.Print(LevelFatal, fmt.Sprint(v...))
	os.Exit(1)
}
func (k *kernel) FatalF(format string, v ...interface{}) {
	k.Print(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func (k *kernel) FatalLn(v ...interface{}) {
	k.Print(LevelFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// 获取堆栈上的调用信息
func getCallers(depth int) []runtime.Frame {
	pc := make([]uintptr, 100)
	n := runtime.Callers(depth, pc)
	frames := runtime.CallersFrames(pc[:n])
	frameList := make([]runtime.Frame, 0)
	for {
		frame, more := frames.Next()
		frameList = append(frameList, frame)
		if !more {
			break
		}
	}
	return frameList
}
