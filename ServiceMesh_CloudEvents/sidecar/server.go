package sidecar

import (
	"fmt"
	//"golang.org/x/text/message"
	"github.com/Open-Twin/CityMesh-ProtoTypes/ServiceMesh/chat"
	"google.golang.org/grpc"
	"log"
	"net"
)
/*
var Timestamp string
var Location string
var Sensortyp string
var SensorID int32
var SensorData string*/

func sidecarServer() {
	// create a TCP Listener on Port 9001
	lis, err := net.Listen("tcp", ":9001")
	// this how you handle errors in Golang
	if err != nil {
		log.Fatalf("Failed to listen on port 9001 %v", err)
	}
	fmt.Println("TEST")

	// this is just a structure that has an interface with needed function SayHello
	// have a look at /chat/chat2.go for more information
	s := chat.Server{}
	//create the GRCP Server
	grpcServer := grpc.NewServer()

	//
	chat.RegisterChatServiceServer(grpcServer, &s)

	// start listening on port 9000 for rpc
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000 %v", err)
	}

}
