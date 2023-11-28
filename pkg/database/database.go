package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Dialect int32

const (
	DbDialectSqlite3 Dialect = 1
)

type DBConfiguration struct {
	Dialect  Dialect `default:"MySQL"`
	Database string
}

// URL returns a connection string for the database.
func (d *DBConfiguration) URL() (url string, err error) {
	switch d.Dialect {
	case DbDialectSqlite3:
		return d.Database, nil
	default:
		return "", fmt.Errorf(" '%v' driver doesn't exist. ", d.Dialect)
	}
}

// Connection return (gorm.DB or error)
func Connection(dbConf DBConfiguration) (db *gorm.DB, cleanup func(), err error) {
	config := &gorm.Config{}
	switch dbConf.Dialect {
	case DbDialectSqlite3:
		db, err = newSqlite(dbConf, config)
	default:
		return nil, func() {}, fmt.Errorf("database dialect %d not supported", dbConf.Dialect)
	}

	if err != nil {
		return
	}
	genDB, err := db.DB()
	cleanup = func() {
		if err := genDB.Close(); err != nil {
			log.Println(err)
		}
	}
	return
}

func newSqlite(dbConf DBConfiguration, conf *gorm.Config) (db *gorm.DB, err error) {
	url, err := dbConf.URL()
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(sqlite.Open(url), conf)
	return
}
