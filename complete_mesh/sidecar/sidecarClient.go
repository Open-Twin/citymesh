package sidecar

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/chat"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	_ "os"
	"strings"
)

func client(mssgg *Message) {
	fmt.Println("Debug: Initialised Client")
	var ClientTimestamp = mssgg.Timestamp
	var ClientLocation = mssgg.Location
	var ClientSensortyp = mssgg.Sensortyp
	var ClientSensorID = mssgg.SensorID
	var ClientSensorData = mssgg.SensorData


	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Zuerst holen wir uns alle Ips aus dem ClientCon File

	var ips []string
	file, err := os.Open("../files/sidecarCon.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		res:= strings.Split(line,";")

		ips = append(ips, res[1])

	}

	target := ips[0]
	fmt.Println(target)


	conn, erro := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	var message chat.Message

	message = chat.Message{
		Timestamp:  ClientTimestamp,
		Location:   ClientLocation,
		Sensortyp:  ClientSensortyp,
		SensorID:   ClientSensorID,
		SensorData: ClientSensorData,
	}

	response, err := c.DataFromSidecar(context.Background(), &message)

	if response != nil {
		log.Printf("Response from Server: %s , ", response.Reply)

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