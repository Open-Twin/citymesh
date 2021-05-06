package sidecar

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	_ "os"
)

const (
	hostname = "Sidecar"
	ipadr = "127.0.0.1"
	rtype = "store"
	listeningport = "9000"
)

func NewServer() {


	// create a TCP Listener on Port 9000
	lis, err := net.Listen("tcp", ":"+listeningport)
	// this how you handle errors in Golang
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 %v", err)
	}
	// this is just a structure that has an interface with needed function SayHello
	s := Server{}
	fmt.Println("Sidecar: GRPC-Server started ... ")
	//create the GRCP Server

	creds, _ := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.MaxSendMsgSize(10*1024*1024), grpc.MaxRecvMsgSize(10*1024*1024))

	// error handling omitted
	//grpcServer := grpc.NewServer()
	RegisterChatServiceServer(grpcServer, &s)
	grpcServer.Serve(lis)
	fmt.Println("Sidecar: after test")
	// start listening on port 9000 for rpc calls
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000 %v", err)
	}

}
