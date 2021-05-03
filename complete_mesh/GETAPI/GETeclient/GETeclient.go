package GETeclient

import (
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	_ "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"

	"encoding/json"
	_ "errors"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat/eladestationen"
	"os"
)

type Properties struct {
	Objectid          int         `json:"OBJECTID"`
	Address           string      `json:"ADDRESS"`
	Bezirk            int         `json:"BEZIRK"`
	City              string      `json:"CITY"`
	Countrycode       string      `json:"COUNTRYCODE"`
	Designation       string      `json:"DESIGNATION"`
	Directpayment     int         `json:"DIRECTPAYMENT"`
	Evseid            string      `json:"EVSEID"`
	HubjectCompatible int         `json:"HUBJECT_COMPATIBLE"`
	Operatorname      string      `json:"OPERATORNAME"`
	Source            string      `json:"SOURCE"`
	SeAnnoCadData     interface{} `json:"SE_ANNO_CAD_DATA"`
}

func Apiclient() (cloudeventmessage sidecar.CloudEvent) {

	ampel, erro := http.Get("https://data.wien.gv.at/daten/geo?service=WFS&request=GetFeature&version=1.1.0&typeName=ogdwien:ELADESTELLEOGD&srsName=EPSG:4326&outputFormat=json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	dataJson := string(body)

	//value := gjson.Get(dataJson, "features")
	//println("Length: ", len(value.Array()))
	//println(value.String())

	messages := make([]*eladestationen.Message, 0, 0)
	//0,0 checken

	datenres := gjson.Get(dataJson, "features")
	datenres.ForEach(func(key, value gjson.Result) bool {
		stand := value.String()
		warnstufen := make([]*eladestationen.Features, 0, 0)

		types := gjson.Get(dataJson, "features.#.type")
		ids := gjson.Get(dataJson, "features.#.id")

		/*println("Length types: ", len(types.Array()))
		println("Text", types.String())

		println("Length ids: ", len(ids.Array()))
		println("Text", ids.String())*/

		geometrytypes := gjson.Get(dataJson, "features.#.geometry.#.type")
		println("Length geometry: ", len(geometrytypes.Array()))
		println("Text", geometrytypes.String())

		geometry := gjson.Get(dataJson, "features.#.geometry")
		geometry.ForEach(func(key, value gjson.Result) bool {
			println("in geometry")

			geometrycoordinates := gjson.Get(dataJson, "features.#.geometry.#.coordinates")

			println("Length coordinates: ", len(geometrycoordinates.Array()))
			println("Text", geometrycoordinates.String())

			return true
		})

		println("Length types: ", len(types.Array()))
		println("Text", types.String())

		println("Length ids: ", len(ids.Array()))
		println("Text", ids.String())

		messages := append(messages,
			&eladestationen.Message{
				Type:          "",
				Totalfeatures: 0,
				Features:      nil,
				Crs:           nil,
			})

		path := "#(Stand==" + stand + ").Warnstufen"
		result := gjson.Get(dataJson, path)
		result.ForEach(func(key, value gjson.Result) bool {
			var properties Properties
			if err := json.Unmarshal([]byte(value.Raw), &properties); err != nil {
				panic(err)
			}
			warnstufen = append(warnstufen,
				&eladestationen.Features{
					Type:         "",
					ID:           "",
					Geometry:     nil,
					GeometryName: "",
					Properties:   nil,
				})
			return true
		})
		messages = append(messages,
			&eladestationen.Message{
				Type:          "",
				Totalfeatures: 0,
				Features:      nil,
				Crs:           nil,
			},
		)
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
