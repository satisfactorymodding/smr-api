package entlogger

import (
	"context"
	"fmt"

	"github.com/Vilsol/slox"
)

/*
Placed here so our [logging.StackRewriter] can exclude it
*/

func EntLogger(ctx context.Context) func(v ...interface{}) {
	return func(v ...interface{}) {
		slox.Info(ctx, fmt.Sprint(v...))
	}
}
