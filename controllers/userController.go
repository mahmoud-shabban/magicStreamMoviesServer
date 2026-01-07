package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/database"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/models"
	"github.com/mahmoud-shabban/magicStreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UsersCollection mongo.Collection = *database.OpenCollection("users")

func HashPassowrd(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
			return
		}

		validator := validator.New()
		if err := validator.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashadPassowrd, err := HashPassowrd(user.Password)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.Password = hashadPassowrd
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		count, err := UsersCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal srever error"})
			return
		}

		if count != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		resutlt, err := UsersCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": resutlt})
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
			return
		}

		validator := validator.New()
		if err := validator.Struct(userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var foundUser models.User

		filter := bson.M{"email": userLogin.Email}

		err := UsersCollection.FindOne(ctx, filter).Decode(&foundUser)

		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found, please signup instead"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "worng email or password"})
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(
			foundUser.Email,
			foundUser.FirstName,
			foundUser.LastName,
			foundUser.Role,
			foundUser.UserID,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		if err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			UserId:          foundUser.UserID,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			Token:           token,
			RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})

	}
}
