package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"wira-assignment/auth"
	"wira-assignment/cache"
	"wira-assignment/config"
	"wira-assignment/ranking"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://77.237.243.104:3000", "https://wira.aizat.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "x-session-id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:          12 * time.Hour,
	}))

	// Public routes
	authRouter := r.Group("/api/auth")
	{
		authRouter.POST("/register", handleRegister)
		authRouter.POST("/login", handleLogin)
		authRouter.POST("/2fa/login/verify", handle2FALogin)
		
		// Session validation endpoint
		authRouter.POST("/validate-session", func(c *gin.Context) {
			var req struct {
				SessionID string `json:"sessionID" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			session, err := auth.ValidateSession(db, req.SessionID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"valid": true,
				"session": session,
			})
		})

		// Logout endpoint
		authRouter.POST("/logout", func(c *gin.Context) {
			sessionID := c.GetHeader("X-Session-ID")
			if sessionID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
				return
			}

			err := auth.DeleteSession(db, sessionID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
		})
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		api.GET("/profile", getProfile)
		api.POST("/cache/clear", handleClearCache)
		api.GET("/rankings", getRankings)
		api.GET("/rankings/:class", getRankingsByClass)
		api.POST("/characters", createCharacter)
		api.PUT("/characters/:id/score", updateScore)
		api.GET("/search", searchRankings)
		api.GET("/classes", func(c *gin.Context) {
			classes, err := rankingRepo.GetClasses()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
				return
			}
			c.JSON(http.StatusOK, classes)
		})

		// 2FA routes
		api.POST("/2fa/enable", handleEnable2FA)
		api.POST("/2fa/verify", handleVerify2FA)
		api.POST("/2fa/disable", handleDisable2FA)
	}

	// Start cleanup goroutine for expired sessions
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			if err := auth.DeleteExpiredSessions(db); err != nil {
				log.Printf("Failed to cleanup expired sessions: %v", err)
			}
		}
	}()

	// Start server
	if err := r.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatal(err)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Sanitize inputs
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

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
	SessionID string `json:"sessionID,omitempty"`
}

func handleLogin(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := auth.AuthenticateUser(db, loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if 2FA is enabled
	twoFactorEnabled, _, err := auth.GetUser2FAStatus(db, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check 2FA status"})
		return
	}

	if twoFactorEnabled {
		c.JSON(http.StatusOK, LoginResponse{
			Requires2FA: true,
		})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Create session
	session, err := auth.CreateSession(db, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:    token,
		User:     user,
		SessionID: session.SessionID,
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
	classIDStr := c.Param("class")
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
	charIDStr := c.Param("id")
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

func handle2FALogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user and 2FA secret
	var user auth.User
	var secret string
	err := db.QueryRow(`
		SELECT acc_id, username, email, password_hash, two_factor_secret 
		FROM accounts 
		WHERE username = $1 AND two_factor_enabled = true`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate TOTP code
	if !auth.ValidateTOTP(secret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Create session
	session, err := auth.CreateSession(db, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		User:      &user,
		SessionID: session.SessionID,
	})
}
