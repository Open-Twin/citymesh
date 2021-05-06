package sidecar

import (
	_ "errors"
	"github.com/Open-Twin/citymesh/complete_mesh/master"
	_ "github.com/gogo/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	_ "os"
	_ "reflect"
)

const (
	grpcport = "9010"

)

func client(cloudmessage *CloudEvent) {

	log.Println("Debug: Initialised Client")

	// create client for GRPC Server
	var conn *grpc.ClientConn

	conn, error := grpc.Dial("master-service.default.svc.cluster.local"+":"+grpcport, grpc.WithInsecure())

	if error != nil {
		log.Println("no server connection could be established cause: %v", error)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := master.NewChatServiceClient(conn)

	message := master.CloudEvent{
		IdService:   cloudmessage.IdService,
		Source:      cloudmessage.Source,
		SpecVersion: cloudmessage.SpecVersion,
		Type:        cloudmessage.Type,
		Attributes:  nil,
		Data:        nil,
		IdSidecar:   cloudmessage.IdSidecar,
		IpService:   cloudmessage.IpService,
		IpSidecar:   cloudmessage.IpSidecar,
		Timestamp:   cloudmessage.Timestamp,
	}

	response, err := c.DataFromSidecar(context.Background(), &message)

	if err != nil {
		log.Println("No connection could be established")
	}

	if response != nil {
		log.Printf("Response from Server: %s , ", response.Message)

	} else {
		log.Println("Debug: Could not establish a connection")
	}
}
