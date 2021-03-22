package master

import (
	"fmt"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"log"
	"net"
)

func Master() {
	// create a TCP Listener on Port 9000
	lis, err := net.Listen("tcp", ":9000")
	// this how you handle errors in Golang
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 %v", err)
	}
	// this is just a structure that has an interface with needed function SayHello
	// have a look at /master/chat2.go for more information
	s := Server{}
	fmt.Println("Sidecar: GRPC-Server started")
	//create the GRCP Server

	creds, _ := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// error handling omitted
	grpcServer.Serve(lis)

	//grpcServer := grpc.NewServer()
	RegisterChatServiceServer(grpcServer, &s)
	// start listening on port 9000 for rpc
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000 %v", err)
	}
}
