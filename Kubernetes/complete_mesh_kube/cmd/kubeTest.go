package main

import (
	"fmt"
	mast "github.com/Open-Twin/citymesh/complete_mesh/master"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const (
	hostname = "Master"
	masterip = "127.0.0.1"
	rtype = "store"
	grpcport = "20000"
)

func main() {

	// create a TCP Listener on Port 9001
	lis, err := net.Listen("tcp", ":"+grpcport)
	// this how you handle errors in Golang
	if err != nil {
		log.Fatalf("Failed to listen on port 9001 %v", err)
	}
	// this is just a structure that has an interface with needed function SayHello
	s := mast.Server{}
	fmt.Println("Master: GRPC-Server started")
	//create the GRCP Server

	creds, _ := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.MaxSendMsgSize(10*1024*1024), grpc.MaxRecvMsgSize(10*1024*1024))
	// error handling omitted

	//grpcServer := grpc.NewServer()
	mast.RegisterChatServiceServer(grpcServer, &s)
	grpcServer.Serve(lis)
	// start listening on port 9001 for rpc
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9001 %v", err)
	}
}

