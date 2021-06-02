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
type Features struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float32 `json:"coordinates"`
	} `json:"geometry"`
	GeometryName string `json:"geometry_name"`
	Properties   struct {
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
	} `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

type Message struct {
	Type          string `json:"type"`
	Totalfeatures int    `json:"totalFeatures"`
	Features      []struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float32 `json:"coordinates"`
		} `json:"geometry"`
		GeometryName string `json:"geometry_name"`
		Properties   struct {
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
		} `json:"properties"`
	} `json:"features"`
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

	//allproperties := make([]*eladestationen.Features, 0, 0)
	//0,0 checken
	allfeatures := make([]*eladestationen.Features, 0, 0)

	//allcoordinates := make([]*Coordinates,0,0)

	datenres := gjson.Get(dataJson, "features.#.id")
	datenres.ForEach(func(key, value gjson.Result) bool {
		//stand := value.String()
		allproperties := make([]*eladestationen.Properties, 0, 0)
		allgeometries := make([]*eladestationen.Geometry, 0, 0)
		id := value.String()
		fmt.Println(id)
		/*types := gjson.Get(dataJson, "features.#.type")
		ids := gjson.Get(dataJson, "features.#.id")*/

		/*println("Length types: ", len(types.Array()))
		println("Text", types.String())

		println("Length ids: ", len(ids.Array()))
		println("Text", ids.String())*/

		//propertiespath := "#(ID==" + value.String() + ").properties"
		//propertiespath := "features.#(id=" + id + ").properties"
		propertiespath := "features.#.properties"
		propertysearch := gjson.Get(dataJson, propertiespath)
		propertysearch.ForEach(func(key, value gjson.Result) bool {
			var properties Properties
			//fmt.Println(value.String())
			if err := json.Unmarshal([]byte(value.Raw), &properties); err != nil {
				panic(err)
			}

			allproperties = append(allproperties,
				&eladestationen.Properties{
					Objectid:          int64(properties.Objectid),
					Address:           properties.Address,
					Bezirk:            int32(properties.Bezirk),
					City:              properties.City,
					Countrycode:       properties.Countrycode,
					Designation:       properties.Designation,
					Directpayment:     int32(properties.Directpayment),
					Evseid:            properties.Evseid,
					HubjectCompatible: int64(properties.HubjectCompatible),
					Operatorname:      properties.Operatorname,
					Source:            properties.Source,
					//SeAnnoCadData:     ptypes.MarshalAny(properties.SeAnnoCadData),
				})
			//fmt.Println(allproperties)
			return true
		})

		//geometryname := gjson.Get(dataJson, "features.#.geometry_name")
		//geometriespath := "features.#.geometry.type"

		geometry := gjson.Get(dataJson, "features.#.geometry")
		geometry.ForEach(func(key, value gjson.Result) bool {
			var geometry Geometry

			if err := json.Unmarshal([]byte(value.Raw), &geometry); err != nil {
				panic(err)
			}

			allgeometries = append(allgeometries,
				&eladestationen.Geometry{
					Type:        geometry.Type,
					Coordinates: geometry.Coordinates,
				})

			fmt.Println(allgeometries)
			return true
		})

		/*var features Features
		featurespath := "features"
		result := gjson.Get(dataJson, featurespath)
		fmt.Println("Result:", result)
		if err := json.Unmarshal([]byte(result.String()), &features); err != nil {
			panic(err)
		}

		allfeatures = append(allfeatures,
			&eladestationen.Features{
				Type:         features.Type,
				ID:           features.ID,
				Geometry:     allgeometries,
				GeometryName: features.GeometryName,
				Properties:   allproperties,
			})*/

		var message Message
		featurespath := ""
		result := gjson.Get(dataJson, featurespath)
		fmt.Println("Result:", result)
		if err := json.Unmarshal([]byte(value.Raw), &message); err != nil {
			panic(err)
		}

		messages = append(messages,
			&eladestationen.Message{
				Type:          message.Type,
				Totalfeatures: int64(message.Totalfeatures),
				Features:      allfeatures,
			})

		/*println("Length types: ", len(types.Array()))
		println("Text", types.String())

		println("Length ids: ", len(ids.Array()))
		println("Text", ids.String())*/

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
		//fmt.Println(text)
		//result, _ := ProtobufToJSON(element)
		//text = text + result

	}

	fmt.Println(messages)

	finishedMessage := sidecar.CloudEvent_TextData{TextData: text}

	fmt.Println(finishedMessage.TextData)

	cloudeventmessage = sidecar.CloudEvent{
		IdService:   "",
		Source:      "E-Ladestationen",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        &finishedMessage,
		IdSidecar:   "01",
		IpService:   "192.168.0.10",
		IpSidecar:   "192.168.0.11",
	}

	fmt.Println(cloudeventmessage.Data)

	/*response, err := c.DataFromService(context.Background(), &cloudeventmessage)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())*/

	return cloudeventmessage

}
