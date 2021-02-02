package sidecar

import (
	"bufio"
	"fmt"
	//"github.com/Open-Twin/CityMesh-ProtoTypes/ServiceMesh/sidecar"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var responseData datastruct

type datastruct []struct {
	string     `json:"Stand"`
	Warnstufen []struct {
		Gkz       string `json:"GKZ"`
		Name      string `json:"Name"`
		Region    string `json:"Region"`
		Warnstufe string `json:"Warnstufe"`
	} `json:"Warnstufen"`
}

func sidecarClient() {

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Zuerst holen wir uns alle Ips aus dem ClientCon File
	// Getting all the ips from the ClientCon File

	var ips []string
	file, err := os.Open("files/clientCon.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		res := strings.Split(line, ";")
		ips = append(ips, res[1])
	}

	target := ips[0]
	fmt.Println(target)

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

	var message sidecar.CloudEvent
	var message2 sidecar.CloudEvent
	var lines []string

	attributes := make([]*CloudEvent_CloudEventAttributeValue, 0, 0)

	msgs := make([]*CloudEvent_CloudEventAttributeValue_Messages, 0, 0)

	msg := make([]*CloudEvent_CloudEventAttributeValue_Message, 0, 0)

	warnstufen := make([]*CloudEvent_CloudEventAttributeValue_Warnstufen, 0, 0)



	attributes = append(attributes, &CloudEvent_CloudEventAttributeValue{
	})


	warnstufen = append(warnstufen, &CloudEvent_CloudEventAttributeValue_Warnstufen{
		Region:    "Tirol",
		GKZ:       "123",
		Name:      "Scheisse",
		Warnstufe: "300",
	})

	msg = append(msg, &CloudEvent_CloudEventAttributeValue_Message{
		Stand:      "2020",
		Warnstufen: warnstufen,
	})

	msgs = append(msgs, &CloudEvent_CloudEventAttributeValue_Messages{
		Message: msg,
	})

	message = CloudEvent{
		IdService:   "s123",
		Source:      "vasco",
		SpecVersion: "1.2",
		Type:        "data",
		Attributes:  attributes,
		Data:        nil,
		IdSidecar:   "",
		IpService:   "",
		IpSidecar:   "",
	}

	fmt.Println(message)

	fmt.Println("Debug: Timer over ")

	for _ = range time.Tick(time.Second * 10) {

		fmt.Println("Debug: Timer over")

		// Sending the Data

		response, err := c.DataFromService(context.Background(), &message)

		if response != nil {
			log.Printf("Response from Server: %s , ", response.Reply)

			/*// Established a connection
			fmt.Println("Sending old data")
			file, err := os.Open("files/safeData.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				// Versuchen die alten Daten zu pushen
				res := strings.Split(line, ";")
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
					log.Printf("Response from Server: %s , ", response.Reply)

				} else {
					lines = append(lines, line)
				}

			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			f, err := os.OpenFile("files/safeData.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, value := range lines {
				fmt.Fprintln(f, value) // print values to f, one per line
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
			fmt.Println("Debug: Oje Keine Verbindung")
			fmt.Println("Debug: Starte Verbindungsversuch zu anderem Sidecar Wiu Wiu")

			newIP(ips)

			// Der weitere Teil ist für das Abspeichern von Datensätzen in ein File zuständig, wenn mit keinem Sidecar Verbindung aufgebaut werden kann
			//Saving the data in a local file
			dataSave()

			*/


		}
		if err != nil {
			log.Printf("Response from Server: %s , ", err)
		}
	}
}

/*func dataSave() {

	fmt.Println("Oje Kein Ausweg, starte Plan B")
	fmt.Println("Der Datensatz wird in ein File gespeichert")

	f, err := os.OpenFile("files/safeData.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString("ce_boolean" + ";" + "ce_integer" + ";" + "192.168.41.132" + ";" + "ce_bytes" + ";" + "S123" + ";" + "ce_uri_ref" + ";" + "Timestamp" + ";" + "\n")
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

func newIP(ips []string) {
	fmt.Println("Beep Boop starte Notfall Loadbalance Protokol")
	for _, newIP := range ips {
		fmt.Println("Sidecar Ip")
		fmt.Println(newIP)
		//conn, erro := grpc.Dial("target", grpc.WithInsecure())
		//defer conn.Close()
		//c := sidecar.NewChatServiceClient(conn)
		//response, err := c.DataFromService(context.Background(), &message2)


			if response != nil {
				fmt.Println("New Target gesichtet")
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Response from Server: %s , ", response.Reply)
				break
			}



	}

}*/
