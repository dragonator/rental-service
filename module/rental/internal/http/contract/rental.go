package contract

// User is a contract for the user object.
type User struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Price is a contract for the price object.
type Price struct {
	Day int64 `json:"day"`
}

// Location is a contract for the location object.
type Location struct {
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	Country   string  `json:"country"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

// Rental is a contract for the rental object.
type Rental struct {
	ID              int32   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	Make            string  `json:"make"`
	Model           string  `json:"model"`
	Year            int32   `json:"year"`
	Length          float32 `json:"length"`
	Sleeps          int32   `json:"sleeps"`
	PrimaryImageURL string  `json:"primary_image_url"`
	Price           Price
	Location        Location
	User            User
}

// ListRentalsQuery is used to decode the query parameters of ListRentals.
type ListRentalsQuery struct {
	Ids      []int32   `schema:"ids"`
	PriceMin *int64    `schema:"price_min"`
	PriceMax *int64    `schema:"price_max"`
	Near     []float32 `schema:"near"`
	Limit    *int      `schema:"limit"`
	Offset   *int      `schema:"offset"`
	Sort     *string   `schema:"sort"`
}

// GetRentalByIDResponse is a server response getting a single rental by id.
type GetRentalByIDResponse struct {
	Rental
}

// ListRentalsResponse is a server response listing rentals by filters.
type ListRentalsResponse []*Rental
