package clientSQL

import (
	"database/sql"
	"log"
	"os"
)

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

