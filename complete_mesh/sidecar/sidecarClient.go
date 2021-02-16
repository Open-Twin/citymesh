package sidecar

import (
	"bufio"
	"fmt"
	//"github.com/Open-Twin/citymesh/complete_mesh/chat"
	//"github.com/Open-Twin/citymesh/service_mesh/smesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	_ "os"
	"strings"
	"encoding/json"
	_ "errors"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/gogo/protobuf/proto"
	"io/ioutil"
	_ "reflect"
	"github.com/tidwall/gjson"
)

type Warning struct {
	Gkz       string `json:"GKZ"`
	Name      string `json:"Name"`
	Region    string `json:"Region"`
	Warnstufe string `json:"Warnstufe"`
}

func client(cloudmessage *CloudEvent) {
	fmt.Println("Debug: Initialised Client")

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
		res := strings.Split(line, ";")

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




	c := NewChatServiceClient(conn)

	/*message := chat.CloudEvent{
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
	}*/

	cloudeventmessage := Apiclient()

	response, err := c.DataFromService(context.Background(), &cloudeventmessage)


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

func Apiclient() (cloudeventmessage CloudEvent) {


	ampel, erro := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
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

		return true
	})



	marshalMessages, err := ptypes.MarshalAny(&dataFormat.Messages{Message: messages})
	if err != nil {
		panic(err)
	}

	data := CloudEvent_ProtoData{ProtoData: marshalMessages}

	cloudeventmessage = CloudEvent{
		IdService:   "",
		Source:      "corona-ampel",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        &data,
		IdSidecar:   "01",
		IpService:   "192.168.0.10",
		IpSidecar:   "192.168.0.11",
	}

	fmt.Print(cloudeventmessage)

	//fmt.Println(cloudeventmessage.Data)

	/*response, err := c.DataFromService(context.Background(), &cloudeventmessage)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())*/

	return cloudeventmessage

}
