package logging

import (
	"bytes"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/spf13/viper"
)

func TestLoggingDev(t *testing.T) {
	buffer := new(bytes.Buffer)

	testza.AssertNoError(t, SetupLogger(buffer))

	PrintLoggingDummies()

	out := regexp.MustCompile(`\x1b\[\d+m`).ReplaceAll(buffer.Bytes(), []byte{})

	lines := bytes.Split(out, []byte("\n"))

	path := filepath.Join("logging", "logging_dummies.go")

	testza.AssertContains(t, string(lines[0]), "INF "+path+":16 A: 1 service=api")
	testza.AssertContains(t, string(lines[1]), "INF "+path+":17 B: 2 service=api")
	testza.AssertContains(t, string(lines[2]), "INF "+path+":20 C: 3 service=api")
	testza.AssertContains(t, string(lines[3]), "INF "+path+":27 D: 4 service=api")

	if t.Failed() {
		println(buffer.String())
	}
}

func TestLoggingProd(t *testing.T) {
	viper.Set("production", true)
	buffer := new(bytes.Buffer)

	testza.AssertNoError(t, SetupLogger(buffer))

	PrintLoggingDummies()

	lines := bytes.Split(buffer.Bytes(), []byte("\n"))

	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.PrintLoggingDummies","file":".*?/logging/logging_dummies.go","line":16},"msg":"A: 1"}`).Match(lines[0]))
	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.PrintLoggingDummies","file":".*?/logging/logging_dummies.go","line":17},"msg":"B: 2"}`).Match(lines[1]))
	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.PrintLoggingDummies.func1","file":".*?/logging/logging_dummies.go","line":20},"msg":"C: 3"}`).Match(lines[2]))
	testza.AssertTrue(t, regexp.MustCompile(`\{"time":".+?","level":"INFO","source":\{"function":"github.com/satisfactorymodding/smr-api/logging.PrintLoggingDummies.func2","file":".*?/logging/logging_dummies.go","line":27},"msg":"D: 4"}`).Match(lines[3]))

	if t.Failed() {
		println(buffer.String())
	}
}
