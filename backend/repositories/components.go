package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/samikshan/kazan/backend/models"
	"github.com/sirupsen/logrus"
)

type ComponentRepo struct {
	db *gorm.DB
}

func NewComponentRepo(db *gorm.DB) *ComponentRepo {
	return &ComponentRepo{
		db: db,
	}
}

// GetByName gets component by name
func (repo *ComponentRepo) GetByName(name string) (*models.Component, error) {
	var c models.Component
	if err := repo.db.Where(&models.Component{Name: name}).Take(&c).Error; err != nil {
		logrus.Error(err)
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// Create creates a new component entry in db
func (repo *ComponentRepo) Create(c *models.Component) error {
	return repo.db.Create(c).Error
}

// Update updates component entry in db
func (repo *ComponentRepo) Update(c *models.Component) error {
	return repo.db.Update(c).Error
}
