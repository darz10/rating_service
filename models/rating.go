package models

import (
	"errors"
	"time"
)

type CreateRatingPlace struct {
	PlaceId int     `json:"place_id" binding:"required"`
	Rating  int     `json:"rating" binding:"required"`
	Review  string  `json:"review"`
	User    float64 `json:"user_id" binding:"required"`
}

func (c *CreateRatingPlace) ValidateRating() error {
	var InvalidValue = errors.New("Invalid rating value, rating should be >=1 and <=5 ")
	if c.Rating < 1 || c.Rating > 5 {
		return InvalidValue
	}
	return nil
}

type RatingPlace struct {
	Id        int       `json:"id" binding:"required"`
	PlaceId   int       `json:"place_id" binding:"required"`
	Rating    int       `json:"rating" binding:"required"`
	Review    string    `json:"review"`
	User      int       `json:"user_id" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}
