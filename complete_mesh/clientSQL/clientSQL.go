package clientSQL

import (
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"github.com/golang/protobuf/ptypes"
	"log"
	"os"

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
	os.Remove("files/sqlite-database.db") // I delete the file to avoid duplicated records.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("files/sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")
	sqliteDatabase, error := sql.Open("sqlite3", "files/sqlite-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database

	createTable2(sqliteDatabase) // Create Database Tables

	insertStorage2(sqliteDatabase, "2021.0", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")
	insertStorage2(sqliteDatabase, "2021.1", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")
	insertStorage2(sqliteDatabase, "2021.2", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")

	// DISPLAY INSERTED RECORDS
	displayStorage2(sqliteDatabase)
	/*
	createTable(sqliteDatabase) // Create Database Tables

	// INSERT RECORDS
	insertStorage(sqliteDatabase, "0001", "Liana Kim", "Bachelor")
	insertStorage(sqliteDatabase, "0002", "Glen Rangel", "Bachelor")
	insertStorage(sqliteDatabase, "0003", "Martin Martins", "Master")
	insertStorage(sqliteDatabase, "0004", "Alayna Armitage", "PHD")
	insertStorage(sqliteDatabase, "0005", "Marni Benson", "Bachelor")
	insertStorage(sqliteDatabase, "0006", "Derrick Griffiths", "Master")
	insertStorage(sqliteDatabase, "0007", "Leigh Daly", "Bachelor")
	insertStorage(sqliteDatabase, "0008", "Marni Benson", "PHD")
	insertStorage(sqliteDatabase, "0009", "Klay Correa", "Bachelor")

	// DISPLAY INSERTED RECORDS
	displayStorage(sqliteDatabase)
	*/

	/*

	// get()
	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Getting all the ips from the ClientCon File
	var ips []string
	file, erre := os.Open("files/clientCon.csv")
	if erre != nil {
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

	// Placeholder until Gateway
	conn, erro := grpc.Dial(":9001", grpc.WithInsecure())
	// real code:
	//conn, err := grpc.Dial(target, grpc.WithInsecure())

	if erro != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)

	}

	// defer runs after the functions finishes
	defer conn.Close()

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
	/*
	for _ = range time.Tick(time.Second * 10) {

		// Sending the Data

		cloudeventmessage := Apiclient()

		response, err := c.DataFromService(context.Background(), &cloudeventmessage)

		fmt.Println(cloudeventmessage.Data)

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
			log.Printf("Debug: Could not establish a connection")
			log.Printf("Debug: Saving data locally")

			newIP(ips)

			//Saving the data in a local file
			DataSave(cloudeventmessage)

		}
		if err != nil {
			log.Printf("Response: %s , ", err)
		}
	}

	 */
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

	data := sidecar.CloudEvent_ProtoData{ProtoData: marshalMessages}

	cloudeventmessage = sidecar.CloudEvent{
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

	//fmt.Print(cloudeventmessage.Data)

	//fmt.Println(cloudeventmessage.Data)

	/*response, err := c.DataFromService(context.Background(), &cloudeventmessage)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())*/

	return cloudeventmessage

}

func createTable(db *sql.DB) {
	createStorageTableSQL := `CREATE TABLE storage (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create storage table...")
	statement, err := db.Prepare(createStorageTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("storage table created")
}

func createTable2(db *sql.DB) {
	createStorageTableSQL := `CREATE TABLE storage (
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

func insertStorage2(db *sql.DB, Timestamp string, IdService string, Source string,SpecVersion string,Type string,IdSidecar string,IpSidecar string,IpService string) {
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

func displayStorage2(db *sql.DB) {
	row, err := db.Query("SELECT * FROM storage ORDER BY Timestamp")
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
	}
}

// We are passing db reference connection from main to our method with other parameters
func insertStorage(db *sql.DB, code string, name string, program string) {
	log.Println("Inserting storage record ...")
	insertStorageSQL := `INSERT INTO storage(code, name, program) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStorageSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayStorage(db *sql.DB) {
	row, err := db.Query("SELECT * FROM storage ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Storage: ", code, " ", name, " ", program)
	}
}

