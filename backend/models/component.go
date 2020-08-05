package models

type Component struct {
	Name   string  `gorm:"primary_key;auto_increment:false"`
	Tracks []Track `gorm:"many2many:tracks_components;"`
	Users  []User  `gorm:"mamy2many:users_instruments"`
}

type ComponentRepo interface {
	GetByName(name string) (*Component, error)
	Create(c *Component) error
	Update(c *Component) error
}
