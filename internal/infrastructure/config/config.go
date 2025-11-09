package config

import (
    "fmt"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Security SecurityConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type ServerConfig struct {
    Port string
}

type SecurityConfig struct {
    SecretKey string
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf("failed to load .env file: %v", err)
    }

    cfg := &Config{
        Database: DatabaseConfig{
            Host:     os.Getenv("DATABASE_HOST"),
            Port:     os.Getenv("DATABASE_PORT"),
            User:     os.Getenv("DATABASE_USER"),
            Password: os.Getenv("DATABASE_PASSWORD"),
            DBName:   os.Getenv("DATABASE_NAME"),
            SSLMode:  os.Getenv("DATABASE_SSL_MODE"),
        },
        Server: ServerConfig{
            Port: os.Getenv("SERVER_PORT"),
        },
        Security: SecurityConfig{
            SecretKey: os.Getenv("SECRET_KEY"),
        },
    }

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("DATABASE_PASSWORD is required")
    }
    if cfg.Security.SecretKey == "" {
        return nil, fmt.Errorf("SECRET_KEY is required")
    }
    if _, err := strconv.Atoi(cfg.Server.Port); err != nil {
        return nil, fmt.Errorf("SERVER_PORT must be numeric: %v", err)
    }
    if _, err := strconv.Atoi(cfg.Database.Port); err != nil {
        return nil, fmt.Errorf("DATABASE_PORT must be numeric: %v", err)
    }

    return cfg, nil
}
