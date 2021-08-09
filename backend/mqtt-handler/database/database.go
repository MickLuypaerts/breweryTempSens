package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"mqtt-handler/data"
)

type DBConnection struct {
	l  *log.Logger
	db *sql.DB
}

func NewDBConnection(host string, port int, user string, password string, dbname string, l *log.Logger) (*DBConnection, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		l.Println(err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		l.Println(err)
		return nil, err
	}
	return &DBConnection{db: db, l: l}, nil
}

func (db *DBConnection) InsertProduct(reading *data.Reading) error {

	sqlStatement := `
					INSERT INTO fermantationBarrelTemperture (temperture, barrelId)
					VALUES ($1, $2)`
	_, err := db.db.Exec(sqlStatement, reading.Temperture, reading.BarrelId)

	if err != nil {
		db.l.Println(err)
		return err
	}
	return nil
}
