package master

import (
	"database/sql"
	//"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"log"
	"os"
)

type Server struct {
}

var Timestamp string
var Location string
var Sensortyp string
var SensorID int32
var SensorData string

var localmessage string

const (
	kafkaConn = "localhost:9092"
	topic     = "topic_test"
)

func (s *Server) mustEmbedUnimplementedChatServiceServer() {
}

func (s *Server) DataFromSidecar(ctx context.Context, message *CloudEvent) (*MessageReply, error) {

	log.Printf("Received message body from client: %s, %s", message.IpSidecar, message.Source)
	fmt.Println(message.Data)

	SafeToFile(message.IpSidecar, message.IdSidecar, message.Timestamp)
	return &MessageReply{Message: "Sidecar Proxy: Received the message"}, nil
}

func (s *Server) HealthCheck(ctx context.Context, health2 *Health) (*MessageReply, error) {
	SafeToFile("192.168.32.12","SIDI123","Heute")
	SafeToFile("192.168.32.14","SIDI456","Morgen")
	// The Healthcheck function returns all the master sidecar ips
	log.Println("Returning all registered entries")
	sidecardata := HealthData()
	fmt.Println(sidecardata)


	return &MessageReply{Message: sidecardata}, nil
}


func initialize() {
	// create producer
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	msg := localmessage

	// publish without goroutene
	publish(msg, producer)

}

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func publish(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}

func SafeToFile(ip string, sid string, tst string) {
	// Saving sidecar metadata in a SQL DB
	if _, err := os.Stat("files/mastermetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating mastermetadata-database.db...")
		file, err := os.Create("files/mastermetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("mastermetadata-database.db created")
	} else {
		log.Println("mastermetadata-database.db already exists")
	}

	sqliteDatabase, error := sql.Open("sqlite3", "files/mastermetadata-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database

	// Creating a table and inserting the metadata
	CreateTable(sqliteDatabase) // Create Database Tables
	InsertStorage(sqliteDatabase,sid,ip,tst)
	DisplayStorage(sqliteDatabase)
}



func CreateTable(db *sql.DB) {
	// Creating a table for the sidecar metadata
	createStorageTableSQL := `CREATE TABLE IF NOT EXISTS metadata (
		"IdService" TEXT NOT NULL PRIMARY KEY,		
		"IpService" TEXT,
		"Timestamp"	TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create metadata table...")
	statement, err := db.Prepare(createStorageTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("metadata table created")
}

func InsertStorage(db *sql.DB, ip string, sid string, tst string) {
	// Inserting the metadata into the DB
	log.Println("Inserting metadata record ...")
	insertStorageSQL := `INSERT OR REPLACE INTO metadata(IdService , IpService, Timestamp) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStorageSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(sid, ip, tst)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func DisplayStorage(db *sql.DB) string {
	// Displaying the entries within the DB
	row, err := db.Query("SELECT * FROM metadata ORDER BY Timestamp")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var sqlresponse string
	for row.Next() { // Iterate and fetch the records from result cursor
		var timestamp string
		var IdService string
		var IpService string
		row.Scan(&IpService, &IdService, &timestamp )
		log.Println("Metadata: TS:", timestamp, ", ID:", IdService, ", IP:", IpService)
		sqlresponse += timestamp +";"+IdService+ ";" +IpService+ ";"
	}
	return sqlresponse
}

func HealthData() string {
	// The Healthdata function replies all enterded sidecar ips to the healthcheck service
	if _, err := os.Stat("files/mastermetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating mastermetadata-database.db...")
		file, err := os.Create("files/mastermetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("mastermetadata-database.db created")
	} else {
		log.Println("mastermetadata-database.db already exists")
	}

	db, error := sql.Open("sqlite3", "files/mastermetadata-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}

	row, err := db.Query("SELECT * FROM metadata ORDER BY Timestamp")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var sqlresponse string

	for row.Next() { // Iterate and fetch the records from result cursor
		var timestamp string
		var IdService string
		var IpService string
		row.Scan(&IpService, &IdService, &timestamp )
		log.Println("Metadata: TS:", timestamp, ", ID:", IdService, ", IP:", IpService)

		sqlresponse+=IpService+";"

	}
	return sqlresponse
}


