package GETKrankenhaus

import (
	"encoding/json"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat/krankehaus"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"github.com/gogo/protobuf/proto"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
)


type Properties struct {
	Objectid               int         `json:"OBJECTID"`
	KrankenhausBezeichnung string      `json:"KRANKENHAUS_BEZEICHNUNG"`
	Bezeichnung            string      `json:"BEZEICHNUNG"`
	Weblink1               string      `json:"WEBLINK1"`
	Adresse                string      `json:"ADRESSE"`
	Bezirk                 int         `json:"BEZIRK"`
	Telefon                string      `json:"TELEFON"`
	SeAnnoCadData          interface{} `json:"SE_ANNO_CAD_DATA"`
}

type Features struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	GeometryName string `json:"geometry_name"`
	Properties   struct {
		Objectid               int         `json:"OBJECTID"`
		KrankenhausBezeichnung string      `json:"KRANKENHAUS_BEZEICHNUNG"`
		Bezeichnung            string      `json:"BEZEICHNUNG"`
		Weblink1               string      `json:"WEBLINK1"`
		Adresse                string      `json:"ADRESSE"`
		Bezirk                 int         `json:"BEZIRK"`
		Telefon                string      `json:"TELEFON"`
		SeAnnoCadData          interface{} `json:"SE_ANNO_CAD_DATA"`
	} `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Message struct {
	Type          string `json:"type"`
	Totalfeatures int    `json:"totalFeatures"`
	Features      []struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		GeometryName string `json:"geometry_name"`
		Properties   struct {
			Objectid               int         `json:"OBJECTID"`
			KrankenhausBezeichnung string      `json:"KRANKENHAUS_BEZEICHNUNG"`
			Bezeichnung            string      `json:"BEZEICHNUNG"`
			Weblink1               string      `json:"WEBLINK1"`
			Adresse                string      `json:"ADRESSE"`
			Bezirk                 int         `json:"BEZIRK"`
			Telefon                string      `json:"TELEFON"`
			SeAnnoCadData          interface{} `json:"SE_ANNO_CAD_DATA"`
		} `json:"properties"`
	} `json:"features"`
	Crs struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
}

type Coordinates struct {
	Coordinates []float64 `json:"coordinates"`
}




func Apiclient() (cloudeventmessage sidecar.CloudEvent) {

	dataJSON, erro := http.Get("https://data.wien.gv.at/daten/geo?service=WFS&request=GetFeature&version=1.1.0&typeName=ogdwien:STATIONENOGD&srsName=EPSG:4326&outputFormat=json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(dataJSON.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	dataJson := string(body)
	
	messages := make([]*krankehaus.Message, 0, 0)
	
	allfeatures := make([]*krankehaus.Features, 0, 0)
	
	datenres := gjson.Get(dataJson, "features.#.id")
	datenres.ForEach(func(key, value gjson.Result) bool {
		//stand := value.String()
		allproperties := make([]*krankehaus.Properties, 0, 0)
		id := value.String()
		fmt.Println(id)
		

		//propertiespath := "#(ID==" + value.String() + ").properties"
		//propertiespath := "features.#(id=" + id + ").properties"
		propertiespath := "features.#.properties"
		propertysearch := gjson.Get(dataJson, propertiespath)
		propertysearch.ForEach(func(key, value gjson.Result) bool {
			var properties Properties
			fmt.Println(value.String())
			if err := json.Unmarshal([]byte(value.Raw), &properties); err != nil {
				panic(err)
			}

			allproperties = append(allproperties,
				&krankehaus.Properties{
					Objectid:       int64(properties.Objectid),
					Krankenhausbez: properties.KrankenhausBezeichnung,
					Bezeichnung:    properties.Bezeichnung,
					Weblink1:       properties.Weblink1,
					Adresse:        properties.Adresse,
					Bezirk:         int32(properties.Bezirk),
					Telefon:        properties.Telefon,
				})
			return true
		})


		geometries := make([]*krankehaus.Geometry, 0, 0)
		var geometriestype krankehaus.Geometry

		//coordinatespath := "features.#.geometry.coordinates"
		geometriespath := "features.#(id=" + id + ").geometry"
		geometriessearch := gjson.Get(dataJson, geometriespath)
		geometriessearch.ForEach(func(key, value gjson.Result) bool {
			var geometry Geometry
			if err := json.Unmarshal([]byte(value.Raw), &geometry); err != nil {
				panic(err)
			}

			fmt.Println(geometry.Coordinates)


			 = geometry.Type

			var coordinatepair []float64

			geometriespath := "features.#(id=" + id + ").geometry.coordinates"
			geometriessearch := gjson.Get(dataJson, geometriespath)
			geometriessearch.ForEach(func(key, value gjson.Result) bool {
				var coordinates Coordinates
				if err := json.Unmarshal([]byte(value.Raw), &coordinates); err != nil {
					panic(err)
				}

				coordinatepair = coordinates.Coordinates
				return true
			})



			return true
		})




		allgeometries := &krankehaus.Geometry{
			Type:        "",
			Coordinates: nil,
		}


		var features Features
		if err := json.Unmarshal([]byte(value.Raw), &features); err != nil {
			panic(err)
		}

		allfeatures = append(allfeatures,
			&krankehaus.Features{
				Type:         features.Type,
				ID:           features.ID,
				Geometry:     geometries,
				GeometryName: features.GeometryName,
				Properties:   allproperties,
			})


		var message Message
		if err := json.Unmarshal([]byte(value.Raw), &message); err != nil {
			panic(err)
		}

		messages = append(messages,
			&krankenhaus.Message{
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
