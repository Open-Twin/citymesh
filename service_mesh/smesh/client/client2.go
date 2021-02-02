package client

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/citymesh/service_mesh/smesh/sidecar"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func Client() {

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Getting all the ips from the ClientCon File
	var ips []string
	file, err := os.Open("../files/clientCon.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//saving all
	for scanner.Scan() {
		line := scanner.Text()
		res:= strings.Split(line,";")
		ips = append(ips, res[1])
	}

	target := ips[0]

	fmt.Println("Connecting to:" +target)

	// Creating a connection with the first ip from the file

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

	var message sidecar.Message
	var message2 sidecar.Message
	var lines []string



	message = sidecar.Message{
		Timestamp:  "16:00:00",
		Location:   "Alte Börse",
		Sensortyp:  "192.168.41.140",
		SensorID:   123,
		SensorData: "S1234",
	}


	for _ = range time.Tick(time.Second * 10) {

		// Sending the Data

		response, err := c.DataFromService(context.Background(), &message)

		if response != nil {
			log.Printf("Response: %s ", response.Reply)

			// Established a connection
			fmt.Println("Sending old data")
			file, err := os.Open("../files/safeData.csv")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				// Trying to push the old data
				res:= strings.Split(line,";")
				message2 = sidecar.Message{
					Timestamp:  res[0],
					Location:   res[1],
					Sensortyp:  res[2],
					SensorID:   123,
					SensorData: res[4],
				}
				response, err := c.DataFromService(context.Background(), &message2)
				if response != nil {
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Response from Sidecar: %s ", response.Reply)

				} else {
					lines = append(lines, line)
				}

			}



			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			f, err := os.OpenFile("../files/safeData.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, value := range lines {
				fmt.Fprintln(f, value)  // print values to f, one per line
			}
			if err != nil {
				fmt.Println(err)
				f.Close()
				return
			}
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}

		} else {
			fmt.Println("Debug: Could not establish a connection")
			fmt.Println("Debug: Saving data locally")

			newIP(ips)

			// Der weitere Teil ist für das Abspeichern von Datensätzen in ein File zuständig, wenn mit keinem Sidecar Verbindung aufgebaut werden kann
			//Saving the data in a local file
			DataSave(message)

		}
		if err != nil {
			log.Printf("Response: %s , ", err)
		}
	}
}

func DataSave(clientMessage sidecar.Message){

	f, err := os.OpenFile("../files/safeData.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString( clientMessage.Location+";"+"Placeholder"+";"+clientMessage.Sensortyp+";"+"Placeholder"+";"+clientMessage.SensorData+";"+"ce_uri_ref"+";"+clientMessage.Timestamp+";"+"\n")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "Debug: Bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func newIP(ips []string){
	fmt.Println("Connecting to new Sidecar")
	for _, newIP := range ips {
		//fmt.Println("Sidecar Ip")
		fmt.Println("Connecting to:"+newIP)
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