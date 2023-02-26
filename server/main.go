package main

import (
	"log"
	"net"

	pb "github.com/mwitter-backend-gRPC/proto/v1/mweets"
	"github.com/mwitter-backend-gRPC/server/routes"
	"google.golang.org/grpc"
)

const portNumber = "9001"

func main() {
	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	routes, err := routes.GetServer()

	if err != nil {
		log.Fatalf("failed get routes: %s", err.Error())
	}

	grpcServer := grpc.NewServer()

	pb.RegisterMweetsServer(grpcServer, routes)

	log.Printf("start gRPC server on %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
