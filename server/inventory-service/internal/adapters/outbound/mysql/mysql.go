package outbound_mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectToMySQL() *sql.DB {
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

	// ping mysql database
	if err = conn.Ping(); err != nil {
		log.Fatal("Error pinging mysql database", err)
	}

	return conn
}
