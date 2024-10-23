package db

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB = DBClassify{}

type (
	DBClassify struct {
		GORM *gorm.DB
	}
)

func (db *DBClassify) Register() error {
	dbType, ok := viper.Get("DB.Type").(string)
	if !ok {
		log.Fatalf("DB.Type is not string")
	}
	user, ok := viper.Get("DB.User").(string)
	if !ok {
		log.Fatalf("DB.User is not string")
	}
	password, ok := viper.Get("DB.Password").(string)
	if !ok {
		log.Fatalf("DB.Password is not string")
	}
	address, ok := viper.Get("DB.Address").(string)
	if !ok {
		log.Fatalf("DB.Address is not string")
	}
	port, ok := viper.Get("DB.Port").(string)
	if !ok {
		log.Fatalf("DB.Port is not string")
	}
	dbName, ok := viper.Get("DB.DBName").(string)
	if !ok {
		log.Fatalf("DB.Port is not string")
	}

	if dbType == "mysql" {
		gormDB, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, password, address, port, dbName)), &gorm.Config{
			CreateBatchSize: 200,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			log.Fatalf("db err: %e", err)
		}
		db.GORM = gormDB
		return nil
	}

	log.Fatalf("%s not supported DB.Type", dbType)

	return nil
}
