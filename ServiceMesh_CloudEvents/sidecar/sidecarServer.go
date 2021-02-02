package sidecar

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	_ "os"
)

func NewServer() {
	// create a TCP Listener on Port 9000
	lis, err := net.Listen("tcp", ":9000")
	// this how you handle errors in Golang
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 %v", err)
	}
	fmt.Println("TEST")

	// this is just a structure that has an interface with needed function SayHello
	// have a look at /chat/chat2.go for more information
	s := Server{}
	fmt.Println("Server Created")
	//create the GRCP Server
	grpcServer := grpc.NewServer()
	fmt.Println("GRPC created")
	RegisterChatServiceServer(grpcServer, &s)
	fmt.Println("GRPC register")
	// start listening on port 9000 for rpc
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000 %v", err)
	}

}
