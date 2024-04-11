package log

import (
	"context"
	"fmt"
	"github.com/kissjingalex/virtpay/internal/util/net/ip"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// const
const (
	DebugLevel = zerolog.DebugLevel
	InfoLevel  = zerolog.InfoLevel
	WarnLevel  = zerolog.WarnLevel
	ErrorLevel = zerolog.ErrorLevel
	FatalLevel = zerolog.FatalLevel
	Disabled   = zerolog.Disabled
)

var (
	logger = func() *zerolog.Logger {
		l := zerolog.New(ConsoleWriter()).With().Timestamp().Str("host", ip.Hostname()).Logger()
		return &l
	}()
	callDepth   = 2
	async       = false
	defaultName = "defo-server"
)

// Level type.
type Level = zerolog.Level // go1.9 type alias

// Log struct.
type Log struct {
	depth   int
	level   Level
	stack   bool
	traceId string
	zlogger *zerolog.Logger
	sampler zerolog.Sampler
}

type ctxKey struct{}

// Ctx is a wrapper for *Log
type Ctx struct {
	logger *Log
}

func init() {
	zerolog.MessageFieldName = "desc"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = ""
}

// NewLog .
func NewLog() *Log {
	return defaultLogger(nil, callDepth, DebugLevel, defaultName)
}

// NewLogger .
func NewLogger(project, env, redisURL string, level zerolog.Level) *Log {
	l := defaultLogger(nil, callDepth, level, project)
	SetEnableAsync()
	SetAttachment(map[string]string{
		"project": project,
		"env":     env,
	})

	w := append([]io.Writer{}, ConsoleWriter())
	if redisURL != "" {
		w = append(w, RedisWriter(RedisConfig{
			Level:     level,
			RedisURL:  redisURL,
			RedisPass: "",
			LogKey:    "wtserver:go:basic:log",
		}))
	}
	SetOutput(w...)
	return l
}

func defaultLogger(b *Log, depth int, level Level, name ...string) *Log {
	if b == nil {
		b = &Log{
			zlogger: logger,
		}
	} else {
		// snapshot
		b = func() *Log {
			nb := *b
			return &nb
		}()
	}
	b.depth = depth
	b.level = level
	if len(name) != 0 {
		fields := map[string]interface{}{
			"name": name[0],
			"ts":   fmt.Sprintf("%d.%d", time.Now().Unix(), time.Now().Nanosecond()/1000),
		}
		b.zlogger = logPointer(b.zlogger.With().Fields(fields).Logger())
	}
	return b
}

func logPointer(z zerolog.Logger) *zerolog.Logger {
	return &z
}

// Debug level msg
func Debug(name ...string) *Log {
	return defaultLogger(nil, callDepth, DebugLevel, name...)
}

// Info level msg
func Info(name ...string) *Log {
	return defaultLogger(nil, callDepth, InfoLevel, name...)
}

// Warn level msg
func Warn(name ...string) *Log {
	return defaultLogger(nil, callDepth, WarnLevel, name...)
}

// Error level msg
func Error(name ...string) *Log {
	return defaultLogger(nil, callDepth, ErrorLevel, name...)
}

// SetOutput set multi log writer, careful, all SetXXX method are non-thread safe.
func SetOutput(w ...io.Writer) {
	switch len(w) {
	case 0:
		return
	case 1:
		*logger = logger.Output(w[0])
	default:
		*logger = logger.Output(zerolog.MultiLevelWriter(w...))
	}
}

// SetLevel set global log max level.
func SetLevel(l Level) {
	zerolog.SetGlobalLevel(l)
}

// SetCallDepth set call depth for show line number.
func SetCallDepth(n int) {
	callDepth = n
}

// SetEnableAsync enables async log, should use Close func to wait all flushed.
func SetEnableAsync() {
	async = true
}

// SetAttachment add global kv to logger
func SetAttachment(kv map[string]string) {
	for k, v := range kv {
		*logger = logger.With().Str(k, v).Logger()
	}
}

// WithContext create a wrapped *Log
func WithContext(ctx ...context.Context) *Ctx {
	if len(ctx) == 0 || ctx[0] == nil {
		t := defaultLogger(nil, callDepth, Disabled)
		return &Ctx{logger: t}
	}

	if l, ok := ctx[0].Value(ctxKey{}).(*Ctx); ok {
		return l
	}

	t := defaultLogger(nil, callDepth, Disabled)
	return &Ctx{logger: t}
}

// Sampler is zerolog Sampler alias
type Sampler = zerolog.Sampler

// BasicSampler is a sampler that will send every Nth events, regardless of
// there level.
type BasicSampler = zerolog.BasicSampler

// BurstSampler lets Burst events pass per Period then pass the decision to
// NextSampler. If Sampler is not set, all subsequent events are rejected.
type BurstSampler = zerolog.BurstSampler

// LevelSampler applies a different sampler for each level.
type LevelSampler = zerolog.LevelSampler

// Sample logger
func Sample(sampler zerolog.Sampler) *Ctx {
	t := defaultLogger(nil, callDepth, Disabled)
	t.sampler = sampler
	return &Ctx{logger: t}
}

// Debug level msg
func (b *Ctx) Debug(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, DebugLevel, name...)
}

// Info level msg
func (b *Ctx) Info(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, InfoLevel, name...)
}

// Warn level msg
func (b *Ctx) Warn(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, WarnLevel, name...)
}

