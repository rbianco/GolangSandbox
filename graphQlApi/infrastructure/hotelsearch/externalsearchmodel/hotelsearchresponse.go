package externalsearchmodel

type HotelSearchResponse struct {
	Token       string  `json:"token"`
	Hotels      []Hotel `json:"hotels"`
	TotalHotels int     `json:"totalHotels"`
}
