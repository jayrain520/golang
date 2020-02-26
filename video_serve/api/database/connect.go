package database

import (
	"database/sql"
	"testgin/api/config"
)
import _ "github.com/go-sql-driver/mysql"

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", config.Set.UserName+":"+
		config.Set.PassWord+"@tcp("+config.Set.DbConnectAddress+
		":"+config.Set.Port+")/"+config.Set.DbName+"?charset="+config.Set.Charset) //?charset=utf-8 config.Set.Charset
	if err != nil {
		panic(err)
	}
	dbConn.SetMaxOpenConns(int(config.Set.SetMaxOpenConns))
	dbConn.SetMaxIdleConns(int(config.Set.SetMaxIdleConns))
}
