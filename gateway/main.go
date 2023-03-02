package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mwitter-backend-gRPC/proto/v1/mweets"
	"google.golang.org/grpc"
)

const (
	portNumber           = "9000"
	gRPCServerPortNumber = "9001"
)

func main() {
	// gRPC gateway와 gRPC server를 이어줌
	ctx := context.Background()
	// HTTP 요청이 오면 그 요청을 gRPC server에 보내기 전 처리 용(미들웨어)
	mux := runtime.NewServeMux()

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	log.Printf("hostIs :: %s", host)

	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err := pb.RegisterMweetsHandlerFromEndpoint(
		ctx,
		mux,
		host+":"+gRPCServerPortNumber,
		options,
	)
	if err != nil {
		log.Fatal("failed to register gRPC gateway: ", err.Error())
	}

	log.Printf("start HTTP server on %s port", portNumber)

	err = http.ListenAndServe(":"+portNumber, mux)

	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}

}
