package main

import (
	"context"
	"os"

	"github.com/satisfactorymodding/smr-api/validation"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	f, _ := os.ReadFile(os.Args[1])

	validation.InitializeValidator()
	info, err := validation.ExtractModInfo(context.Background(), f, true, "N/A")
	if err != nil {
		panic(err)
	}
	_, err = validation.ExtractMetadata(context.Background(), f, info.GameVersion, info.ModReference)
	if err != nil {
		panic(err)
	}
}
