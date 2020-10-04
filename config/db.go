package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "learngo"
)

var Db *sql.DB

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	Db, err = sql.Open("postgres", psqlInfo)
	CheckErr(err)

	err = Db.Ping()
	CheckErr(err)

	fmt.Println("Successfully connected to database!")

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