// Error level msg
func (b *Ctx) Error(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, ErrorLevel, name...)
}

// KV is log kv pairs.
func (b *Ctx) KV(key string, val string) *Log {
	b.logger.KV(key, val)
	return b.logger
}

// Context add wappered *Log to context
func (b *Log) Context() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, &Ctx{logger: b})
}

// KV is log kv pairs.
func (b *Log) KV(key string, val string) *Log {
	b.zlogger = logPointer(b.zlogger.With().Str(key, val).Logger())
	return b
}

func (b *Log) TraceID(traceId string) *Log {
	b.traceId = traceId
	return b
}

// Stack enables stack trace
func (b *Log) Stack() *Log {
	b.stack = true
	return b
}

// Msg output
func (b *Log) Msg(msg ...interface{}) {
	b.depth++
	switch len(msg) {
	case 0:
		return
	case 1:
		switch v := msg[0].(type) {
		case string:
			b.Msgf(v)
		default:
			b.Msgf("%v", v)
		}
	default:
		fmtStr := strings.Repeat("%v, ", len(msg))
		b.Msgf(fmtStr[:len(fmtStr)-2], msg...) // shrink last ', '
	}
}

// Msgf formatted output
func (b *Log) Msgf(msg string, v ...interface{}) {

	if b.depth != 0 {
		msg = callInfo(b.depth) + msg
	}

	var l = *b.zlogger
	if b.stack {
		v = append(v, TakeStacktrace(b.depth+1))
		// l = l.With().Str("stack", TakeStacktrace(b.depth+1)).Logger()
	}

	if b.sampler != nil {
		l = l.Sample(b.sampler)
	}
	if b.traceId != "" {
		l.WithLevel(b.level).Msgf(fmt.Sprintf("traceId:[%s] ", b.traceId)+msg, v...)
	} else {
		l.WithLevel(b.level).Msgf(fmt.Sprintf(msg, v...))
	}

}

func callInfo(n int) string {
	funcName, file, line, ok := runtime.Caller(n)
	if ok {
		if i := strings.Index(file, "src/"); i != 0 {
			return "[" + file[i+4:] + ":" + strconv.Itoa(line) + " " + runtime.FuncForPC(funcName).Name() + "] "
		}
		return "[" + file + ":" + strconv.Itoa(line) + " " + runtime.FuncForPC(funcName).Name() + "] "
	}
	return ""
}

// 实现zk logger

// Printf logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Printf(format string, args ...interface{}) {
	b.Msgf(format, args)
}

// 实现grpclog v2

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (b *Log) Info(args ...interface{}) {
	b.Msg(args)
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (b *Log) Infoln(args ...interface{}) {
	b.Msg(args)
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Infof(format string, args ...interface{}) {
	b.Msgf(format, args)
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (b *Log) Warning(args ...interface{}) {
	b.level = WarnLevel
	b.Msg(args)
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (b *Log) Warningln(args ...interface{}) {
	b.level = WarnLevel
	b.Msg(args)
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Warningf(format string, args ...interface{}) {
	b.level = WarnLevel
	b.Msgf(format, args)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (b *Log) Error(args ...interface{}) {
	b.level = ErrorLevel
	b.Msg(args)
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (b *Log) Errorln(args ...interface{}) {
	b.level = ErrorLevel
	b.Msg(args)
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Errorf(format string, args ...interface{}) {
	b.level = ErrorLevel
	b.Msgf(format, args)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (b *Log) Fatal(args ...interface{}) {
	b.level = FatalLevel
	b.Msg(args)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (b *Log) Fatalln(args ...interface{}) {
	b.level = FatalLevel
	b.Msg(args)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (b *Log) Fatalf(format string, args ...interface{}) {
	b.level = FatalLevel
	b.Msgf(format, args)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (b *Log) V(l int) bool {
	return true
}

var stacktracePool = sync.Pool{
	New: func() interface{} {
		return newProgramCounters(64)
	},
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}

var bufferPool = NewBytesPool()

// TakeStacktrace is helper func to take snap short of stack trace.
func TakeStacktrace(optionalSkip ...int) string {
	skip := 2
	if len(optionalSkip) != 0 {
		skip = optionalSkip[0]
	}

	buff := bufferPool.Get()
	defer buff.Free()

	programCounters := stacktracePool.Get().(*programCounters)
	defer stacktracePool.Put(programCounters)

	var numFrames int
	for {
		// Skip the call to runtime.Counters and takeStacktrace so that the
		// program counters start at the caller of takeStacktrace.
		numFrames = runtime.Callers(skip, programCounters.pcs)
		if numFrames < len(programCounters.pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		programCounters = newProgramCounters(len(programCounters.pcs) * 2)
	}

	frames := runtime.CallersFrames(programCounters.pcs[:numFrames])

	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	i := 0
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i != 0 {
			buff.AppendByte('\n')
		}
		i++
		buff.AppendString(frame.Function)
		buff.AppendByte('\n')
		buff.AppendByte('\t')
		buff.AppendString(frame.File)
		buff.AppendByte(':')
		buff.AppendInt(int64(frame.Line))
	}

	return buff.String()
}

var closeFuncList []func() error

func registerCloseFunc(f func() error) {
	closeFuncList = append(closeFuncList, f)
}

// Close flushes all log writer
func Close() {
	for i := range closeFuncList {
		closeFuncList[i]()
	}
}
