package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func runInK8S() bool {
	return strings.ToLower(os.Getenv("LogEnv")) == "k8s"
}

// ConsoleConf is conf for console writer.
type ConsoleConf struct {
	Level   Level
	Format  string // empty or json
	NoColor bool
	Out     io.Writer
}

// ConsoleWriter gens console with conf provided.
func ConsoleWriter(consoleConf ...ConsoleConf) io.Writer {
	var conf ConsoleConf

	if len(consoleConf) == 0 {
		conf = ConsoleConf{
			Level:   DebugLevel,
			Out:     os.Stdout,
			NoColor: runInK8S(),
		}
	} else {
		conf = consoleConf[0]
	}

	if conf.Out == nil {
		conf.Out = os.Stdout
	}

	if async {
		w := NewAsyncWriter(conf.Level, conf, diodes.NewManyToOne(1024, diodes.AlertFunc(func(missed int) {
			fmt.Printf("Console dropped %d messages\n", missed)
		})), 1*time.Second)
		registerCloseFunc(w.Close)
		return w
	}
	return conf
}

const (
	cReset    = 0
	cBold     = 1
	cRed      = 31
	cGreen    = 32
	cYellow   = 33
	cBlue     = 34
	cMagenta  = 35
	cCyan     = 36
	cGray     = 37
	cDarkGray = 90
)

var consoleBufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 100))
	},
}

// Write write data to writer.
func (w ConsoleConf) Write(p []byte) (n int, err error) {
	if w.Format == "json" {
		w.Out.Write(p)
		n = len(p)
		return
	}

	var event map[string]interface{}
	err = json.Unmarshal(p, &event)
	if err != nil {
		return
	}
	buf := consoleBufPool.Get().(*bytes.Buffer)
	defer consoleBufPool.Put(buf)
	lvlColor := cReset
	level := "?????"
	if l, ok := event[zerolog.LevelFieldName].(string); ok {
		if !w.NoColor {
			lvlColor = levelColor(l)
		}
		level = strings.ToUpper(l)
	}
	if _, ok := event[zerolog.TimestampFieldName]; ok {
		event[zerolog.TimestampFieldName] = time.Now().Format("2006-01-02 15:04:05.999999")
	}
	fmt.Fprintf(buf, "%s |%s| %s",
		colorize(event[zerolog.TimestampFieldName], cDarkGray, !w.NoColor),
		colorize(level, lvlColor, !w.NoColor),
		colorize(event[zerolog.MessageFieldName], cReset, !w.NoColor))
	fields := make([]string, 0, len(event))
	for field := range event {
		switch field {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName:
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		fmt.Fprintf(buf, " %s=", colorize(field, cCyan, !w.NoColor))
		switch value := event[field].(type) {
		case string:
			if needsQuote(value) {
				buf.WriteString(strconv.Quote(value))
			} else {
				buf.WriteString(value)
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			fmt.Fprint(buf, value)
		default:
			b, err := json.Marshal(value)
			if err != nil {
				fmt.Fprintf(buf, "[error: %v]", err)
			} else {
				fmt.Fprint(buf, string(b))
			}
		}
	}
	buf.WriteByte('\n')
	buf.WriteTo(w.Out)

	n = len(p)
	return
}

// WriteLevel write data to writer with level info provided
func (w ConsoleConf) WriteLevel(level Level, p []byte) (n int, err error) {
	if level < w.Level {
		return len(p), nil
	}
	return w.Write(p)
}

func colorize(s interface{}, color int, enabled bool) string {
	if !enabled {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

func levelColor(level string) int {
	switch level {
	case "debug":
		return cMagenta
	case "info":
		return cGreen
	case "warn":
		return cYellow
	case "error", "fatal", "panic":
		return cRed
	default:
		return cReset
	}
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}
