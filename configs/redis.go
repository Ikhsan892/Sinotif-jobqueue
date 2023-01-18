package configs

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	RedisHost     string
	RedisPort     uint
	RedisUsername string
	RedisPassword string
	RedisDatabase int
}

func LoadRedisConfig() RedisConfig {
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		panic("Redis Config Cannot convert err")
	}

	redisDb, errDb := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if errDb != nil {
		panic("RedisConfig Cannot convert errDB")
	}

	return RedisConfig{
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     uint(redisPort),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDatabase: redisDb,
	}
}
