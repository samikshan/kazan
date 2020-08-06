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
	if err := repo.db.Where(&models.Track{Model: gorm.Model{ID: id}}).Preload("Instruments").First(&t).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Create creates a new track entry in db
func (repo *TrackRepo) Create(t *models.Track, composer *models.User) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(t).Error; err != nil {
			return err
		}

		if err := tx.Model(composer).Association("Tracks").Append(t).Error; err != nil {
			return err
		}

		return nil
	})
}

// Update updates user entry in db
func (repo *TrackRepo) Update(t *models.Track) error {
	return repo.db.Update(t).Error
}

// UpdateInstruments updates instrument tags associated with a track
func (repo *TrackRepo) UpdateInstruments(t *models.Track) error {
	return repo.db.Model(t).Association("Instruments").Append(t.Instruments).Error
}

// GetTracksByInstrument retrieves all tracks that DO NOT have the listed instruments
func (repo *TrackRepo) GetTracksByInstrument(instrumentNames []string) ([]*models.Instrument, error) {
	var instruments []*models.Instrument
	if err := repo.db.Not("name", instrumentNames).Preload("Tracks").Find(&instruments).Error; err != nil {
		return nil, err
	}

	return instruments, nil
}
