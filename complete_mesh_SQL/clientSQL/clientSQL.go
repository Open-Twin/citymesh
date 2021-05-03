package clientSQL

import (
	"bufio"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"strings"
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


func Client() {
	//os.Remove("files/sqlite-database.db") // I delete the file to avoid duplicated records.
	//os.Create("files/sqlite-database.db")

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


	sqliteDatabase, error := sql.Open("sqlite3", "files/sqlite-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	sqliteDatabase.Exec("PRAGMA journal_mode=WAL;")
	CreateTable(sqliteDatabase) // Create Database Tables

	//insertStorage(sqliteDatabase, "", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")
	//insertStorage(sqliteDatabase, "2021.1", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")
	//insertStorage(sqliteDatabase, "2021.2", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")

	// DISPLAY INSERTED RECORDS
	//displayStorage(sqliteDatabase)

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Getting all the ips from the ClientCon File
	var ips []string
	ips = GetIPs()

	target := ips[0]

	log.Printf("Connecting to:" + target)

	// Creating a connection with the first ip from the file

	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}
	conn, err = grpc.Dial(":9000", grpc.WithTransportCredentials(creds))


	// Placeholder until Gateway
	//conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
	// real code:
	//conn, err := grpc.Dial(target, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", err)

	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := sidecar.NewChatServiceClient(conn)

	for _ = range time.Tick(time.Second * 10) {

		// Sending the Data

		cloudeventmessage := Apiclient()

		response, err := c.DataFromService(context.Background(), &cloudeventmessage)

		if response != nil {
			log.Printf("Response: %s ", response.Message)

			// Established a connection
			log.Printf("Sending old data")

			if error != nil {
				log.Panic(error)

			} else {

				row, err := sqliteDatabase.Query("SELECT * FROM storage ORDER BY Timestamp")
				if err != nil {
					log.Fatal(err)
				}
				defer row.Close()
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

					message2 := sidecar.CloudEvent{
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

					response, err := c.DataFromService(context.Background(), &message2)

					if response != nil {
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("Response from Sidecar: %s ", response.Message)
						if error != nil {
							log.Panic(error)

						} else {
							DeleteStorage(sqliteDatabase,message2.Timestamp)
						}


					}
				}

			}

		} else {
			log.Printf("Debug: Could not establish a connection")
			log.Printf("Debug: Saving data locally")

			newIP(ips)

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

func GetIPs() []string {
	var ips []string
	file, erre := os.Open("files/clientCon.csv")
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
		IpSidecar:   "192.168.0.11",
		Timestamp: 	time.Now().Format(time.RFC850),
	}

	return cloudeventmessage

}


func CreateTable(db *sql.DB) {
	createStorageTableSQL := `CREATE TABLE IF NOT EXISTS storage (
		"Timestamp" TEXT NOT NULL PRIMARY KEY,		
		"IdService" TEXT,
		"Source" TEXT,
		"SpecVersion" TEXT,
		"Type" TEXT,
		"IdSidecar" TEXT,
		"IpService" TEXT,
		"IpSidecar"	TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create storage table...")
	statement, err := db.Prepare(createStorageTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("storage table created")
}

func InsertStorage(db *sql.DB, Timestamp string, IdService string, Source string,SpecVersion string,Type string,IdSidecar string,IpSidecar string,IpService string) {

	log.Println("Inserting storage record ...")
	insertStorageSQL := `INSERT INTO storage(Timestamp , IdService, Source, SpecVersion, Type, IdSidecar, IpService, IpSidecar ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertStorageSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(Timestamp, IdService, Source, SpecVersion, Type, IdSidecar, IpService, IpSidecar)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func DeleteStorage(db *sql.DB, tst string) {
	log.Println("Delete storage record ...")
	deleteStorageSQL := `DELETE FROM storage WHERE Timestamp = (?)`

	statement, err := db.Prepare(deleteStorageSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(tst)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Deleting from storage record ...")
}


func DisplayStorage(db *sql.DB) string {
	row, err := db.Query("SELECT * FROM storage ORDER BY Timestamp")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var sqlresponse string
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

		sqlresponse += timestamp +";"+IdService+ ";" +Source+ ";"  +SpecVersion+ ";" +Type+ ";" +IdSidecar+ ";" +IpService+ ";" +IpSidecar
		//log.Println(sqlresponse)
	}
	return sqlresponse
}



