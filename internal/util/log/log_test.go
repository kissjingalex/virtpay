package log

import (
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

var sLog = Sample(&BasicSampler{N: 2})                          //sampler log
var cLog = Error("abc").KV("a1", "b1").KV("a2", "b2").Context() // log ctx

func load() {
	SetAttachment(map[string]string{
		"common-tag": "all log must have this",
	})
}

func TestLog(t *testing.T) {
	load()
	//SetOutput(ConsoleWriter(), RedisWriter())
	//SetOutput(ConsoleWriter(ConsoleConf{
	//	Level:ErrorLevel,
	//}))

	Debug().Msg("test")
	Debug("hello").Msg("test with name=hello")
	Debug("kvLog").KV("x", "y").Msg("test with kv(x=y) and name(kvLog)")
	Debug().Msg("test")

	fmt.Println()

	WithContext(cLog).Warn().Msg("ctx log")
	WithContext(cLog).Warn("@override").Msg("with name override")
	WithContext(cLog).Warn().KV("a1", "x1").Msg("kv override")

	fmt.Println()

	sLog.Info().Msg("this should not output")
	sLog.Info().Msg("with no name")
	sLog.Info().KV("hello", "world").Msg("this should not output")
	sLog.Info("").KV("model", "test").Msg("with empty name")
	sLog.Info().Stack().Msg("this should not output")
	sLog.Info("test-info").Msg("with name test-info")
	sLog.Info().Msg("this should not output")
	sLog.Info().Msg("with no name")

	fmt.Println()

	Error().Msg("x", "y")

	SetCallDepth(2)
	SetAttachment(map[string]string{
		"type": "kg-api-dev",
	})

	begin := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(4)

	go func(wg *sync.WaitGroup) {
		Debug("debug").Msg("debug msg")
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		Info("info").Msg("info msg")
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		Warn("warn").Msg("warn msg")
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		Error("error").Msg("error msg")
		wg.Done()
	}(wg)

	wg.Wait()
	fmt.Println("using", time.Since(begin))

	Debug().KV("key1", "val1").KV("key2", "val2").Msgf("debug test kv")
	Debug().Msgf("hello:%s, time:%s", "123", time.Now())

	hello()
	funcA()

	fmt.Println()

	m1 := Error("abc").KV("a1", "b1").KV("a2", "b2").Context()
	m2 := WithContext(m1).KV("m1", "n1").Context()
	m3 := WithContext(m2).Error("uvw").KV("m2", "n2").Context()

	WithContext(m3).Info("xyz").Msg("hello")

	fmt.Println()

	smp := Sample(&BasicSampler{N: 2})
	smp.Info("name").Msg("1")
	smp.Info("name").Msg("2")
	smp.Info("name").Msg("3")

	fmt.Println()

	bSample := Sample(&BurstSampler{
		Burst:       3,
		Period:      time.Second,
		NextSampler: &BasicSampler{N: 2},
	})

	bSample.Info("info").Msg("1")
	bSample.Info("info").Msg("2")
	bSample.Info("info").Msg("3")
	bSample.Info("info").Msg("4")
	bSample.Info("info").Msg("5")

	<-time.After(time.Second)
	bSample.Debug("debug").Msg("1")
	bSample.Debug("debug").Msg("2")
	bSample.Debug("debug").Msg("3")
	bSample.Debug("debug").Msg("4")
	bSample.Debug("debug").Msg("5")
}

func hello() {
	Debug().Msg("hello world")
}

func funcA() {
	Debug().Msg("funcA")
	funcB()
}

func funcB() {
	Debug("stack").Stack().Msg("funcB")
}

func BenchmarkDebugStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		debug.Stack()
	}
}

func BenchmarkTakeStacktrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TakeStacktrace(2)
	}
}
