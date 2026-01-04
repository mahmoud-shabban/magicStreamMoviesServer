package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/database"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	moviesCollection = database.OpenCollection("movies")
	validate         = validator.New()
)

func GetMovies() gin.HandlerFunc {

	return func(c *gin.Context) {
		var movies []models.Movie

		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		// moviesCollection := database.OpenCollection("movies")

		if moviesCollection == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "faild to retrieve movies"})
			return
		}

		cursor, err := moviesCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "faild to retrieve movies"})
			return
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "faild to retrieve movies"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"movies": movies})

	}
}

func GetMovieByID() gin.HandlerFunc {

	return func(c *gin.Context) {
		var movie models.Movie

		id := c.Param("id")

		if len(id) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"movie": "id should be vaild imdb_id"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		// moviesCollection := database.OpenCollection("movies")

		if moviesCollection == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("faild to retrieve movie with id %s", id)})
			return
		}

		result := moviesCollection.FindOne(ctx, bson.M{"imdb_id": id})
		if result.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("faild to get movie with id %s: %s", id, result.Err().Error())})
			return
		}

		err := result.Decode(&movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("faild to retrieve movie with id %s: %s", id, err.Error())})
			return
		}

		c.JSON(http.StatusOK, gin.H{"movie": movie})
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		// bodyBytes, err := io.ReadAll(c.Request.Body)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// 	return
		// }
		var tmpMovie struct {
			ImdbID      string         `bson:"imdb_id" json:"imdb_id" validate:"required"`
			Title       string         `bson:"title" json:"title" validate:"required,min=2,max=500"`
			PosterPath  string         `bson:"poster_path" json:"poster_path" validate:"required,url"`
			YoutubeID   string         `bson:"youtube_id" json:"youtube_id" validate:"required"`
			Genre       []models.Genre `bson:"genre" json:"genre" validate:"required,dive"`
			AdminReview string         `bson:"admin_review" json:"admin_review" validate:"required"`
			Ranking     models.Ranking `bson:"ranking" json:"ranking" validate:"required"`
		}

		// err = json.Unmarshal(bodyBytes, &tmpMovie)
		err := c.ShouldBindJSON(&tmpMovie)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err = validate.Struct(tmpMovie)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": "validation failed", "details": err})
			return
		}

		movie := models.Movie{
			ImdbID:      tmpMovie.ImdbID,
			Title:       tmpMovie.Title,
			PosterPath:  tmpMovie.PosterPath,
			YoutubeID:   tmpMovie.YoutubeID,
			Genre:       tmpMovie.Genre,
			AdminReview: tmpMovie.AdminReview,
			Ranking:     tmpMovie.Ranking,
		}

		// collection := database.OpenCollection("movies")
		result, err := moviesCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": "validation failed", "details": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"movie": result})

	}
}
