package models

import (
	"github.com/jinzhu/gorm"
)

type Auth struct {
	gorm.Model
	LookupKey  string `gorm:"unique_index;not null"`
	IV         string `gorm:"not null"`
	CipherText string `gorm:"not null"`
}

type AuthRepo interface {
	Create(auth *Auth) error
	// Update(auth *Auth) error
	GetByLookupKey(key string) (*Auth, error)
}
