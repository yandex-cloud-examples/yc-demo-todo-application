package database
import (
    "fmt"
    "log"
    "time"
    "todo/config"
    "gorm.io/driver/postgres"
    "gorm.io/plugin/dbresolver"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func ConnectRW(c *config.DbInstance) *gorm.DB {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", 
                c.Host, c.User, c.Password, c.Database, c.Port, c.SslMode, c.TZ)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        PrepareStmt: false,
        Logger: logger.Default.LogMode(logger.Info),
        SkipDefaultTransaction: true,
        DisableAutomaticPing: true,
    })
    if err != nil {
        log.Fatal(fmt.Sprintf("Failed to connect to database: %s\n", err))
    }
    log.Print(fmt.Sprintf("Connected to %s\n", c.Host))
    sqlDB, err := db.DB()
    sqlDB.SetMaxIdleConns(50)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Minute * 3)
    return db
}

func ConnectRO(c *config.DbInstance, db *gorm.DB) *gorm.DB {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", 
                c.Host, c.User, c.Password, c.Database, c.Port, c.SslMode, c.TZ)
    db.Use(dbresolver.Register(
        dbresolver.Config{
            Replicas: []gorm.Dialector{postgres.Open(dsn),},
            Policy: dbresolver.RandomPolicy{},
            TraceResolverMode: true,
    }))
    log.Print(fmt.Sprintf("Read-only replica %s\n", c.Host))
    return db
}
