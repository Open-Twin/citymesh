package test

import (
	"database/sql"
	"github.com/Open-Twin/citymesh/complete_mesh/clientSQL"
	//"github.com/Open-Twin/citymesh/complete_mesh/master"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"log"
	"os"
	"testing"
	"time"
)

func TestClientSidecarCom(t *testing.T) {
	go sidecar.NewServer()
	go clientSQL.Client()

	time.Sleep(30 * time.Second)

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

	responsi := clientSQL.DisplayStorage(sqliteDatabase)
	print(responsi)

	if responsi != "" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "")
	}

}
/*
func TestClientSidecarMetadata(t *testing.T) {
	os.Remove("files/sqlite-database.db")
	os.Remove("files/sqlite-database.db")

	go master.Master()

	time.Sleep(30 * time.Second)

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

	responsi := clientSQL.DisplayStorage(sqliteDatabase)
	print(responsi)

	if responsi != "" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "")
	}


}

 */

