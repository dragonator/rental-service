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
	IDs      []int32
	PriceMin *int64
	PriceMax *int64
	Near     *Location
	OrderBy  *string
}

// RentalSortFields defines allowed fields for sorting.
var RentalSortFields = map[string]*struct{}{
	"ids":       nil,
	"price_min": nil,
	"price_max": nil,
	"near":      nil,
	"sort":      nil,
	"limit":     nil,
	"offset":    nil,
}
