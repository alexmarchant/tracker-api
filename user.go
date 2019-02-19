package main

import (
  "time"
)

type User struct {
  ID uint `gorm:"primary_key" json:"id"`
  CreatedAt time.Time `json:"createdAt"`
  UpdatedAt time.Time `json:"updatedAt"`
  DeletedAt *time.Time `json:"-"`
  Email string `gorm:"UNIQUE_INDEX;UNIQUE;NOT NULL" json:"email"`
  PasswordHash string `gorm:"NOT NULL" json:"-"`
}
