package consumer

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sinotif/configs"
	"sinotif/pkg/utils"
)

func DBError() {
	if msg := recover(); msg != nil {
		utils.Error(utils.DATABASE, msg)
	} else {
		utils.Info(utils.DATABASE, "Database Connected successfully")
	}
}

func ConnectionDB(cfg configs.DBConfig) *gorm.DB {
	defer DBError()

	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port)

	utils.Debug(utils.DATABASE, dataSourceName)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		panic(err.Error())
	}

	return db

}
