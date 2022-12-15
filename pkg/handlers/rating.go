package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ratingBookingService/models"
	"ratingBookingService/pkg/repository"
	"ratingBookingService/pkg/services"
	"strconv"
)

type RatingHandler struct {
	repository *repository.RatingPostgres
}

func CreateNewRatingHandler(r *repository.RatingPostgres) *RatingHandler {
	return &RatingHandler{repository: r}
}

// @BasePath /api

// @Summary Get list of ratings
// @Tags ratings
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param        place    query     int  false  "getting ratings for place"
// @Success 200 {object} []models.RatingPlace
// @Router /ratings [get]
func (h *RatingHandler) GetRatingsAll(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	place, exists := queryParams["place"]
	if exists && len(place) == 1 {
		placeId, err := strconv.Atoi(place[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ratings, err := h.repository.GetRatingsByPlaceId(placeId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, ratings)
	} else {
		ratings, err := h.repository.GetAll()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unable to execute the query"})
			return
		}
		c.IndentedJSON(http.StatusOK, ratings)
	}
}

// @BasePath /api

// @Summary Get rating by id
// @Tags ratings
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Success 200 {object} models.RatingPlace
// @Router /ratings/{id} [get]
func (h *RatingHandler) GetRatingByID(c *gin.Context) {
	ratingId := c.Param("id")
	id, err := strconv.Atoi(ratingId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	rating, err := h.repository.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, rating)
}

// @BasePath /api

// @Summary Create rating
// @Tags ratings
// @Schemes
// @Description
// @Param input body models.CreateRatingPlace true "Create rating"
// @Accept json
// @Produce json
// @Success 200 {object} models.RatingPlace
// @Router /ratings [post]
func (h *RatingHandler) CreateRating(c *gin.Context) {
	user, exists := c.Keys["User"]
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	var createRating *models.CreateRatingPlace
	err := c.ShouldBindJSON(&createRating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createRating.User = user.(float64)
	err = createRating.ValidateRating()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = services.GetPlace(createRating.PlaceId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	rating, err := h.repository.Create(createRating)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to execute the query"})
		return
	}
	c.IndentedJSON(http.StatusCreated, rating)
}

// @BasePath /api

// @Summary Update rating
// @Tags ratings
// @Schemes
// @Description
// @Param input body models.CreateRatingPlace false "Update rating"
// @Accept json
// @Produce json
// @Success 200 {object} models.RatingPlace
// @Router /ratings/{id} [patch]
func (h *RatingHandler) UpdateRating(c *gin.Context) {
	user, exists := c.Keys["User"]
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	var updateRating *models.CreateRatingPlace
	ratingId := c.Param("id")
	id, err := strconv.Atoi(ratingId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	err = c.ShouldBindJSON(&updateRating)
	updateRating.User = user.(float64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	isOwner := CheckOwnerRating(updateRating.User, id, h.repository)
	if !isOwner {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only owner can influence the rating"})
		return
	}
	err = updateRating.ValidateRating()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	rating, err := h.repository.Update(id, updateRating)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to execute the query"})
		return
	}
	c.IndentedJSON(http.StatusOK, rating)
}

// @BasePath /api

// @Summary Update rating
// @Tags ratings
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Success 204 {object} models.RatingPlace
// @Router /ratings/{id} [delete]
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	user, exists := c.Keys["User"]
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	ratingId := c.Param("id")
	id, err := strconv.Atoi(ratingId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	isOwner := CheckOwnerRating(user.(float64), id, h.repository)
	if !isOwner {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only owner can influence the rating"})
		return
	}
	err = h.repository.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to execute the query"})
		return
	}
	c.Writer.WriteHeader(204)
}
