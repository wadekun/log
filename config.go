// author: wsfuyibing <websearch@163.com>
// date: 2021-02-22

package log

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/fuyibing/log/adapters"
	"github.com/fuyibing/log/interfaces"
)

// 日志配置.
type configuration struct {
	AdapterName string             `yaml:"adapter"`
	LevelName   string             `yaml:"level"`
	Adapter     interfaces.Adapter `yaml:"-"`
	Level       interfaces.Level   `yaml:"-"`
	Handler     interfaces.Handler `yaml:"-"`
	TimeFormat  string             `yaml:"time"`
	SpanId      string             `yaml:"span-id"`
	SpanVersion string             `yaml:"span-version"`
	TraceId     string             `yaml:"trace-id"`
	debugOn     bool
	infoOn      bool
	warnOn      bool
	errorOn     bool
}

// 创建配置实例.
func newConfiguration() interfaces.ConfigInterface {
	// 1. basic
	o := &configuration{
		Adapter:     interfaces.DefaultAdapter,
		Level:       interfaces.DefaultLevel,
		TimeFormat:  interfaces.DefaultTimeFormat,
		SpanId:      interfaces.DefaultSpanId,
		SpanVersion: interfaces.DefaultSpanVersion,
		TraceId:     interfaces.DefaultTraceId,
	}
	// 2. extensions.
	for _, file := range []string{"./tmp/log.yaml", "../tmp/log.yaml", "./config/log.yaml", "../config/log.yaml"} {
		if o.LoadYaml(file) == nil {
			break
		}
	}
	// 3. ended.
	return o
}

func (o *configuration) DebugOn() bool { return o.debugOn }
func (o *configuration) InfoOn() bool  { return o.infoOn }
func (o *configuration) WarnOn() bool  { return o.warnOn }
func (o *configuration) ErrorOn() bool { return o.errorOn }

func (o *configuration) GetHandler() interfaces.Handler {
	return o.Handler
}

// 读取时间格式.
func (o *configuration) GetTimeFormat() string { return o.TimeFormat }

// 读取请求链名称.
func (o *configuration) GetTrace() (string, string, string) { return o.TraceId, o.SpanId, o.SpanVersion }

// 从YAML加载配置.
func (o *configuration) LoadYaml(file string) (err error) {
	var body []byte
	// 1. open file.
	if body, err = ioutil.ReadFile(file); err != nil {
		return
	}
	// 2. parse yaml.
	if err = yaml.Unmarshal(body, o); err != nil {
		return
	}
	// 2.1 adapter name.
	if o.AdapterName != "" {
		o.SetAdapter(o.AdapterName)
	}
	if o.Adapter > interfaces.AdapterOff {
		o.useAdapter()
	}
	// 2.2 level name.
	if o.LevelName != "" {
		o.SetLevel(o.LevelName)
	}
	if o.Level > interfaces.LevelOff {
		o.useLevel()
	}
	// 3. end.
	return
}

// 设置适配器名称.
func (o *configuration) SetAdapter(str string) interfaces.ConfigInterface {
	str = strings.ToUpper(str)
	for adapter, name := range interfaces.AdapterTexts {
		if name == str {
			o.Adapter = adapter
			o.AdapterName = name
			o.useAdapter()
			break
		}
	}
	return o
}

// 设置日志回调.
func (o *configuration) SetHandler(handler interfaces.Handler) interfaces.ConfigInterface {
	o.Handler = handler
	return o
}

// 设置级别名称.
func (o *configuration) SetLevel(str string) interfaces.ConfigInterface {
	str = strings.ToUpper(str)
	for level, name := range interfaces.LevelTexts {
		if name == str {
			o.Level = level
			o.LevelName = name
			o.useLevel()
			break
		}
	}
	return o
}

// 设置时间格式.
func (o *configuration) SetTimeFormat(format string) interfaces.ConfigInterface {
	o.TimeFormat = format
	return o
}

// Use adapter.
func (o *configuration) useAdapter() {
	switch o.Adapter {
	case interfaces.AdapterTerm:
		o.SetHandler(adapters.NewTerm().Run)
	case interfaces.AdapterFile:
		o.SetHandler(adapters.NewFile().Run)
	case interfaces.AdapterRedis:
		o.SetHandler(adapters.NewRedis().Run)
	}
}

// Use level.
func (o *configuration) useLevel() {
	o.debugOn = o.Level >= interfaces.LevelDebug
	o.infoOn = o.Level >= interfaces.LevelInfo
	o.warnOn = o.Level >= interfaces.LevelWarn
	o.errorOn = o.Level >= interfaces.LevelError
}
