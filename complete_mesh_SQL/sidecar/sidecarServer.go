package sidecar

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	// this is just a structure that has an interface with needed function SayHello
	// have a look at /master/chat2.go for more information
	s := Server{}
	fmt.Println("0")
	fmt.Println("Sidecar: GRPC-Server started ... ")
	//create the GRCP Server
	fmt.Println("1")
	creds, _ := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.MaxSendMsgSize(10*1024*1024), grpc.MaxRecvMsgSize(10*1024*1024))
	// error handling omitted
	fmt.Println("2")
	//grpcServer := grpc.NewServer()
	RegisterChatServiceServer(grpcServer, &s)
	grpcServer.Serve(lis)
	fmt.Println("Sidecar: after test")
	// start listening on port 9000 for rpc
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000 %v", err)
	}

}
