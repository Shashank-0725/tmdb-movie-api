# 🎬 Go Movie Watchlist & Recommendation API

A professional REST API built with Go that integrates with the TMDB API to allow users to search movies, manage watchlists, rate movies, and receive personalized recommendations.

This project focuses on backend architecture, authentication, caching strategy, and external API integration.

## 📌 Problem Explanation

Users often struggle to track movies they want to watch and discover new content aligned with their preferences.

This API solves that by:
- Allowing users to create accounts
- Searching real-time movie data from TMDB
- Saving movies into personal watchlists
- Rating watched movies
- Generating personalized recommendations based on genre preferences
- Reducing redundant external API calls through a caching mechanism

## 🚀 Key Features
- User Registration & Login (JWT Authentication)
- Search Movies using TMDB API
- Add Movies to Watchlist
- Rate Movies (1–5)
- Personalized Recommendations
- In-memory caching for external API responses
- Modular project structure using Go best practices
- MongoDB persistence for user data

## 🏗️ Project Architecture

The project follows a clean modular structure:

```
movie-api/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── auth/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── tmdb/
│
├── go.mod
└── .env
```

**Architecture Design Principles**
- Separation of concerns
- Modular internal packages
- Clean HTTP handler layer
- Centralized configuration management
- External API abstraction layer (tmdb package)

## 🔐 Authentication System

JWT-based authentication is implemented.

**Flow:**
- User registers
- Password is hashed using bcrypt
- User logs in
- JWT token is generated
- Protected routes require Authorization header

Example header:
```
Authorization: <JWT_TOKEN>
```

## 🎥 External API Integration (TMDB)

The system integrates with TMDB using an API key.

**Used endpoints:**
- Movie search
- Movie details
- Genre-based discovery

TMDB handles real movie metadata, while the API stores only user-specific data.

## ⚡ Caching Strategy

To reduce repeated TMDB calls:
- In-memory map cache is implemented
- Each cached response has expiration time (10 minutes)
- If cached data exists and is valid → return from cache
- Otherwise → fetch from TMDB and store in cache

**Why In-Memory Cache?**
- Simple and fast
- Suitable for single-instance backend
- Avoids unnecessary external API calls

This improves:
- Performance
- Rate-limit safety
- Response time

## 🎯 Recommendation Algorithm

The recommendation system works as follows:
- Fetch user's rated movies
- Select movies rated ≥ 4
- Extract genre from first highly-rated movie
- Call TMDB discover endpoint using that genre
- Return recommended movies

This creates a basic personalized recommendation engine.

## 🛠️ How to Run

1️⃣ **Install Go & MongoDB**

Ensure:
- Go installed
- MongoDB running locally or remote URI available

2️⃣ **Create .env**
```
PORT=8080
MONGO_URI=mongodb://localhost:27017
JWT_SECRET=your_secret_key
TMDB_API_KEY=your_tmdb_api_key
```

3️⃣ **Run Server**

From project root:
```
go mod tidy
go run ./cmd
```

Server runs at:
```
http://localhost:8080
```

## 📡 Sample API Usage

**Register**
```
POST /register
{
  "email": "test@mail.com",
  "password": "123456"
}
```

**Login**
```
POST /login
```
Returns JWT token.

**Search Movies**
```
GET /movies/search?q=batman
```

**Add to Watchlist (Protected)**
```
POST /watchlist/278
```
Header:
```
Authorization: <token>
```

**Rate Movie (Protected)**
```
POST /rate/278
{
  "rating": 5
}
```

**Get Recommendations (Protected)**
```
GET /recommend
```

## 🧠 Design Decisions

- Clean Modular Structure
- MongoDB for users, watchlists, ratings
- JWT Authentication for stateless auth
- In-Memory Cache for TMDB responses

## 🔍 Database Schema Overview

**Collections:**

- users
  - _id
  - email
  - password (hashed)

- watchlists
  - user_id
  - movie_id
  - title
  - release_date
  - genres
  - poster_path
  - vote_average

- ratings
  - user_id
  - movie_id
  - title
  - release_date
  - genres
  - user_rating

## 🧪 Testing

You can test using:
- Postman
- curl
- Browser for search endpoint

## 🤖 AI Assistance Disclosure

This project was developed with assistance from AI tools to accelerate development and improve documentation clarity.

AI assistance was used for:
- Structuring project architecture
- Generating initial boilerplate code
- Improving README documentation clarity

However, the following were performed independently:
- Understanding all generated code
- Implementing and verifying authentication flow
- Designing caching logic
- Debugging and testing API endpoints
- Structuring final modular architecture

AI was used as a productivity assistant, while system design understanding and implementation decisions remain my own.
