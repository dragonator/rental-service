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
var RentalSortFields = []string{
	"id",
	"name",
	"type",
	"make",
	"model",
	"year",
	"length",
	"sleeps",
	"price_per_day",
}

// SortFieldAllowed checks whether the given field in allowed to sort rentals by.
func SortFieldAllowed(field string) bool {
	for _, f := range RentalSortFields {
		if field == f {
			return true
		}
	}

	return false
}
