package models

import "time"

// Rental - a model for the rental entity.
type Rental struct {
	ID              int32
	UserID          int32
	Name            string
	Type            string
	Description     string
	Sleeps          int32
	PricePerDay     int64
	HomeCity        string
	HomeState       string
	HomeZip         string
	HomeCountry     string
	VehicleMake     string
	VehicleModel    string
	VehicleYear     int32
	VehicleLength   float32
	Latitude        float32
	Longitude       float32
	PrimaryImageURL string
	Created         time.Time
	Updated         time.Time

	User *User
}

// Rentals is a slice of Rental objects.
type Rentals []*Rental
