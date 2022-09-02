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
	_, err := validation.ExtractModInfo(context.Background(), f, true, true, "N/A")

	if err != nil {
		panic(err)
	}
}
