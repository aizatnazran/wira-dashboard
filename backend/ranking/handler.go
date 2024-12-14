package ranking

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetClasses(c *gin.Context) {
	classes, err := h.repo.GetClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) GetUserCharacters(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	characters, err := h.repo.GetUserCharacters(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, characters)
}

func (h *Handler) CreateCharacter(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		ClassID int `json:"class_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.repo.CreateCharacter(userID.(int), req.ClassID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "character created successfully"})
}

func (h *Handler) GetRankings(c *gin.Context) {
	// Get query parameters with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")
	classID, _ := strconv.Atoi(c.DefaultQuery("classId", "0"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var rankings []RankingEntry
	var totalCount int
	var err error

	rankings, totalCount, err = h.repo.GetRankings(classID, page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	response := RankingResponse{
		Rankings:    rankings,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateScore(c *gin.Context) {
	var req struct {
		CharID int `json:"char_id" binding:"required"`
		Score  int `json:"score" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.repo.UpdateScore(req.CharID, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "score updated successfully"})
}
