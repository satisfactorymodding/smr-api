//go:build tools
// +build tools

package smr

import _ "github.com/99designs/gqlgen"
import _ "github.com/swaggo/swag/cmd/swag"

//go:generate go run github.com/99designs/gqlgen generate
//go:generate go run github.com/swaggo/swag/cmd/swag init --generalInfo cmd/api/serve.go
//go:generate protoc -I./proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative proto/parser/parser.proto
