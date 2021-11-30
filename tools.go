//go:build tools
// +build tools

package smr

import _ "github.com/99designs/gqlgen"
import _ "github.com/swaggo/swag/cmd/swag"

//go:generate go run github.com/99designs/gqlgen generate
//go:generate go run github.com/swaggo/swag/cmd/swag init --generalInfo cmd/api/serve.go
