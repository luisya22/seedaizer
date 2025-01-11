package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/luisya22/seedaizer/dbscanner"
)

var (
	version = "0.0.1"
)

func main() {
	dsn := "root:Local1234567890@tcp(localhost:3306)/mcs_ticketing_system"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	scner := dbscanner.DBScanner{Db: db}

	schema, err := scner.Scan()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(schema)
}

/**
Support multiple database engines
Allow users to configure connection details securely

CLI and web client


CLI client:
one command  should save. Other command should execute the prompt
set credentials on .env, json or yaml
select tables via input parameter array
send data to prompt
receive data and print it to: stdout, .sql file or execute directly


Pregenerate a large amount of data. The AI will only map a column
to a field type.

**/
