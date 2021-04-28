package sidecar

import (
	"bufio"
	_ "errors"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/master"
	_ "github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/credentials"

	//"github.com/Open-Twin/citymesh/service_mesh/smesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	_ "os"
	_ "reflect"
	"strings"
)

func client(cloudmessage *CloudEvent) {
	fmt.Println("Debug: Initialised Client")

	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Zuerst holen wir uns alle Ips aus dem ClientCon File

	var ips []string
	file, err := os.Open("files/sidecarCon.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ";")

		ips = append(ips, res[1])

	}

	target := ips[0]
	fmt.Println(target)

	creds, _ := credentials.NewClientTLSFromFile("cert/server.crt", "")
	conn, error := grpc.Dial(":9010", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", error)
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
		log.Fatal(err)
	}

	if response != nil {
		log.Printf("Response from Server: %s , ", response.Message)

	} else {
		fmt.Println("Debug: Could not establish a connection")

		for _, newIP := range ips {
			fmt.Println("Debug: New Sidecar Ip:")
			fmt.Println(newIP)
			//conn, erro := grpc.Dial("target", grpc.WithInsecure())
			//defer conn.Close()
			//c := sidecar.NewChatServiceClient(conn)
			//response, err := c.DataFromSidecar(context.Background(), &message2)

			/*
				if response != nil {
					fmt.Println("New Target gesichtet")
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Response from Server: %s ", response.Reply)
					break
				}

			*/
		}
	}
}
