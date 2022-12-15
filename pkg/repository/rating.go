package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"ratingBookingService/models"
	"ratingBookingService/pkg/services"
	"strings"
)

type RatingPostgres struct {
	db *sql.DB
}

func CreateNewRatingRepository(db *sql.DB) *RatingPostgres {
	return &RatingPostgres{db: db}
}

func (r *RatingPostgres) Create(c *models.CreateRatingPlace) (*models.RatingPlace, error) {
	var rating *models.RatingPlace = new(models.RatingPlace)
	_, err := services.GetPlace(c.PlaceId)
	if err != nil {
		return nil, errors.New("place_id doesn't exists")
	}
	sqlStatement := `INSERT INTO rating_place (place_id, review, rating, user_id) VALUES ($1, $2, $3, $4) RETURNING *`
	row := r.db.QueryRow(sqlStatement, c.PlaceId, c.Review, c.Rating, c.User)
	err = row.Scan(&rating.Id, &rating.PlaceId, &rating.Rating, &rating.Review, &rating.CreatedAt, &rating.User)
	if err != nil {
		return nil, errors.New("Unable to execute the query")
	}
	return rating, nil
}

func (r *RatingPostgres) GetAll() ([]*models.RatingPlace, error) {
	var ratings []*models.RatingPlace
	query := fmt.Sprintf("SELECT * FROM rating_place")
	rows, err := r.db.Query(query)
	if err != nil {
		log.Fatal("[ratingsAll] Error db query ", err)
	}
	for rows.Next() {
		rate := &models.RatingPlace{}
		err = rows.Scan(&rate.Id, &rate.PlaceId, &rate.Rating, &rate.Review, &rate.CreatedAt, &rate.User)
		ratings = append(ratings, rate)
	}
	return ratings, err
}

func (r *RatingPostgres) GetById(ratingId int) (*models.RatingPlace, error) {
	var rating *models.RatingPlace
	row := r.db.QueryRow("SELECT * FROM rating_place WHERE id = $1", ratingId)
	err := row.Scan(&rating.Id, &rating.PlaceId, &rating.Rating, &rating.Review, &rating.CreatedAt, &rating.User)
	return rating, err
}

func (r *RatingPostgres) Update(ratingId int, c *models.CreateRatingPlace) (*models.RatingPlace, error) {
	var rating *models.RatingPlace = new(models.RatingPlace)
	setValues := make([]string, 0)

	if c.Rating != 0 {
		setValues = append(setValues, fmt.Sprintf("rating=%v", c.Rating))
	}
	if c.PlaceId != 0 {
		setValues = append(setValues, fmt.Sprintf("place_id=%v", c.PlaceId))
	}
	if c.Review != "" {
		setValues = append(setValues, fmt.Sprintf(`review='%s'`, c.Review))
	}
	if c.User != 0 {
		setValues = append(setValues, fmt.Sprintf(`user_id=%v`, c.User))
	}
	setQuery := strings.Join(setValues, ", ")

	sqlStatement := fmt.Sprintf("UPDATE rating_place SET %s WHERE id=$1 RETURNING *", setQuery)
	row := r.db.QueryRow(sqlStatement, ratingId)
	err := row.Scan(&rating.Id, &rating.PlaceId, &rating.Rating, &rating.Review, &rating.CreatedAt)
	if err != nil {
		return nil, errors.New("Unable to execute the query")
	}
	return rating, nil
}

func (r *RatingPostgres) DeleteByID(ratingId int) error {
	_, err := r.db.Exec("DELETE FROM rating_place WHERE id=$1", ratingId)
	if err != nil {
		return err
	}
	return nil
}

func (r *RatingPostgres) GetRatingsByPlaceId(placeId int) ([]*models.RatingPlace, error) {
	var ratings []*models.RatingPlace
	rows, err := r.db.Query("SELECT * FROM rating_place WHERE place_id = $1", placeId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rate := &models.RatingPlace{}
		err = rows.Scan(&rate.Id, &rate.PlaceId, &rate.Rating, &rate.Review, &rate.CreatedAt, &rate.User)
		ratings = append(ratings, rate)
	}
	return ratings, err
}
