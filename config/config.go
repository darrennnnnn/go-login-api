package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	Redis  RedisConfig
}

type JWTConfig struct {
	Secret []byte
}

type ServerConfig struct {
	ServerPort string
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		JWT: JWTConfig{
			Secret: []byte(getenv("JWT_SECRET")),
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
		Redis: RedisConfig{
			Host:     getenvDefault("REDIS_HOST", "localhost"),
			Port:     getenvDefault("REDIS_PORT", "8051"),
			Password: getenvDefault("REDIS_PASSWORD", "password"),
			DB:       getenvDefault("REDIS_DB", "0"),
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

func getenvDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
