package mysql_conn

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/LeonLow97/internal/pkg/config"
)

func ConnectToMySQL(cfg config.Config) *sql.DB {
	// MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQLConfig.User,
		cfg.MySQLConfig.Password,
		cfg.MySQLConfig.Host,
		cfg.MySQLConfig.Port,
		cfg.MySQLConfig.Database,
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
