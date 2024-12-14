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
		c.Set("claims", claims)
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
			auth.POST("/2fa/login/verify", handle2FALoginVerify)
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
			protected.POST("/cache/clear", handleClearCache)

			// 2FA endpoints
			protected.POST("/2fa/enable", handleEnable2FA)
			protected.POST("/2fa/verify", handleVerify2FA)
			protected.POST("/2fa/disable", handleDisable2FA)
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
		if strings.Contains(err.Error(), "username already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		if strings.Contains(err.Error(), "email already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}


type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token,omitempty"`
	User     *auth.User `json:"user,omitempty"`
	Requires2FA bool `json:"requires_2fa,omitempty"`
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

	// Check if user has 2FA enabled
	enabled, _, err := auth.GetUser2FAStatus(db, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check 2FA status"})
		return
	}

	if enabled {
		c.JSON(http.StatusOK, LoginResponse{
			Requires2FA: true,
		})
		return
	}

	token, err := auth.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

type Verify2FALoginRequest struct {
	Username string `json:"username" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func handle2FALoginVerify(c *gin.Context) {
	var req Verify2FALoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	var user auth.User
	var secret string
	err := db.QueryRow(`
        SELECT acc_id, username, email, two_factor_secret 
        FROM accounts 
        WHERE username = $1
    `, req.Username).Scan(&user.ID, &user.Username, &user.Email, &secret)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate 2FA code
	if !auth.ValidateTOTP(secret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  &user,
	})
}

func getProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var profile struct {
		ID               int       `json:"id"`
		Username         string    `json:"username"`
		Email            string    `json:"email"`
		CreatedAt        time.Time `json:"created_at"`
		TwoFactorEnabled bool      `json:"two_factor_enabled"`
	}

	err := db.QueryRow(`
        SELECT acc_id, username, email, created_at, two_factor_enabled
        FROM accounts
        WHERE acc_id = $1
    `, userID).Scan(&profile.ID, &profile.Username, &profile.Email, &profile.CreatedAt, &profile.TwoFactorEnabled)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
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

func handleClearCache(c *gin.Context) {
    ctx := c.Request.Context()
    err := cache.ClearAll(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cache"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Cache cleared successfully"})
}

type Enable2FARequest struct {
    Password string `json:"password"`
}

type Verify2FARequest struct {
    Code   string `json:"code"`
    Secret string `json:"secret"`
}

func handleEnable2FA(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    userClaims := claims.(*auth.Claims)
    var req Enable2FARequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Verify password
    var passwordHash string
    err := db.QueryRow("SELECT password_hash FROM accounts WHERE acc_id = $1", userClaims.UserID).Scan(&passwordHash)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
        return
    }

    if !auth.CheckPasswordHash(req.Password, passwordHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
        return
    }

    // Generate 2FA secret
    secret, err := auth.GenerateSecret()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate secret"})
        return
    }

    // Create QR code
    qrURL := fmt.Sprintf("otpauth://totp/Wira:%s?secret=%s&issuer=Wira", userClaims.Username, secret)
    
    c.JSON(http.StatusOK, gin.H{
        "secret": secret,
        "qr_url": qrURL,
    })
}

func handleVerify2FA(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    userClaims := claims.(*auth.Claims)
    var req Verify2FARequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Validate TOTP code using the provided secret
    if !auth.ValidateTOTP(req.Secret, req.Code) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
        return
    }

    // Enable 2FA with the validated secret
    if err := auth.Enable2FA(db, userClaims.UserID, req.Secret); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enable 2FA"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

func handleDisable2FA(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    userClaims := claims.(*auth.Claims)
    err := auth.Disable2FA(db, userClaims.UserID)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "2FA not enabled for this user"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disable 2FA"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}
