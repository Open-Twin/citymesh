package chat

import (
	"fmt"

	"google.golang.org/grpc"
	"log"
	"net"
)

func Master() {
	// create a TCP Listener on Port 9001
	lis, err := net.Listen("tcp", ":9001")
	// this how you handle errors in Golang
	if err != nil {
		log.Fatalf("Failed to listen on port 9001 %v", err)
	}
	fmt.Println("Starting Master Service")

	// this is just a structure that has an interface with needed function SayHello
	// have a look at /chat/chat2.go for more information
	s := Server{}
	//create the GRCP Server
	grpcServer := grpc.NewServer()

	//
	RegisterChatServiceServer(grpcServer, &s)

	// start listening on port 9000 for rpc
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000 %v", err)
	}

}
