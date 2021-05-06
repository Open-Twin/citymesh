package main

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/master"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	masterGrpcPort = "9010"
	sidecarGrpcPort = "9000"
)

func main() {
	ips := GetIPs()

	var report string

	for i,value := range ips {


		report+= strconv.Itoa(i+1) +") Master: " +value+ "\n"
		// Placeholder until Gateway
		conn, erro := grpc.Dial(value+":"+masterGrpcPort, grpc.WithInsecure())


		if erro != nil {
			report += "No server connection could be established! \n"
			log.Println("no server connection could be established cause: %v", erro)
		}

		// defer runs after the functions finishes
		defer conn.Close()

		c := master.NewChatServiceClient(conn)

		var message master.Health

		message = master.Health{Health: "-"}

		response, err := c.HealthCheck(context.Background(), &message)

		if err != nil {
			log.Println("no server connection could be established cause: %v", err)
		}

		if response != nil{

			fmt.Println(response.Message)
			res := strings.Split(response.Message, ";")

			for i, value := range res {

				report+= strconv.Itoa(i+1) +") Sidecar: " +value+ "\n"

				fmt.Println(value)

				// Placeholder until Gateway
				conn, erro := grpc.Dial(":"+sidecarGrpcPort, grpc.WithInsecure())
				//conn, erro := grpc.Dial(value+":9000", grpc.WithInsecure())

				if erro != nil {
					log.Println("no server connection could be established cause: %v", erro)
					report += "No server connection could be established \n"
				}

				// defer runs after the functions finishes
				defer conn.Close()

				c := sidecar.NewChatServiceClient(conn)

				var message sidecar.Health

				message = sidecar.Health{Health: "abc"}

				responseSide, err := c.HealthCheck(context.Background(), &message)

				if err != nil {
					log.Println("no server connection could be established cause: %v", err)
				}
				report += responseSide.Message
			}
		}
	}
	fmt.Println(report)
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

func GetIPs() []string {
	var ips []string
	file, erre := os.Open("files/meshMasters.csv")
	if erre != nil {
		log.Fatal(erre)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//saving all
	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ";")
		ips = append(ips, res[1])
	}
	return ips
}

