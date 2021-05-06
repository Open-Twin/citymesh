package clientSQL

import (
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"

	_ "github.com/gogo/protobuf/proto"
	"io/ioutil"
	"net/http"
	_ "reflect"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
)
type Warning struct {
	Gkz       string `json:"GKZ"`
	Name      string `json:"Name"`
	Region    string `json:"Region"`
	Warnstufe string `json:"Warnstufe"`
}

const (
	Servicename = "Service123"
	Serviceip = "127.0.0.1"
	Rtype = "store"
	gRPCport= "9000"

)


func Client() {

	// Creating and opening a sqlite database
	if _, err := os.Stat("files/sqlite-database.db"); os.IsNotExist(err) {
		log.Println("Creating sqlite-database.db...")
		file, err := os.Create("files/sqlite-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sqlite-database.db created")
	}else {
		log.Println("sqlite-database.db already exists")
	}

	// Opening the created database
	sqliteDatabase, error := sql.Open("sqlite3", "files/sqlite-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database

	// Starting the database in WAL-mode to enable quicker access
	sqliteDatabase.Exec("PRAGMA journal_mode=WAL;")


	// Creating a Table for the database
	CreateTable(sqliteDatabase)

	// Create client for GRPC Server
	var conn *grpc.ClientConn

	// connecting to the sidecar
	conn, _ = grpc.Dial("sidecar-service.default.svc.cluster.local"+":"+gRPCport, grpc.WithInsecure())


	// Defer runs after the functions finishes
	defer conn.Close()
	var c sidecar.ChatServiceClient
	c = sidecar.NewChatServiceClient(conn)

	// For loop which sends data to the sidecars

	for _ = range time.Tick(time.Second * 10) {

		// Gathering data from the open data API
		cloudeventmessage := Apiclient()

		// Sending the gathered data with the the rpc method "DataFromService"
		response, err := c.DataFromService(context.Background(), &cloudeventmessage)

		// Checking if the sidecar received the message
		if response != nil {
			log.Printf("Response: %s ", response.Message)

			// Established a connection
			log.Printf("Sending old data")

			if err != nil {
				log.Println(err)

			} else {

				// If the first message has been received by the sidecar the service tries to send all the stored data rows
				row, err := sqliteDatabase.Query("SELECT * FROM storage ORDER BY Timestamp")
				if err != nil {
					log.Fatal(err)
				}
				defer row.Close()
				// Iterating through all stored entries and sending them to a sidecar
				for row.Next() { // Iterate and fetch the records from result cursor
					var timestamp string
					var IdService string
					var Source string
					var SpecVersion string
					var Type string
					var IdSidecar string
					var IpService string
					var IpSidecar string
					row.Scan(&timestamp, &IdService, &Source, &SpecVersion, &Type, &IdSidecar, &IpService, &IpSidecar )
					log.Println("Storage: ", timestamp, " ", IdService, " ", Source, " ", SpecVersion, " ", Type, " ", IdSidecar, " ", IpService, " ", IpSidecar)

					// Creating a new cloudevent from the sqlite data
					var message2 sidecar.CloudEvent
					message2 = sidecar.CloudEvent{
						IdService:   IdService,
						Source:      Source,
						SpecVersion: SpecVersion,
						Type:        Type,
						Attributes:  nil,
						Data:        nil,
						IdSidecar:   IdSidecar,
						IpService:   IpService,
						IpSidecar:   IpSidecar,
						Timestamp:   timestamp,
					}

					// Sending the new data row
					response, err := c.DataFromService(context.Background(), &message2)

					// If the response from the sidecar is not empty the data row gets deleted from the DB
					if response != nil {
						log.Printf("Response from Sidecar: %s ", response.Message)
						if err != nil {
							log.Fatal(err)
						} else {
							// deleting the entry with the right Timestamp
							DeleteStorage(sqliteDatabase,message2.Timestamp)
						}
					}
				}
			}
		} else {

			// No Connection could be established
			log.Printf("Debug: Could not establish a connection")
			log.Printf("Debug: Saving data locally")



			if error != nil {
				log.Panic(error)
			} else {
				//Saving the data in a local file
				InsertStorage(sqliteDatabase,cloudeventmessage.Timestamp,cloudeventmessage.IdService,cloudeventmessage.Source, cloudeventmessage.SpecVersion,cloudeventmessage.Type,cloudeventmessage.IdSidecar, cloudeventmessage.IpSidecar,cloudeventmessage.IpService)
			}

			DisplayStorage(sqliteDatabase)


		}
		if err != nil {
			log.Printf("Response: %s , ", err)
		}
	}
}

func Apiclient() (cloudeventmessage sidecar.CloudEvent) {

	ampel, erro := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	ampelJson := string(body)

	value := gjson.Get(ampelJson, "#.Stand")
	println("Length: ", len(value.Array()))

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

	data := sidecar.CloudEvent_ProtoData{ProtoData: marshalMessages}

	cloudeventmessage = sidecar.CloudEvent{
		IdService:   "S03",
		Source:      "corona-ampel",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        &data,
		IdSidecar:   "01",
		IpService:   "192.168.0.14",
		IpSidecar:   "-",
		Timestamp: 	time.Now().Format(time.RFC850),
	}

	return cloudeventmessage

}

