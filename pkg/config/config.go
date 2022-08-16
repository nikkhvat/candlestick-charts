package config

import (
	"log"

	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

type Config struct {
	Port             string `json:"port"`
	PostgresHost     string `json:"postgres_host"`
	PostgresUser     string `json:"postgres_user"`
	PostgresPassword string `json:"postgres_password"`
	PostgresDbname   string `json:"postgres_dbname"`
	PostgresPort     string `json:"postgres_port"`
	PostgresSslmode  string `json:"postgres_sslmode"`
	PostgresTimezone string `json:"postgres_timezone"`
	Key              string `json:"key"`
	Nominals         string `json:"nominals"`
}

var config Config

func initConfig() map[string]string {
	var err error
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Panic(err)
	}

	return envs
}

func GetConfig() Config {
	onceBody := func() {
		envs := initConfig()

		config = Config{
			Key:              envs["key"],
			Nominals:         envs["nominals"],
			Port:             envs["port"],
			PostgresHost:     envs["postgres_host"],
			PostgresUser:     envs["postgres_user"],
			PostgresPassword: envs["postgres_password"],
			PostgresDbname:   envs["postgres_dbname"],
			PostgresPort:     envs["postgres_port"],
			PostgresSslmode:  envs["postgres_sslmode"],
			PostgresTimezone: envs["postgres_timezone"],
		}
	}

	once.Do(onceBody)

	return config
}
