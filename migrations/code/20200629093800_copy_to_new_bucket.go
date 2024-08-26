package code

import (
	"github.com/lab259/go-migration"
)

func init() {
	migration.NewCodeMigration(
		func(_ interface{}) error {
			// Antiquated
			return nil
		},
	)
}
