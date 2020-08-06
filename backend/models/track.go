package models

import "github.com/jinzhu/gorm"

type Track struct {
	gorm.Model
	CID           string `gorm:"notnull"`
	Title         string
	ComposerID    uint
	Composer      User
	ParentTrackID uint
	ParentTrack   *Track
	nForks        uint         `gorm:"default:0"`
	Forks         []Track      `gorm:"foreignkey:ParentTrackID"`
	Instruments   []Instrument `gorm:"many2many:tracks_instruments"`
}

func (t *Track) AfterCreate(tx *gorm.DB) error {
	return tx.Model(t).Association("Instruments").Append(t.Instruments).Error
}

type TrackRepo interface {
	GetByTrackID(id uint) (*Track, error)
	Create(t *Track, composer *User) error
	Update(t *Track) error
	GetTracksByInstrument(instruments []string) ([]*Instrument, error)
}
