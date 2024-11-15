// config.go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func InitDB() *gorm.DB {
    dsn := "host=localhost user=postgres password=postgres dbname=notes_app port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database:", err)
    }
    db.AutoMigrate(&User{}, &Note{})
    return db
}
