package storage

// Pagination specifies a pagination for the request.
type Pagination struct {
	Limit  *int
	Offset *int
}

// Location represents a location by latitude and longitude values.
type Location struct {
	Latitude  float32
	Longitude float32
}

// RentalFilters is a filters type to be used for listing rentals.
type RentalFilters struct {
	Pagination
	IDs      []string
	PriceMin *int64
	PriceMax *int64
	Near     *Location
	OrderBy  *string
}
