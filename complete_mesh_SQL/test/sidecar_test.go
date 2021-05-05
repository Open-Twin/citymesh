package test

import (
	"database/sql"
	csql "github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"log"
	"os"
	"testing"
)

func TestSidecarStorage(t *testing.T) {
	os.Remove("files/sidecarmetadata-database.db")
	if _, err := os.Stat("files/sidecarmetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating sidecarmetadata-database.db...")
		file, err := os.Create("files/sidecarmetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sidecarmetadata-database.db created")
	}else {
		log.Println("sidecarmetadata-database.db already exists")
	}


	sqliteDatabase, error := sql.Open("sqlite3", "files/sidecarmetadata-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	sqliteDatabase.Exec("PRAGMA journal_mode=WAL;")
	csql.CreateTable(sqliteDatabase) // Create Database Tables

	csql.InsertStorage(sqliteDatabase,"Service123ID","192.168.41.32","2021.1")
	csql.InsertStorage(sqliteDatabase,"Service456ID","192.168.41.34","2021.2")

	responsi := csql.DisplayStorage(sqliteDatabase)
	log.Println(responsi)

	if responsi != "2021.1;192.168.41.32;Service123ID|2021.2;192.168.41.34;Service456ID|" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;192.168.41.32;Service123ID|2021.2;192.168.41.34;Service456ID|")
	}
}

func TestSidecarUpdate(t *testing.T) {
	os.Remove("files/sidecarmetadata-database.db")
	if _, err := os.Stat("files/sidecarmetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating sidecarmetadata-database.db...")
		file, err := os.Create("files/sidecarmetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sidecarmetadata-database.db created")
	}else {
		log.Println("sidecarmetadata-database.db already exists")
	}


	sqliteDatabase, error := sql.Open("sqlite3", "files/sidecarmetadata-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	sqliteDatabase.Exec("PRAGMA journal_mode=WAL;")
	csql.CreateTable(sqliteDatabase) // Create Database Tables

	csql.InsertStorage(sqliteDatabase,"Service123ID","192.168.41.32","2021.1")
	csql.InsertStorage(sqliteDatabase,"Service456ID","192.168.41.34","2021.2")

	csql.InsertStorage(sqliteDatabase,"Service123ID","192.168.41.42","2021.15")
	csql.InsertStorage(sqliteDatabase,"Service456ID","192.168.41.44","2021.25")
	responsi := csql.DisplayStorage(sqliteDatabase)
	log.Println(responsi)

	if responsi != "2021.1;192.168.41.32;Service123ID|2021.15;192.168.41.42;Service123ID|2021.2;192.168.41.34;Service456ID|2021.25;192.168.41.44;Service456ID|" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;192.168.41.32;Service123ID|2021.15;192.168.41.42;Service123ID|2021.2;192.168.41.34;Service456ID|2021.25;192.168.41.44;Service456ID|")
	}
}


func TestSidecarGetIps(t *testing.T) {
	ips := csql.GetIp()

	if ips[0] != "192.168.61.192:9000" {
		t.Errorf("Value 0 correct got: %s, want: %s.", ips[0], "192.168.61.192:9000")
	}
	if ips[1] != "192.168.45.172:9000" {
		t.Errorf("Value 0 correct got: %s, want: %s.", ips[1], "192.168.45.172:9000")
	}
	if ips[2] != "192.168.14.123:9000" {
		t.Errorf("Value 0 correct got: %s, want: %s.", ips[2], "192.168.14.123:9000")
	}
	if ips[3] != "192.168.33.156:9000" {
		t.Errorf("Value 0 correct got: %s, want: %s.", ips[2], "192.168.33.156:9000")
	}
}

