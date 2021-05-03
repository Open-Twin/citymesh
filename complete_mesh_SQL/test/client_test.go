package test

import (
	"database/sql"
	csql "github.com/Open-Twin/citymesh/complete_mesh/clientSQL"
	"log"
	"os"
	"testing"
)

func TestDataStorage(t *testing.T) {
	os.Remove("files/sqlite-database.db")
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
	csql.CreateTable(sqliteDatabase) // Create Database Tables

	csql.InsertStorage(sqliteDatabase, "2021.1", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")
	csql.InsertStorage(sqliteDatabase, "2021.2", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")

	responsi := csql.DisplayStorage(sqliteDatabase)
	log.Println(responsi)

	if responsi != "2021.1;S123;Corona Ampel;1.0;Json;000;192.168.0.1;123.123.123.1232021.2;S123;Corona Ampel;1.0;Json;000;192.168.0.1;123.123.123.123" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;S123;Corona Ampel;1.0;Json;000;192.168.0.1;123.123.123.1232021.2;S123;Corona Ampel;1.0;Json;000;192.168.0.1;123.123.123.123\n")
	}
}

func TestDataDelete(t *testing.T) {
	os.Remove("files/sqlite-database.db")
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
	csql.CreateTable(sqliteDatabase) // Create Database Tables

	csql.InsertStorage(sqliteDatabase, "2021.1", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")
	csql.InsertStorage(sqliteDatabase, "2021.2", "S123", "Corona Ampel","1.0","Json","000","123.123.123.123","192.168.0.1")

	csql.DeleteStorage(sqliteDatabase,"2021.2")

	responsi := csql.DisplayStorage(sqliteDatabase)
	log.Println(responsi)

	if responsi != "2021.1;S123;Corona Ampel;1.0;Json;000;192.168.0.1;123.123.123.123" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;S123;Corona Ampel;1.0;Json;000;192.168.0.1;123.123.123.123")
	}
}


func TestGetIps(t *testing.T) {
	ips := csql.GetIPs()

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
