package main

import (
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {

	// Placeholder until Gateway
	conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
	// real code:
	//conn, err := grpc.Dial(target, grpc.WithInsecure())

	if erro != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)

	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := sidecar.NewChatServiceClient(conn)

	var message sidecar.Health

	message = sidecar.Health{Health: "abc"}
	response, err := c.HealthCheck(context.Background(), &message)
	var data string
	if response != nil {
		fmt.Sprintf("Healthcheck Report: %s ", response.Message)
		fmt.Println(data)
		saveHealthCheckCSV(response.Message)

	}

	if err != nil {
		log.Printf("Response from Server: %s , ", err)
	}

}

func saveHealthCheckCSV(data string) {

	f, err := os.OpenFile("../files/healthCheckReport.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func saveHealthCheckSQL(data string) {

	f, err := os.OpenFile("../files/healthCheckReport.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

