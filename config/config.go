package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB DBConfig
	JWT JWTConfig
}

type JWTConfig struct {
	Secret string
}

type ServerConfig struct{
	ServerPort string
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		JWT: JWTConfig{
			Secret: getenv("JWT_SECRET"),
		},
		Server: ServerConfig{
			ServerPort: getenv("SERVER_PORT"),
		},
		DB: DBConfig{
			DBUser:     getenv("DB_USER"),
			DBPassword: getenv("DB_PASSWORD"),
			DBHost:     getenv("DB_HOST"),
			DBPort:     getenv("DB_PORT"),
			DBName:     getenv("DB_NAME"),
		},
	}
}

func getenv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        panic(key + " is required")
		}
    return v
}
