package GETAPI

import (
	"encoding/json"
	_ "errors"
	"fmt"
	_ "github.com/Open-Twin/CityMesh-ProtoTypes/ServiceMesh/chat"
	"github.com/Open-Twin/CityMesh-ProtoTypes/ServiceMesh/dataFormat"
	//"golang.org/x/net/context"
	_ "github.com/gogo/protobuf/proto"
	"io/ioutil"
	"net/http"

	"google.golang.org/grpc"
	"log"
	"os"
	_ "reflect"

	"github.com/tidwall/gjson"
)

type Warning struct {
	Gkz       string `json:"GKZ"`
	Name      string `json:"Name"`
	Region    string `json:"Region"`
	Warnstufe string `json:"Warnstufe"`
}

func Apiclient() {

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", err)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	//c := chat.NewChatServiceClient(conn)

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

		path := "#(Stand=="+stand+").Warnstufen"
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

	fmt.Println(messages)

	//attributes := make([]*sidecar.CloudEvent_CloudEventAttributeValue, 0, 0)


	/*marshalMessages, err := ptypes.MarshalAny(&messages)
	if err != nil {
		panic(err)
	}

	protodata := append(protodata, &sidecar.CloudEvent_ProtoData{ProtoData: marshalMessages})



	cloudeventmessage := sidecar.CloudEvent{
		IdService:   "",
		Source:      "corona-ampel",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        protodata,
		IdSidecar:   "01",
		IpService:   "192.168.0.10",
		IpSidecar:   "192.168.0.11",
	}



	for _, element := range messages {
		response, err := c.SayHello(context.Background(), &element)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", response.GetMessage())
	}*/
}

func APIData() ([]byte, error){
	ampel, erro := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	return body, err
}
