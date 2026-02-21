package handlers

import (
	"context"
	"fmt"

	"movie-api/internal/auth"
	"movie-api/internal/database"
	"movie-api/internal/tmdb"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const dbName = "movieapi"

// ================= AUTH =================

func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.BindJSON(&input)

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	result, _ := database.Client.
		Database(dbName).
		Collection("users").
		InsertOne(context.TODO(), bson.M{
			"email":    input.Email,
			"password": string(hash),
		})

	c.JSON(200, gin.H{"user_id": result.InsertedID})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.BindJSON(&input)

	users := database.Client.Database(dbName).Collection("users")

	var user bson.M
	users.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)

	err := bcrypt.CompareHashAndPassword(
		[]byte(user["password"].(string)),
		[]byte(input.Password),
	)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	objectID := user["_id"].(primitive.ObjectID)
	token, _ := auth.GenerateToken(objectID.Hex())

	c.JSON(200, gin.H{"token": token})
}

// ================= SEARCH =================

func Search(c *gin.Context) {
	query := c.Query("q")
	apiKey := c.MustGet("tmdbKey").(string)

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s",
		apiKey,
		query,
	)

	data, _ := tmdb.Fetch(url)
	c.Data(200, "application/json", data)
}

// ================= RECOMMEND =================

func Recommend(c *gin.Context) {
	userID := c.GetString("user_id")
	apiKey := c.MustGet("tmdbKey").(string)

	ratings := database.Client.Database(dbName).Collection("ratings")

	cursor, _ := ratings.Find(context.TODO(), bson.M{
		"user_id":     userID,
		"user_rating": bson.M{"$gte": 4},
	})

	var rated []bson.M
	cursor.All(context.TODO(), &rated)

	if len(rated) == 0 {
		c.JSON(200, gin.H{"message": "Rate some movies first"})
		return
	}

	first := rated[0]
	genres := first["genres"].(primitive.A)
	genre := genres[0].(bson.M)
	genreID := fmt.Sprintf("%v", genre["id"])

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/discover/movie?api_key=%s&with_genres=%s",
		apiKey,
		genreID,
	)

	data, _ := tmdb.Fetch(url)
	c.Data(200, "application/json", data)
}

// ================= WATCHLIST =================

func AddWatchlist(c *gin.Context) {
	userID := c.GetString("user_id")
	movieID := c.Param("id")
	apiKey := c.MustGet("tmdbKey").(string)

	watchlists := database.Client.Database(dbName).Collection("watchlists")

	// Prevent duplicates
	count, _ := watchlists.CountDocuments(context.TODO(), bson.M{
		"user_id":  userID,
		"movie_id": movieID,
	})

	if count > 0 {
		c.JSON(400, gin.H{"message": "Already in watchlist"})
		return
	}

	// Fetch movie details
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%s?api_key=%s",
		movieID,
		apiKey,
	)

	data, _ := tmdb.Fetch(url)

	var movieData bson.M
	bson.UnmarshalExtJSON(data, true, &movieData)

	watchlists.InsertOne(context.TODO(), bson.M{
		"user_id":      userID,
		"movie_id":     movieID,
		"title":        movieData["title"],
		"release_date": movieData["release_date"],
		"genres":       movieData["genres"],
		"poster_path":  movieData["poster_path"],
		"vote_average": movieData["vote_average"],
	})

	c.JSON(200, gin.H{"message": "Added to watchlist"})
}

func GetWatchlist(c *gin.Context) {
	userID := c.GetString("user_id")

	cursor, _ := database.Client.
		Database(dbName).
		Collection("watchlists").
		Find(context.TODO(), bson.M{"user_id": userID})

	var results []bson.M
	cursor.All(context.TODO(), &results)

	c.JSON(200, results)
}

// ================= RATE =================

func RateMovie(c *gin.Context) {
	userID := c.GetString("user_id")
	movieID := c.Param("id")
	apiKey := c.MustGet("tmdbKey").(string)

	var input struct {
		Rating int `json:"rating"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(400, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	ratings := database.Client.Database(dbName).Collection("ratings")

	// Prevent duplicate rating
	count, _ := ratings.CountDocuments(context.TODO(), bson.M{
		"user_id":  userID,
		"movie_id": movieID,
	})

	if count > 0 {
		c.JSON(400, gin.H{"message": "Already rated"})
		return
	}

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%s?api_key=%s",
		movieID,
		apiKey,
	)

	data, _ := tmdb.Fetch(url)

	var movieData bson.M
	bson.UnmarshalExtJSON(data, true, &movieData)

	ratings.InsertOne(context.TODO(), bson.M{
		"user_id":      userID,
		"movie_id":     movieID,
		"title":        movieData["title"],
		"release_date": movieData["release_date"],
		"genres":       movieData["genres"],
		"user_rating":  input.Rating,
	})

	c.JSON(200, gin.H{"message": "Rated successfully"})
}

func GetRated(c *gin.Context) {
	userID := c.GetString("user_id")

	cursor, _ := database.Client.
		Database(dbName).
		Collection("ratings").
		Find(context.TODO(), bson.M{"user_id": userID})

	var results []bson.M
	cursor.All(context.TODO(), &results)

	c.JSON(200, results)
}