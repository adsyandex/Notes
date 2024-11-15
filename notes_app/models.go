// models.go
package main

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
    Notes    []Note
}

type Note struct {
    gorm.Model
    Title     string
    Content   string
    UserID    uint
    ExpiresAt *time.Time // опциональное время истечения
}
