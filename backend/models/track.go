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
	Forks         []Track     `gorm:"foreignkey:ParentTrackID"`
	Components    []Component `gorm:"mamy2many:tracks_components"`
}

type TrackRepo interface {
	GetByTrackID(id uint) (*Track, error)
	Create(t *Track) error
	Update(t *Track) error
}
