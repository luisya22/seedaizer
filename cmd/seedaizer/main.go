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

// TODO: Select databases and send prompt to open ai to generate insert statements.
// TODO: Execute insert statements
/**
Support multiple database engines
Allow users to configure connection details securely

(DONE) Retrieve metadata like tables, columns, datatypes and relationships (foreign keys)

CLI and web client


CLI client:
one command  should save. Other command should execute the prompt
set credentials on .env, json or yaml
select tables via input parameter array
send data to prompt
receive data an print it to: stdout, .sql file or execute directly


Executor:
how to select tables. Maybe using foreign keys to find childrens and then getting some data from parents e.g users dsata because i need user_id
build the prompt
send it to open ai
**/
