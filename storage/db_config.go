package storage

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//InitDB initialize DB connection
func InitDB(dbtype string, dbname string) (*gorm.DB, error) {
	db, err := gorm.Open(dbtype, dbname)
	if err != nil {
		errors.Wrap(err, "failed to connect to db")
		log.Println(err)
		return nil, err
	}

	db.LogMode(true)
	if !db.HasTable(&Book{}) {
		db.CreateTable(&Book{}).Set("gorm:table_options", "ENGINE=InnoDB")
	}
	return db, nil
}
