package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"unique;not null"`
	WalletAddress string `gorm:"not null"`
}

type UserRepo interface {
	Create(u *User) error
	Update(u *User) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

// func (u *User) HashPassword(plain string) error {
// 	// Generate "hash" to store from user password
// 	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.WithError(err).Error("failed to generate password hash")
// 		return err
// 	}

// 	u.Password = string(hash)
// 	return nil
// }

// func (u *User) CheckPassword(plain string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
// 	return err == nil
// }

// func (u *User) FFSCreate() (string, string, error) {
// 	ma, err := multiaddr.NewMultiaddr(backend.Cfg.PowGrpcHost)
// 	if err != nil {
// 		log.WithError(err).Error("error parsing multiaddress")
// 		return "", "", err
// 	}
// 	client, err := pow.NewClient(ma, grpc.WithInsecure())
// 	if err != nil {
// 		log.WithError(err).Error("failed to create powergate client")
// 		return "", "", err
// 	}

// 	ffsID, ffsToken, err := client.FFS.Create(context.Background())
// 	if err != nil {
// 		log.WithError(err).Error("failed to create powergate ffs instance")
// 		return "", "", err
// 	}

// 	err = client.Close()
// 	if err != nil {
// 		log.WithError(err).Error("failed to close powergate client")
// 	}

// 	return ffsID, ffsToken, nil
// }
