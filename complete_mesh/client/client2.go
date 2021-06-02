package client

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	_ "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"

	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	"encoding/json"
	_ "errors"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	_ "github.com/golang/protobuf/ptypes"

	"io/ioutil"
	"net/http"
	_ "reflect"

	"github.com/tidwall/gjson"
)

type Warning struct {
	Gkz       string `json:"GKZ"`
	Name      string `json:"Name"`
	Region    string `json:"Region"`
	Warnstufe string `json:"Warnstufe"`
}

func Client() {

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Getting all the ips from the ClientCon File
	var ips []string
	file, err := os.Open("files/clientCon.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//saving all
	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ";")
		ips = append(ips, res[1])
	}

	target := ips[0]

	fmt.Println("Connecting to:" + target)

	// Creating a connection with the first ip from the file

	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}
	conn, err = grpc.Dial(":9000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	print("Test after connection")

	/*// Placeholder until Gateway
	conn, erro := grpc.Dial(":9000", grpc.WithInsecure())

	// real code:
	//conn, err := grpc.Dial(target, grpc.WithInsecure())

	if erro != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)

	}

	// defer runs after the functions finishes
	defer conn.Close()*/

	c := sidecar.NewChatServiceClient(conn)

	var lines []string

	/*message := sidecar.CloudEvent{
		IdService:   "S123",
		Source:      "corona-ampel",
		SpecVersion: "1.1",
		Type:        "JSON",
		Attributes:  nil,
		Data:        nil,
		IdSidecar:   "",
		IpService:   "123.123.123.123",
		IpSidecar:   "",
		Timestamp:   "2021",
	}*/

	for _ = range time.Tick(time.Second * 10) {

		// Sending the Data

		cloudeventmessage := Apiclient()

		response, err := c.DataFromService(context.Background(), &cloudeventmessage)

		//fmt.Println(cloudeventmessage.Data)

		//fmt.Println("Hat es funktioniert?" + response.Message)
		if response != nil {
			log.Printf("Response: %s ", response.Message)

			// Established a connection
			fmt.Println("Sending old data")
			file, err := os.Open("files/saveData.csv")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				// Trying to push the old data
				res := strings.Split(line, ";")
				fmt.Println(res)
				message2 := sidecar.CloudEvent{
					IdService:   res[0],
					Source:      res[1],
					SpecVersion: res[2],
					Type:        res[3],
					Attributes:  nil,
					Data:        nil,
					IdSidecar:   "",
					IpService:   res[4],
					IpSidecar:   "",
					Timestamp:   res[5],
				}
				response, err := c.DataFromService(context.Background(), &message2)
				if response != nil {
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Response from Sidecar: %s ", response.Message)

				} else {
					lines = append(lines, line)
				}

			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			f, err := os.OpenFile("files/saveData.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
			fmt.Println("Debug: Could not establish a connection")
			fmt.Println("Debug: Saving data locally")

			newIP(ips)

			// Der weitere Teil ist für das Abspeichern von Datensätzen in ein File zuständig, wenn mit keinem Sidecar Verbindung aufgebaut werden kann
			//Saving the data in a local file
			DataSave(cloudeventmessage)

		}
		if err != nil {
			log.Printf("Response: %s , ", err)
		}
	}
}

func DataSave(clientMessage sidecar.CloudEvent) {

	f, err := os.OpenFile("files/saveData.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(clientMessage.IdService + ";" + "Placeholder" + ";" + clientMessage.Source + ";" + "Placeholder" + ";" + clientMessage.SpecVersion + ";" + "ce_uri_ref" + ";" + clientMessage.Type + ";" + clientMessage.IdSidecar + ";" + clientMessage.IpService + ";" + clientMessage.IpSidecar + ";" + clientMessage.Timestamp + ";\n")
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

func newIP(ips []string) {
	fmt.Println("Connecting to new Sidecar")
	for _, newIP := range ips {
		//fmt.Println("Sidecar Ip")
		fmt.Println("Connecting to:" + newIP)
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

func URL(url string) (ampelJson string) {
	ampel, erro := http.Get(url)

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	ampelJSON := string(body)

	return ampelJSON
}

func Apiclient() (cloudeventmessage sidecar.CloudEvent) {

	ampel, erro := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	ampelJson := string(body)

	value := gjson.Get(ampelJson, "#.Stand")
	println("Length: ", len(value.Array()))
	println(value.String())

	messages := make([]*dataFormat.Message, 0, 0)
	//0,0 checken

	datenres := gjson.Get(ampelJson, "#.Stand")
	datenres.ForEach(func(key, value gjson.Result) bool {
		stand := value.String()
		warnstufen := make([]*dataFormat.Warnstufen, 0, 0)

		path := "#(Stand==" + stand + ").Warnstufen"
		result := gjson.Get(ampelJson, path)
		result.ForEach(func(key, value gjson.Result) bool {
			var warning Warning
			if err := json.Unmarshal([]byte(value.Raw), &warning); err != nil {
				panic(err)
			}
			warnstufen = append(warnstufen,
				&dataFormat.Warnstufen{
					Region:    warning.Region,
					GKZ:       warning.Gkz,
					Name:      warning.Name,
					Warnstufe: warning.Warnstufe,
				})
			return true
		})
		messages = append(messages,
			&dataFormat.Message{Stand: stand,
				Warnstufen: warnstufen,
			})
		//jsonmessage,err := ProtobufToJSON(messages)

		return true
	})

	print("Message", messages)

	/*marshalMessages, err := ptypes.MarshalAny(&dataFormat.Messages{Message: messages})
	if err != nil {
		panic(err)
	}*/

	//data := sidecar.CloudEvent_ProtoData{ProtoData: marshalMessages}

	//jsondata,err := ProtobufToJSON(messages)

	text := ""
	for _, element := range messages {
		text = text + string(proto.MarshalTextString(element))
		//result, _ := ProtobufToJSON(element)
		//text = text + result

	}

	finishedMessage := sidecar.CloudEvent_TextData{TextData: text}

	cloudeventmessage = sidecar.CloudEvent{
		IdService:   "",
		Source:      "corona-ampel",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        &finishedMessage,
		IdSidecar:   "01",
		IpService:   "192.168.0.10",
		IpSidecar:   "192.168.0.11",
	}

	//fmt.Println(cloudeventmessage.Data)

	/*response, err := c.DataFromService(context.Background(), &cloudeventmessage)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())*/

	return cloudeventmessage

}

func ProtobufToJSON(message dataFormat.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		Indent:          "  ",
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	b, err := marshaler.Marshal(&message)
	return string(b), err
}
