package databases

import (
	"log"

	"ethical-be/app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(conf *config.Conf) (*gorm.DB, error) {
	var dsn string

	switch conf.App.Mode {
	case "staging":
		dsn = "host=" + conf.Db_staging.Host + " user=" + conf.Db_staging.User + " password=" + conf.Db_staging.Pass + " dbname=" + conf.Db_staging.Name + " port=" + conf.Db_staging.Port + " sslmode=disable TimeZone=Asia/Jakarta"
		log.Println(conf.App.Name, "runing on", conf.App.Mode, "mode")
	case "production":
		dsn = "host=" + conf.Db_prod.Host + " user=" + conf.Db_prod.User + " password=" + conf.Db_prod.Pass + " dbname=" + conf.Db_prod.Name + " port=" + conf.Db_prod.Port + " sslmode=disable TimeZone=Asia/Jakarta"
		log.Println(conf.App.Name, "runing on", conf.App.Mode, "mode")
	default:
		dsn = "host=" + conf.Db.Host + " user=" + conf.Db.User + " password=" + conf.Db.Pass + " dbname=" + conf.Db.Name + " port=" + conf.Db.Port + " sslmode=disable TimeZone=Asia/Jakarta"
		log.Println(conf.App.Name, "runing on", conf.App.Mode, "mode")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("database not connected")
	} else {
		log.Println("database connected")
	}
	return db, err

}
