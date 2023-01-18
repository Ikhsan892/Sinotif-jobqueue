package configs

import (
	"github.com/joho/godotenv"
	"os"
	"sinotif/pkg/utils"
)

type Config struct {
	DB     DBConfig
	Redis  RedisConfig
	Report ReportConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		utils.Error(utils.DATABASE, "Error Loading configuration ", err)
		os.Exit(1)
	}

	return &Config{
		DB:     LoadDBConfig(),
		Redis:  LoadRedisConfig(),
		Report: LoadReportConfig(),
	}
}
