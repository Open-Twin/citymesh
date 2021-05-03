package sidecar

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tatsushid/go-fastping"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Server struct {
}

func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	//panic("implement me")
}

var Timestamp string
var Location string
var Sensortyp string
var SensorID int32
var SensorData string

var localmessage string

var msg *CloudEvent

func (s *Server) DataFromService(ctx context.Context, message *CloudEvent) (*MessageReply, error) {
	//log.Printf("Received message body from client %s , %s , %s , %s , %s ", message.IdService, message.Source, message.SpecVersion, message.Type, message.IdService, message.IpSidecar, message.IpSidecar, message.Timestamp, message.Data)
	//fmt.Printf( "Received message body from client")
	//log.Printf("Received: %s", message.String())
	message.IdSidecar = "SD01"
	message.IpSidecar = "123.123.123.123"
	msg = message
	fmt.Println("MessageID:"+message.IdService)
	SafeToFile(message.IdService, message.IpService, message.Timestamp)
	go client(msg)
	return &MessageReply{Message: "Sidecar Proxy: Received the message"}, nil
}

func (s *Server) HealthCheck(ctx context.Context, health *Health) (*MessageReply, error) {
	fmt.Println("Sending pings to all registered entries")

	file, err := os.Open("files/sidecarConfig.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	var response string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	for _, value := range lines {
		res := strings.Split(value, ";")
		fmt.Println(res)
		p := fastping.NewPinger()
		ra, err := net.ResolveIPAddr("ip4:icmp", "www.google.com")
		//ra, err := net.ResolveIPAddr("ip4:icmp", res[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		p.AddIPAddr(ra)
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			//response = fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)

			response = fmt.Sprintf("%s ID: %s IP Addr: %s receive, RTT: %v; \n ", response, res[0], addr.String(), rtt)

		}
		p.OnIdle = func() {

		}
		err = p.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	return &MessageReply{Message: response}, nil
}

func SafeToFile(ip string, sid string, tst string) {
	if _, err := os.Stat("files/sqlite-database.db"); os.IsNotExist(err) {
		log.Println("Creating sidecarmetadata-database.db...")
		file, err := os.Create("files/sidecarmetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sidecarmetadata-database.db created")
	} else {
		log.Println("sidecarmetadata-database.db already exists")
	}

	sqliteDatabase, error := sql.Open("sqlite3", "files/sidecarmetadata-database.db") // Open the created SQLite File

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
		//log.Println(sqlresponse)
	}
	return sqlresponse
}

