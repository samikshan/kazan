package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/samikshan/kazan/backend/models"
	"github.com/sirupsen/logrus"
)

type InstrumentRepo struct {
	db *gorm.DB
}

func NewInstrumentRepo(db *gorm.DB) *InstrumentRepo {
	return &InstrumentRepo{
		db: db,
	}
}

// GetByName gets instrument by name
func (repo *InstrumentRepo) GetByName(name string) (*models.Instrument, error) {
	var i models.Instrument
	if err := repo.db.Where(&models.Instrument{Name: name}).Take(&i).Error; err != nil {
		logrus.Error(err)
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

// Create creates a new instrument entry in db
func (repo *InstrumentRepo) Create(i *models.Instrument) error {
	return repo.db.Create(i).Error
}

// Update updates instrument entry in db
func (repo *InstrumentRepo) Update(i *models.Instrument) error {
	return repo.db.Update(i).Error
}
