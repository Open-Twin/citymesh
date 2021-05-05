package clientSQL

import (
	"bufio"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/DDNS"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
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

const (
	Servicename = "Service123"
	Serviceip = "127.0.0.1"
	Rtype = "store"
	gRPCport= "9000"

)


func Client() {

	// Registering the service at the DDNS cluster
	var ipsidecar net.IP
	ipsidecar = conDNS()
	fmt.Println(ipsidecar)

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

	// Getting all the ips from the ClientCon File
	var ips []string
	ips = GetIPs()

	target := ips[0]

	log.Printf("Connecting to:" + target)


	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}
	conn, err = grpc.Dial(":"+gRPCport, grpc.WithTransportCredentials(creds))

	//conn, err = grpc.Dial(sidecarip+":"+gRPCport, grpc.WithTransportCredentials(creds))

	//conn, err := grpc.Dial(":9000", grpc.WithInsecure())


	if err != nil {
		log.Fatalf("no server connection could be established cause: %v", err)

	}

	// Defer runs after the functions finishes
	defer conn.Close()

	c := sidecar.NewChatServiceClient(conn)

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
				log.Panic(err)

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

			// Trying other Ip addresses
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


func newIP(ips []string) error {
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
					return err
				}
				log.Printf("Response from Server: %s , ", response.Reply)
				break
			}

		*/

	}
	return nil
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
func OpenSQL() (*sql.DB,error){
	// Method for opening or creating a sqlite DB
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

	return sqliteDatabase, error
}

func CreateTable(db *sql.DB) {
	// Createing a SQL Table which holds all the data the service is sending
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
	// Inserting service data which could not be send into the sqlite DB
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

	// Deleting existing entries by their Timestamp
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

	// Diplaying and returning all present entries
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

func conDNS() net.IP {

	// Registering the service at the ddns Network
	var sidecarip net.IP
	// Service erreichbar + Ip from Service + Type des Services Store etc...
	ddns.Register(Servicename,Serviceip,Rtype)

	ip, err := ddns.Query("Sidecar")
	fmt.Println("IP API:")
	fmt.Println(ip)
	if err != nil || len(ip) == 0 {
		fmt.Println("Panik im Schweinestall.")
	}else {
		sidecarip = ip[0]
		fmt.Println(sidecarip)
	}
	return sidecarip
}