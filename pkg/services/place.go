package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ratingBookingService/models"
)

// Getting data of place by placeId
func GetPlace(placeId int) (*models.Place, error) {
	var place *models.Place
	url := fmt.Sprintf("http://127.0.0.1:8000/v1/places/%v", placeId)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(response.Body).Decode(&place)
	if err != nil {
		return nil, err
	}
	return place, err
}
