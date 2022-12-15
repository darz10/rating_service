package handlers

import "ratingBookingService/pkg/repository"

func CheckOwnerRating(userId float64, currRatingId int, r *repository.RatingPostgres) bool {
	rating, err := r.GetById(currRatingId)
	if err != nil {
		return false
	}
	if rating.User == int(userId) {
		return true
	}
	return false
}
