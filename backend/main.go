package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"wira-assignment/auth"
	"wira-assignment/cache"
	"wira-assignment/config"
	"wira-assignment/ranking"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"math"
	"fmt"
)

var (
	db          *sql.DB
	rankingRepo *ranking.Repository
	cfg         *config.Config
)

func init() {
	var err error
	// Load configuration
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize JWT key
	auth.InitJWTKey(cfg.JWTSecret)

	// Initialize Redis
	err = cache.InitRedis(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	}

	// Initialize database connection
	db, err = sql.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	rankingRepo = ranking.NewRepository(db)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize repositories and handlers
	rankingHandler := ranking.NewHandler(rankingRepo)

	api := r.Group("/api")
	{
		// Public routes
		api.GET("/classes", rankingHandler.GetClasses)

		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handleRegister)
			auth.POST("/login", handleLogin)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware())
		{
			protected.GET("/rankings", getRankings)
			protected.GET("/rankings/:classId", getRankingsByClass)
			protected.GET("/rankings/search", searchRankings)
			protected.GET("/characters", rankingHandler.GetUserCharacters)
			protected.POST("/characters", rankingHandler.CreateCharacter)
			protected.GET("/profile", getProfile)
			protected.PUT("/characters/:charId/score", updateScore)
		}
	}

	// Handle 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Not found"})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := auth.CreateUser(db, req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func handleLogin(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	user, err := auth.AuthenticateUser(db, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := auth.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func getProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var profile struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}

	// Try to get from cache
	cacheKey := fmt.Sprintf("profile:%v", userID)
	ctx := c.Request.Context()
	err := cache.Get(ctx, cacheKey, &profile)
	if err == nil {
		c.JSON(http.StatusOK, profile)
		return
	}

	// Get from database
	err = db.QueryRow(`
        SELECT acc_id, username, email, created_at
        FROM accounts
        WHERE acc_id = $1
    `, userID).Scan(&profile.ID, &profile.Username, &profile.Email, &profile.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
	}

	// Cache the result for 5 minutes
	err = cache.Set(ctx, cacheKey, profile, 5*time.Minute)
	if err != nil {
		log.Printf("Warning: Failed to cache profile: %v", err)
	}

	c.JSON(http.StatusOK, profile)
}

func getRankings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	class := c.Query("class")

	var classID int
	var err error

	if class != "all" && class != "" {
		classID, err = rankingRepo.GetClassIDByName(class)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class"})
			return
		}
	}

	rankings, total, err := rankingRepo.GetRankings(classID, page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rankings"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"rankings":     rankings,
		"total":        total,
		"current_page": page,
		"total_pages":  totalPages,
	})
}

func getRankingsByClass(c *gin.Context) {
	classIDStr := c.Param("classId")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	rankings, total, err := rankingRepo.GetRankings(classID, page, limit, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rankings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings": rankings,
		"total":    total,
	})
}

type CreateCharacterRequest struct {
	ClassID int `json:"class_id" binding:"required"`
}

func createCharacter(c *gin.Context) {
	var req CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	err := rankingRepo.CreateCharacter(userID.(int), req.ClassID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Character created successfully"})
}

type UpdateScoreRequest struct {
	Score int `json:"score" binding:"required"`
}

func updateScore(c *gin.Context) {
	charIDStr := c.Param("charId")
	charID, err := strconv.Atoi(charIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	var req UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rankingRepo.UpdateScore(charID, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Score updated successfully"})
}

func searchRankings(c *gin.Context) {
	username := c.Query("username")
	classIDStr := c.Query("classId")
	
	if username == "" && classIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide search parameters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Search rankings endpoint"})
}
