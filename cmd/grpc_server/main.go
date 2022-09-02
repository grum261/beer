package main

import (
	"context"
	"log"

	grpcserver "github.com/grum261/beer/internal/app/grpc_server"
)

///home/grum231/prog/src/github.com/grum261/beer/internal/app/grpc_server

func main() {
	log.Fatal(grpcserver.Run(context.Background()))
}
