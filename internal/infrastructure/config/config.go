package config

import (
    "fmt"
    "log"
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
        log.Println("No .env file found, using system environment variables")
    }

    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "payment_gateway_manjo"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
        Security: SecurityConfig{
            SecretKey: getEnv("SECRET_KEY", ""),
        },
    }

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("DB_PASSWORD is required")
    }
    if cfg.Security.SecretKey == "" {
        return nil, fmt.Errorf("SECRET_KEY is required")
    }
    if _, err := strconv.Atoi(cfg.Server.Port); err != nil {
        return nil, fmt.Errorf("SERVER_PORT must be numeric: %v", err)
    }
    if _, err := strconv.Atoi(cfg.Database.Port); err != nil {
        return nil, fmt.Errorf("DB_PORT must be numeric: %v", err)
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
