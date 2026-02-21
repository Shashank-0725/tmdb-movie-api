package main

import (
	"movie-api/internal/auth"
	"movie-api/internal/config"
	"movie-api/internal/database"
	"movie-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	cfg := config.Load()

	auth.Secret = []byte(cfg.JWTSecret)
	database.Init(cfg.MongoURI)

	r := gin.Default()
	r.Use(cors.Default())

	// pass TMDB key into context
	r.Use(func(c *gin.Context) {
		c.Set("tmdbKey", cfg.TMDBAPIKey)
		c.Next()
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/movies/search", handlers.Search)

	protected := r.Group("/")
	protected.Use(auth.Middleware())
	protected.GET("/recommend", handlers.Recommend)
	protected.POST("/rate/:id", handlers.RateMovie)
	protected.GET("/rated", handlers.GetRated)
	protected.POST("/watchlist/:id", handlers.AddWatchlist)
	protected.GET("/watchlist", handlers.GetWatchlist)

	r.Run(":" + cfg.Port)
}