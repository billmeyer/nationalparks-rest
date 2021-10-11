package db

import (
	"database/sql"
	"fmt"
)

func OpenDB(username string, password string, hostname string, port int, dbname string) (*sql.DB, error) {
	var connString = GetConnectionString(username, password, hostname, port, dbname)
	return sql.Open("mysql", connString)
}

func GetConnectionString(username string, password string, hostname string, port int, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true", username, password, hostname, port, dbname)
}

