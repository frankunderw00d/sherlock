package log

import "io"

type (
	Level int // 级别
	Flag  int // Header 构建标识
)

const (
	LevelAll   Level = iota // 开启全级别日志
	LevelTrace              // 追踪开发细粒度日志
	LevelDebug              // 追踪开发日志
	LevelInfo               // 信息日志，强调应用程序的运行过程中，输出程序运行的一些重要信息，但是不能滥用，避免打印过多的日志
	LevelWarn               // 警告日志，表明会出现潜在错误的情形，有些信息不是错误信息，但是也要给的一些提示
	LevelError              // 错误日志，虽然发生错误事件，但仍然不影响系统的继续运行
	LevelFatal              // 致命错误日志，严重的错误事件将会导致应用程序的退出。重大错误，这种级别可以直接停止程序
	LevelOff                // 关闭全级别日志
)

const (
	FlagNone  Flag = 1 << iota // 不构建 Header
	FlagUTC                    // UTC时间
	FlagDate                   // 2006-01-02
	FlagTime                   // 2006-01-02 15:04:05
	FlagMTime                  // 2006-01-02 15:04:05.000
	FlagLPath                  // 打印语句长路径，此 Flag 触发 Caller 调用堆栈信息
	FlagSPath                  // 打印语句短路径，此 Flag 触发 Caller 调用堆栈信息
	FlagLine                   // 打印语句行号，此 Flag 触发 Caller 调用堆栈信息

	FlagDefault = FlagTime | FlagSPath | FlagLine // 默认组合
)

var (
	defaultKernel Kernel
)

func init() {
	defaultKernel = NewKernel()
	defaultKernel.SetSkipDepth(DefaultSkipCallDepth + 2) // 为了构建对外直接提供的打印函数，需要再加 1 获取外界调用堆栈信息
}

func (l Level) String() string {
	switch l {
	case LevelAll:
		return "ALL"
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelOff:
		return "OFF"
	default:
		return "UNDEFINED"
	}
}

// 设置过滤层级
func SetFilterLevel(level Level) {
	defaultKernel.SetFilterLevel(level)
}

// 设置构建标识
func SetFlag(flag Flag) {
	defaultKernel.SetFlag(flag)
}

// 设置跳过调用深度
func SetSkipDepth(depth int) {
	defaultKernel.SetSkipDepth(depth)
}

// 设置格式化者
func SetFormatter(formatter Formatter) {
	defaultKernel.SetFormatter(formatter)
}

// 设置输出钩子，最终的打印输出方
func SetHook(hook io.WriteCloser) {
	defaultKernel.SetHook(hook)
}

func Trace(v ...interface{}) {
	defaultKernel.Trace(v...)
}

func TraceF(format string, v ...interface{}) {
	defaultKernel.TraceF(format, v...)
}

func TraceLn(v ...interface{}) {
	defaultKernel.TraceLn(v...)
}

func Debug(v ...interface{}) {
	defaultKernel.Debug(v...)
}

func DebugF(format string, v ...interface{}) {
	defaultKernel.DebugF(format, v...)
}

func DebugLn(v ...interface{}) {
	defaultKernel.DebugLn(v...)
}

func Info(v ...interface{}) {
	defaultKernel.Info(v...)
}

func InfoF(format string, v ...interface{}) {
	defaultKernel.InfoF(format, v...)
}

func InfoLn(v ...interface{}) {
	defaultKernel.InfoLn(v...)
}

func Warn(v ...interface{}) {
	defaultKernel.Warn(v...)
}

func WarnF(format string, v ...interface{}) {
	defaultKernel.WarnF(format, v...)
}

func WarnLn(v ...interface{}) {
	defaultKernel.WarnLn(v...)
}

func Error(v ...interface{}) {
	defaultKernel.Error(v...)
}

func ErrorF(format string, v ...interface{}) {
	defaultKernel.ErrorF(format, v...)
}

func ErrorLn(v ...interface{}) {
	defaultKernel.ErrorLn(v...)
}

func Fatal(v ...interface{}) {
	defaultKernel.Fatal(v...)
}

func FatalF(format string, v ...interface{}) {
	defaultKernel.FatalF(format, v...)
}

func FatalLn(v ...interface{}) {
	defaultKernel.FatalLn(v...)
}
