package master

import (
	"database/sql"
	//"bufio"
	"fmt"
	"github.com/Shopify/sarama"
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
	panic("implement me")
}

func (s *Server) DataFromSidecar(ctx context.Context, message *CloudEvent) (*MessageReply, error) {
	//log.Printf("Received message body from client %s , %s , %s , %s , %s ", message.IdService, message.Source, message.SpecVersion, message.Type, message.IdService, message.IpSidecar, message.IpSidecar, message.Timestamp, message.Data)
	log.Printf("Received message body from client: %s, %s", message.IpSidecar, message.Source)
	fmt.Println(message.Data)
	//log.Printf( "Received message body from client")
	message.IdSidecar = "sd123"
	message.IpSidecar = "123.123.123.123"
	//msg := message
	SafeToFile(message.IdSidecar, message.IpSidecar, message.Timestamp)
	return &MessageReply{Message: "Sidecar Proxy: Received the message"}, nil
}

func (s *Server) HealthCheck(ctx context.Context, health2 *Health) (*MessageReply, error) {
	/*
		fmt.Println("Sending pings to all registered entries")

		file, err := os.Open("files/master2Config.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		var lines []string
		var report string
		var health sidecar.Health



		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		for _, value := range lines {
			res := strings.Split(value, ";")
			report = report + res[0] + "; \n"
			fmt.Println(res)
			// get()
			// create client for GRPC Server
			var conn *grpc.ClientConn
			// Placeholder until Gateway
			conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
			//conn, err := grpc.Dial(target+":9000", grpc.WithInsecure())

			if erro != nil {
				log.Fatalf("no server connection could be established cause: %v", erro)

			}

			// defer runs after the functions finishes
			defer conn.Close()

			c := sidecar.NewChatServiceClient(conn)

			response, err := c.HealthCheck(context.Background(), &health)

			if response != nil {
				log.Printf("Response from Sidecar: %s ", response.Reply)
			}

			if err != nil {
				fmt.Println(err)
			}

			report = report + response.Reply

		}
		return &MessageReply{Reply: report}, nil
	*/

	return &MessageReply{Message: "Hello"}, nil
}

func (s *Server) GetMessage() string {
	return localmessage
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
	if _, err := os.Stat("files/sqlite-database.db"); os.IsNotExist(err) {
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

	CreateTable(sqliteDatabase) // Create Database Tables
	InsertStorage(sqliteDatabase,sid,ip,tst)
	DisplayStorage(sqliteDatabase)
}



func CreateTable(db *sql.DB) {
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
		row.Scan(&IdService, &IpService, &timestamp )
		log.Println("Metadata: TS:", timestamp, ", ID:", IdService, ", IP:", IpService)
		sqlresponse += timestamp +";"+IdService+ ";" +IpService+ ";"
	}
	return sqlresponse
}

