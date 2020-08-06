package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/samikshan/kazan/backend/models"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (repo *UserRepo) Create(u *models.User) error {
	return repo.db.Create(u).Error
}

func (repo *UserRepo) Update(u *models.User) error {
	// return repo.db.Model()
	return repo.db.Save(u).Error
}

func (repo *UserRepo) UpdateInstruments(u *models.User) error {
	return repo.db.Model(u).Association("Instruments").Append(u.Instruments).Error
}

func (repo *UserRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := repo.db.Where(&models.User{Username: username}).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}

// GetByUserID gets user by user id
func (repo *UserRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := repo.db.Where(&models.User{Model: gorm.Model{ID: id}}).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByWalletAddr gets user by user id
func (repo *UserRepo) GetByWalletAddr(addr string) (*models.User, error) {
	var u models.User
	if err := repo.db.Preload("Instruments").Where(&models.User{WalletAddress: addr}).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
