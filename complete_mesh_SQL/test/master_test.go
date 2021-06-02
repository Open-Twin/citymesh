package test


import (
	"database/sql"
	csql "github.com/Open-Twin/citymesh/complete_mesh/master"
	"log"
	"os"
	"testing"
)

// Tests if the sqlite db can storage messages
func TestMasterStorage(t *testing.T) {
	os.Remove("files/mastermetadata-database.db")
	if _, err := os.Stat("files/mastermetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating mastermetadata-database.db...")
		file, err := os.Create("files/mastermetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("mastermetadata-database.db created")
	}else {
		log.Println("mastermetadata-database.db already exists")
	}


	sqliteDatabase, error := sql.Open("sqlite3", "files/mastermetadata-database.db") // Open the created SQLite File

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

	if responsi != "2021.1;Service123ID;192.168.41.32;2021.2;Service456ID;192.168.41.34;" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;Service123ID;192.168.41.32;2021.2;Service456ID;192.168.41.34;")
	}
}

// Testing if the sqlite db can handle upper and lower case
func TestMasterStoragePosError(t *testing.T) {
	os.Remove("files/mastermetadata-database.db")
	if _, err := os.Stat("files/mastermetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating mastermetadata-database.db...")
		file, err := os.Create("files/mastermetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("mastermetadata-database.db created")
	}else {
		log.Println("mastermetadata-database.db already exists")
	}


	sqliteDatabase, error := sql.Open("sqlite3", "files/mastermetadata-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	sqliteDatabase.Exec("PRAGMA journal_mode=WAL;")
	csql.CreateTable(sqliteDatabase) // Create Database Tables

	csql.InsertStorage(sqliteDatabase,"Service123ID","192.168.41.32","2021.1")
	csql.InsertStorage(sqliteDatabase,"service123ID","192.168.41.34","2021.2")

	responsi := csql.DisplayStorage(sqliteDatabase)
	log.Println(responsi)

	if responsi != "2021.1;Service123ID;192.168.41.32;2021.2;service123ID;192.168.41.34;" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;Service123ID;192.168.41.32;2021.2;service123ID;192.168.41.34;")
	}
}

func TestMasterUpdate(t *testing.T) {
	os.Remove("files/mastermetadata-database.db")
	if _, err := os.Stat("files/mastermetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating mastermetadata-database.db...")
		file, err := os.Create("files/mastermetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("mastermetadata-database.db created")
	}else {
		log.Println("mastermetadata-database.db already exists")
	}


	sqliteDatabase, error := sql.Open("sqlite3", "files/mastermetadata-database.db") // Open the created SQLite File

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

	if responsi != "2021.1;Service123ID;192.168.41.32;2021.15;Service123ID;192.168.41.42;2021.2;Service456ID;192.168.41.34;2021.25;Service456ID;192.168.41.44;" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;Service123ID;192.168.41.32;2021.15;Service123ID;192.168.41.42;2021.2;Service456ID;192.168.41.34;2021.25;Service456ID;192.168.41.44;")
	}
}


func TestMasterUpdatePosError(t *testing.T) {
	os.Remove("files/mastermetadata-database.db")
	if _, err := os.Stat("files/mastermetadata-database.db"); os.IsNotExist(err) {
		log.Println("Creating mastermetadata-database.db...")
		file, err := os.Create("files/mastermetadata-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("mastermetadata-database.db created")
	}else {
		log.Println("mastermetadata-database.db already exists")
	}


	sqliteDatabase, error := sql.Open("sqlite3", "files/mastermetadata-database.db") // Open the created SQLite File

	if error != nil {
		log.Panic(error)
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	sqliteDatabase.Exec("PRAGMA journal_mode=WAL;")
	csql.CreateTable(sqliteDatabase) // Create Database Tables

	csql.InsertStorage(sqliteDatabase,"Service123ID","192.168.41.32","2021.1")
	csql.InsertStorage(sqliteDatabase,"Service456ID","192.168.41.34","2021.2")

	csql.InsertStorage(sqliteDatabase,"Service123ID","192.168.41.32","2021.1")
	csql.InsertStorage(sqliteDatabase,"Service456ID","192.168.41.34","2021.2")
	responsi := csql.DisplayStorage(sqliteDatabase)
	log.Println(responsi)

	if responsi != "2021.1;Service123ID;192.168.41.32;2021.2;Service456ID;192.168.41.34;" {
		t.Errorf("Value 0 correct got: %s, want: %s.", responsi, "2021.1;Service123ID;192.168.41.32;2021.2;Service456ID;192.168.41.34;")
	}
}