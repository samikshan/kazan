package models

type Instrument struct {
	Name   string  `gorm:"primary_key;auto_increment:false"`
	Tracks []Track `gorm:"many2many:tracks_instruments"`
	Users  []User  `gorm:"many2many:users_instruments"`
}

type InstrumentRepo interface {
	GetByName(name string) (*Instrument, error)
	Create(i *Instrument) error
	Update(i *Instrument) error
}
