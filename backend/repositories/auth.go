package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/samikshan/kazan/backend/models"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (repo *AuthRepo) Create(auth *models.Auth) error {
	return repo.db.Create(auth).Error
}

// func (repo *UserRepo) Update(u *models.User) error {
// 	return repo.db.Save(u).Error
// }

func (repo *AuthRepo) GetByLookupKey(key string) (*models.Auth, error) {
	var auth models.Auth
	if err := repo.db.Where(&models.Auth{LookupKey: key}).First(&auth).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &auth, nil
}
