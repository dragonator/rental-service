package contract

type GetRentalByIDRequest struct{}

type User struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Price struct {
	Day int64 `json:"day"`
}

type Location struct {
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	Country   string  `json:"country"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type GetRentalByIDResponse struct {
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
}
