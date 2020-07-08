package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	pow "github.com/textileio/powergate/api/client"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

	"github.com/samikshan/kazan/backend"
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
	Update(u *User) error
	GetByID(id uint) (*User, error)
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

func (u *User) FFSCreate() (string, string, error) {
	ma, err := multiaddr.NewMultiaddr(backend.Cfg.PowGrpcHost)
	if err != nil {
		log.WithError(err).Error("error parsing multiaddress")
		return "", "", err
	}
	client, err := pow.NewClient(ma, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Error("failed to create powergate client")
		return "", "", err
	}

	ffsID, ffsToken, err := client.FFS.Create(context.Background())
	if err != nil {
		log.WithError(err).Error("failed to create powergate ffs instance")
		return "", "", err
	}

	err = client.Close()
	if err != nil {
		log.WithError(err).Error("failed to close powergate client")
	}

	return ffsID, ffsToken, nil
}
