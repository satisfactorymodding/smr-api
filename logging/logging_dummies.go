package logging

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Vilsol/slox"
)

// PrintLoggingDummies logs messages for testing purposes.
// It's in a separate file from the actual logging tests because the tests care
// about what file line numbers the calls are from. Having them in a separate file
// makes editing the tests easier without breaking the line numbers.
func PrintLoggingDummies() {
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
}
