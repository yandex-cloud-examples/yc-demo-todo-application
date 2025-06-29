package config

import (
    "fmt"
    "log"
    "os"
)

type DbInstance struct {
    Host           string
    Port           string
    User           string
    Password       string
    Database       string
    SslMode        string
    TZ             string
}

type DbConfig struct {
    RW DbInstance
    RO DbInstance
}

func getEnv(key, fallback string) string {
    value, ok := os.LookupEnv(key)
    if !ok {
        return fallback
    }
    return value
}

func getRqEnv(key string) string {
    value, ok := os.LookupEnv(key)
    if !ok {
        log.Fatal(fmt.Sprintf("Environment variable %s is undefined\n", key))
    }
    return value
}

func NewDbConfig() *DbConfig {
    c := DbConfig{}
    c.RW.Host = getRqEnv("DB_HOST")
    c.RW.Port = getEnv("DB_PORT", "6432")
    c.RW.User = getRqEnv("DB_USER")
    c.RW.Password = getRqEnv("DB_PASSWORD")
    c.RW.Database = getRqEnv("DB_DATABASE")
    c.RW.SslMode = getEnv("DB_SSL_MODE", "prefer")
    c.RW.TZ = getEnv("DB_TZ", "UTC")
    c.RO.Host = getEnv("RO_DB_HOST", c.RW.Host)
    c.RO.Port = getEnv("RO_DB_PORT", c.RW.Port)
    c.RO.User = getEnv("RO_DB_USER", c.RW.User)
    c.RO.Password = getEnv("RO_DB_PASSWORD", c.RW.Password)
    c.RO.Database = getEnv("RO_DB_DATABASE", c.RW.Database)
    c.RO.SslMode = getEnv("RO_DB_SSL_MODE", c.RW.SslMode)
    c.RO.TZ = getEnv("RO_DB_TZ", c.RW.TZ)
    return &c
}
