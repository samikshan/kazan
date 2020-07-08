package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique_index;not null"`
	Name     string `gorm:"default:'creator'"`
	Password string `gorm:"not null"`
	FFSToken string
	Token    string
}

type UserRepo interface {
	Create(u *User) error
	GetByUserID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

func (u *User) HashPassword(plain string) error {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("failed to generate password hash")
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
