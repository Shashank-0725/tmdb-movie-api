package config

import (
	"os"
)

type Config struct {
	Port        string
	MongoURI    string
	JWTSecret   string
	TMDBAPIKey  string
}

func Load() *Config {
	return &Config{
		Port:       os.Getenv("PORT"),
		MongoURI:   os.Getenv("MONGO_URI"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		TMDBAPIKey: os.Getenv("TMDB_API_KEY"),
	}
}