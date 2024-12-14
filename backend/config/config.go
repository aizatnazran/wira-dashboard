package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    JWTSecret  string
    RedisHost  string
    RedisPort  string
    RedisPassword string
    ServerPort    string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
    }

    config := &Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        DBSSLMode:  os.Getenv("DB_SSLMODE"),
        JWTSecret:  os.Getenv("JWT_SECRET"),
        RedisHost:  os.Getenv("REDIS_HOST"),
        RedisPort:  os.Getenv("REDIS_PORT"),
        RedisPassword: os.Getenv("REDIS_PASSWORD"),
        ServerPort:    os.Getenv("SERVER_PORT"),
    }

    // Set default server port if not specified
    if config.ServerPort == "" {
        config.ServerPort = "8080"
    }

    return config, nil
}

func (c *Config) GetDBConnString() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.DBUser,
        c.DBPassword,
        c.DBHost,
        c.DBPort,
        c.DBName,
        c.DBSSLMode,
    )
}
