package sidecar

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/CityMesh-ProtoTypes/ServiceMesh/chat"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	_ "os"
	"strings"
)

func client(mssgg *CloudEvent) {
	fmt.Println("TEST")

	var ClientTimestamp = mssgg.Timestamp
	var ClientSource = mssgg.Source
	var ClientType = mssgg.Type
	var ClientSidecarIP = mssgg.IpSidecar
	var ClientSidecarID = mssgg.IdSidecar
	var ClientServiceIP = mssgg.IpService
	var ClientServiceID = mssgg.IdService

	// create client for GRPC Server
	var conn *grpc.ClientConn

	//
	//
	// Zuerst holen wir uns alle Ips aus dem ClientCon File
	//
	//

	var ips []string
	file, err := os.Open("./files/sidecarCon.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Linie:" + line)
		res := strings.Split(line, ";")

		ips = append(ips, res[1])

	}

	target := ips[0]
	fmt.Println(target)

	fmt.Println("Variablen")

	conn, erro := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	var message chat.CloudEvent

	message = chat.CloudEvent{
		IdService:   ClientServiceID,
		Source:      ClientSource,
		SpecVersion: "123",
		Type:        ClientType,
		Attributes:  nil,
		Data:        nil,
		IdSidecar:   ClientSidecarID,
		IpService:   ClientServiceIP,
		IpSidecar:   ClientSidecarIP,
		Timestamp:   ClientTimestamp,
	}

	fmt.Println("Message erstellt und am Schicken")

	response, err := c.DataFromService(context.Background(), &message)

	if response != nil {
		log.Printf("Response from Server: %s , ", response.Message)

	} else {
		fmt.Println("Oje Keine Verbindung")
		fmt.Println("Starte Verbindungsversuch zu anderem Sidecar Wiu Wiu")

		fmt.Println("Beep Boop starte Notfall Loadbalance Protokol")
		for _, newIP := range ips {
			fmt.Println("New Sidecar Ip:")
			fmt.Println(newIP)
			//conn, erro := grpc.Dial("target", grpc.WithInsecure())
			//defer conn.Close()
			//c := sidecar.NewChatServiceClient(conn)
			//response, err := c.DataFromService(context.Background(), &message2)

			/*
				if response != nil {
					fmt.Println("New Target gesichtet")
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Response from Server: %s , ", response.Reply)
					break
				}

			*/
		}
	}
}
