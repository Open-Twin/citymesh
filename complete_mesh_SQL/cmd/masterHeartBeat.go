package main

import (
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/master"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
)

func main() {

	// Placeholder until Gateway
	conn, erro := grpc.Dial(":9010", grpc.WithInsecure())
	// real code:
	//conn, err := grpc.Dial(target, grpc.WithInsecure())

	if erro != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := master.NewChatServiceClient(conn)

	var message master.Health

	message = master.Health{Health: "abc"}

	response, err := c.HealthCheck(context.Background(), &message)

	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", err)
	}

	fmt.Println(response.Message)
	res := strings.Split(response.Message, ";")

	for _,value := range res{
		fmt.Println(value)


		// Placeholder until Gateway
		conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
		//conn, erro := grpc.Dial(value+":9000", grpc.WithInsecure())


		if erro != nil {
			log.Fatalf("no server connection could be established cause: %v", erro)
		}

		// defer runs after the functions finishes
		defer conn.Close()

		c := sidecar.NewChatServiceClient(conn)

		var message sidecar.Health

		message = sidecar.Health{Health: "abc"}

		responseSide, err := c.HealthCheck(context.Background(), &message)

		if err != nil {
			log.Fatalf("no server connection could be established cause: %v", err)
		}

		fmt.Println(responseSide)
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

