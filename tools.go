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
//go:generate go run github.com/jmattheis/goverter/cmd/goverter@v1.0.0 gen -g wrapErrors -g ignoreUnexported ./conversion
//go:generate go run github.com/99designs/gqlgen@v0.17.39 generate
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.2 init --generalInfo cmd/api/serve.go --output ./generated/docs
//go:generate go run tools.go

func main() {
	generateEnt()
}

func generateEnt() {
	err := entc.Generate("./db/schema", &gen.Config{
		Target:  "./generated/ent",
		Package: "github.com/satisfactorymodding/smr-api/generated/ent",
		Features: []gen.Feature{
			gen.FeatureModifier,
			gen.FeatureIntercept,
			gen.FeatureSnapshot,
			gen.FeatureExecQuery,
			gen.FeatureUpsert,
		},
		IDType: &field.TypeInfo{
			Type: field.TypeString,
		},
	})

	if err != nil {
		panic(err)
	}
}
