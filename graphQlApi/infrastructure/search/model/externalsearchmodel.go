package externalsearchmodel

type Hotel struct {
	ID    int    `json:"hotelId"`
	Name  string `json:"name"`
	Title string `json:"title"`
	URI   string `json:"uri"`
}

type HotelSearchResponse struct {
	Token       string  `json:"token"`
	Hotels      []Hotel `json:"hotels"`
	TotalHotels int     `json:"totalHotels"`
}
