// Packager hotelsearch contains functions for working with hotel search services
package hotelsearch

import (
	"encoding/json"
	"fmt"

	"github.com/rbianco/GolangSandbox/graphQlApi/domain/model"
	"github.com/rbianco/GolangSandbox/graphQlApi/infrastructure/hotelsearch/externalsearchmodel"
	resty "gopkg.in/resty.v1"
)

// Retrieve hotel from hotel search service
func Hotel(searchEndpoing string, s string) []model.Hotel {

	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"organizationId": "1",
			"pageSize":       "20",
			"currentPage":    "1",
			"culture":        "es-mx",
			"checkIn":        "2019-01-17",
			"getFilters":     "true",
		}).
		SetHeader("Accept", "application/json").
		Get(searchEndpoing + s)

	if resp.IsError() {
		return []model.Hotel{}
	}

	var hotels externalsearchmodel.HotelSearchResponse
	var responseBody = resp.Body()
	json.Unmarshal(responseBody, &hotels)

	fmt.Println("response", resp.IsSuccess())
	fmt.Println("error", err)

	var result []model.Hotel
	for _, hotel := range hotels.Hotels {
		result = append(result, model.Hotel{
			ID:    hotel.ID,
			Name:  hotel.Name,
			Title: hotel.Title,
			URI:   hotel.URI,
		})
	}

	return result
}
