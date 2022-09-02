package main

import (
	"context"
	"log"

	"github.com/grum261/beer/internal/app/gateway"
)

func main() {
	log.Fatal(gateway.Run(context.Background()))
}
