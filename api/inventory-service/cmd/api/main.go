package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var inventoryServicePort = os.Getenv("SERVICE_PORT")

type application struct {
}

func main() {
	app := application{}

	_ = app.connectToDB()

	r := app.routes()

	log.Println("Inventory service listening on port", inventoryServicePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", inventoryServicePort), r); err != nil {
		log.Fatal(err)
	}
}

func (app *application) connectToDB() *sql.DB {
	// MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	// open a connection to MySQL database
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening connection to mysql database in inventory service", err)
	}
	defer conn.Close()

	// ping mysql database
	if err = conn.Ping(); err != nil {
		log.Fatal("Error pinging mysql database", err)
	}

	return conn
}
