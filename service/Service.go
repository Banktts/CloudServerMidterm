package service

import (
	"database/sql"
	"os"
)

func connectSqlDB() *sql.DB {
	var sqlDb, errDb = sql.Open("mysql", os.Getenv("MYSQL_PATH"))
	if errDb != nil {
		panic(errDb.Error())
	}
	return sqlDb
}
