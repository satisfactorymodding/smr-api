//go:build tools
// +build tools

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

// go generate -tags tools -x ./...

//go:generate protoc -I./proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative proto/parser/parser.proto
//go:generate go run github.com/99designs/gqlgen generate
//go:generate go run github.com/swaggo/swag/cmd/swag init --generalInfo cmd/api/serve.go
//go:generate go run tools.go

func main() {
	generateEnt()
}

func generateEnt() {
	err := entc.Generate("./db/schema", &gen.Config{
		Target:  "./ent",
		Package: "github.com/satisfactorymodding/smr-api/ent",
		Features: []gen.Feature{
			gen.FeatureLock,
			gen.FeatureModifier,
			gen.FeatureUpsert,
			gen.FeaturePrivacy,
			gen.FeatureIntercept,
			gen.FeatureSnapshot,
		},
		IDType: &field.TypeInfo{
			Type: field.TypeString,
		},
	})

	if err != nil {
		panic(err)
	}
}
