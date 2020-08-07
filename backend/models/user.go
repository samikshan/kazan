package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username      string       `gorm:"unique;not null"`
	WalletAddress string       `gorm:"unique;not null"`
	Tracks        []Track      `gorm:"foreignkey:ComposerID"`
	Instruments   []Instrument `gorm:"many2many:users_instruments"`
	DisplayName   string
}

type UserRepo interface {
	Create(u *User) error
	Update(u *User) error
	UpdateInstruments(u *User) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByWalletAddr(addr string) (*User, error)
}
