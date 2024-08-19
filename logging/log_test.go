package logging

import (
	"bytes"
	"context"
	"log/slog"
	"regexp"
	"sync"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/slox"
	"github.com/spf13/viper"
)

func TestLoggingDev(t *testing.T) {
	buffer := new(bytes.Buffer)

	testza.AssertNoError(t, SetupLogger(buffer))

	slog.Info("A: 1")
	slox.Info(context.Background(), "B: 2")

	func() {
		slox.Info(context.Background(), "C: 3")
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		slox.Info(context.Background(), "D: 4")
	}()
	wg.Wait()

	out := regexp.MustCompile(`\x1b\[\d+m`).ReplaceAll(buffer.Bytes(), []byte{})

	lines := bytes.Split(out, []byte("\n"))

	testza.AssertContains(t, string(lines[0]), "INF logging/log_test.go:21 A: 1 service=api")
	testza.AssertContains(t, string(lines[1]), "INF logging/log_test.go:22 B: 2 service=api")
	testza.AssertContains(t, string(lines[2]), "INF logging/log_test.go:25 C: 3 service=api")
	testza.AssertContains(t, string(lines[3]), "INF logging/log_test.go:32 D: 4 service=api")

	if t.Failed() {
		println(buffer.String())
	}
}

func TestLoggingProd(t *testing.T) {
	viper.Set("production", true)
	buffer := new(bytes.Buffer)

	testza.AssertNoError(t, SetupLogger(buffer))

	slog.Info("A: 1")
	slox.Info(context.Background(), "B: 2")

	func() {
		slox.Info(context.Background(), "C: 3")
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		slox.Info(context.Background(), "D: 4")
	}()
	wg.Wait()

	lines := bytes.Split(buffer.Bytes(), []byte("\n"))

	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.TestLoggingProd","file":".*?/logging/log_test.go","line":56},"msg":"A: 1"}`).Match(lines[0]))
	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.TestLoggingProd","file":".*?/logging/log_test.go","line":57},"msg":"B: 2"}`).Match(lines[1]))
	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.TestLoggingProd.func1","file":".*?/logging/log_test.go","line":60},"msg":"C: 3"}`).Match(lines[2]))
	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.TestLoggingProd.func2","file":".*?/logging/log_test.go","line":67},"msg":"D: 4"}`).Match(lines[3]))

	if t.Failed() {
		println(buffer.String())
	}
}
