package psql

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	User        string
	Password    string
	Host        string
	Port        string
	DBName      string
	MaxConn     int
	MaxIdleConn int
}

func (p *Database) ConnString() string {

	schema := "postgres://"

	if p.User != "" {
		schema += p.User
	} else {
		schema += "postgres"
	}

	if p.Password != "" {
		schema += ":" + p.Password
	}

	if p.Host != "" {
		schema += "@" + p.Host
	} else {
		schema += "@localhost"
	}

	if p.Port != "" {
		schema += ":" + p.Port
	}

	if p.DBName != "" {
		schema += "/" + p.DBName
	}

	return schema
}

func Init(p *Database) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", p.ConnString())
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(p.MaxConn)
	db.SetMaxIdleConns(p.MaxIdleConn)
	return db, err
}

func SqlDateTime(d string) string {
	// '0000-00-00' or '0000-00-00 00:00:00'
	Moscow, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.Parse(time.RFC3339, d)
	return t.In(Moscow).String()[:19]
}
