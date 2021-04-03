package log

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type (
	// 格式化定义
	Formatter interface {
		Format(Level, Flag, time.Time, []runtime.Frame, string) string
	}

	// 字符串格式化定义实现
	StringFormatter struct {
		prefix []string
	}

	// JSON 格式化定义实现
	JSONFormatter struct{}
)

const ()

var ()

// 实例化字符串格式化定义实现
func NewStringFormatter(prefix []string) Formatter {
	return &StringFormatter{prefix: prefix}
}

// 实例化JSON 格式化定义实现
func NewJSONFormatter() Formatter {
	return &JSONFormatter{}
}

// 构建 Header 标识，默认采取跳过的 k.depth 层数后的第 0 层堆栈数据
func (s *StringFormatter) Format(level Level, flag Flag, t time.Time, frames []runtime.Frame, str string) string {
	// 为输出字符串添加换行
	if str[len(str)-1] != '\n' {
		str += "\n"
	}

	// 构建 Header
	header := ""
	if flag&FlagNone == 0 { // 构建 header
		header += fmt.Sprintf("[%s]", level.String())

		if flag&FlagUTC != 0 { // UTC 时间
			t = t.UTC()
		}

		if flag&(FlagDate|FlagTime|FlagMTime) != 0 {
			header += "["

			year, month, day := t.Date()
			header += fmt.Sprintf("%4d-%02d-%02d", year, month, day) // date

			if flag&(FlagTime|FlagMTime) != 0 {
				hour, minute, second := t.Clock()
				header += fmt.Sprintf(" %02d:%02d:%02d", hour, minute, second) // time

				if flag&FlagMTime != 0 {
					header += fmt.Sprintf(".%03d", t.Nanosecond()/1e6) // milliseconds
				}
			}

			header += "]"
		}

		if flag&(FlagLPath|FlagSPath|FlagLine) != 0 {
			header += "["

			if flag&(FlagLPath|FlagSPath) != 0 { // path
				path := frames[0].Function

				if flag&FlagSPath != 0 { // short path
					cPath := path
					for i := len(cPath) - 1; i > 0; i-- {
						if cPath[i] == '/' {
							path = cPath[i+1:]
							break
						}
					}
				}
				header += path
			}

			if flag&FlagLine != 0 { // line
				if flag&(FlagLPath|FlagSPath) != 0 { // in case of require path and line, add ':'
					header += ":"
				}
				header += fmt.Sprintf("%d", frames[0].Line)
			}

			header += "]"
		}

		if s.prefix != nil && len(s.prefix) > 0 {
			for _, value := range s.prefix {
				header += fmt.Sprintf("[%s]", value)
			}
		}

		return header + " " + str
	} else {
		return str
	}
}

// 构建 Header 标识，默认采取跳过的 k.depth 层数后的第 0 层堆栈数据
func (j *JSONFormatter) Format(level Level, flag Flag, t time.Time, frames []runtime.Frame, str string) string {
	type A struct {
		Level Level  `json:"level"`
		Time  string `json:"time"`
		File  string `json:"file"`
		Line  int    `json:"line"`
		Log   string `json:"log"`
	}

	a := &A{
		Level: level,
		Log:   str,
	}
	if flag&FlagNone == 0 { // 构建 A 结构体
		if flag&FlagUTC != 0 { // UTC 时间
			t = t.UTC()
		}

		if flag&(FlagDate|FlagTime|FlagMTime) != 0 {
			timeString := ""
			year, month, day := t.Date()
			timeString += fmt.Sprintf("%4d-%02d-%02d", year, month, day) // date

			if flag&(FlagTime|FlagMTime) != 0 {
				hour, minute, second := t.Clock()
				timeString += fmt.Sprintf(" %02d:%02d:%02d", hour, minute, second) // time

				if flag&FlagMTime != 0 {
					timeString += fmt.Sprintf(".%03d", t.Nanosecond()/1e6) // milliseconds
				}
			}

			a.Time = timeString
		}

		if flag&(FlagLPath|FlagSPath) != 0 { // path
			path := frames[0].Function

			if flag&FlagSPath != 0 { // short path
				cPath := path
				for i := len(cPath) - 1; i > 0; i-- {
					if cPath[i] == '/' {
						path = cPath[i+1:]
						break
					}
				}
			}

			a.File = path
		}

		if flag&FlagLine != 0 { // line
			a.Line = frames[0].Line
		}
	}

	data, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(data) + "\n"
}
