package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/samikshan/kazan/backend/models"
)

type TrackRepo struct {
	db *gorm.DB
}

func NewTrackRepo(db *gorm.DB) *TrackRepo {
	return &TrackRepo{
		db: db,
	}
}

// GetByTrackID gets user by track id
func (repo *TrackRepo) GetByTrackID(id uint) (*models.Track, error) {
	var t models.Track
	if err := repo.db.Where(&models.Track{Model: gorm.Model{ID: id}}).First(&t).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Create creates a new track entry in db
func (repo *TrackRepo) Create(t *models.Track) error {
	return repo.db.Create(t).Error
}

// Update updates user entry in db
func (repo *TrackRepo) Update(t *models.Track) error {
	return repo.db.Update(t).Error
}
