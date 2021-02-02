package main

import (
	"fmt"
	"github.com/Open-Twin/CityMesh-ProtoTypes/Jakob_GRPC/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", err)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := sidecar.NewChatServiceClient(conn)

	var message sidecar.ProtoSchlecht

	message = sidecar.ProtoSchlecht{}

	fmt.Println("Wir starten den Ping")

	response, err := c.HealthCheck(context.Background(), &message)
	log.Printf("Response from Sidecar: %s , ", response.Reply)
	if err != nil {
		log.Printf("Response from Sidecar: %s , ", err)
	}

}

/*
func main() {

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", err)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := sidecar.NewChatServiceClient(conn)

	var message sidecar.Message


	message = sidecar.Message{
		Timestamp:  "420.420.420.420",
		Location:   "Transdanubien",
		Sensortyp:  "Reumannplatz",
		SensorID:   123,
		SensorData: "StrasserGasse 8-12",
	}

	fmt.Println(message)



	for i := 0; i < 10; i++ {
		response, err := c.DataFromService(context.Background(), &message)
		log.Printf("Response from Server: %s , ", response.Reply)
		if err != nil {
			log.Printf("Response from Server: %s , ", err)
		}
	}


}


*/
