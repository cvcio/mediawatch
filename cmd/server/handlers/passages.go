package handlers

import (
	"net/http"
	"strconv"

	"github.com/cvcio/mediawatch/models/passage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
	"github.com/cvcio/mediawatch/pkg/web"
	"go.opencensus.io/trace"
)

// Passages is the handler struct for org related enbpoints
type Passages struct {
	DB  *db.MongoDB
	log *zap.Logger
}

// NewPassagesHandler returns a new instance of Passages
func NewPassagesHandler(log *zap.Logger, db *db.MongoDB) (*Passages, error) {
	return &Passages{
		DB:  db,
		log: log,
	}, nil
}

// Create new passage
func (u *Passages) Create(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Passages.Create")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	var newPassage *passagesv2.Passage
	if err := c.ShouldBindJSON(&newPassage); err != nil {
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	created, err := passage.Create(ctx, u.DB, newPassage)
	if err != nil {
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, created)
}

// GetPassages returns a list of passages based on query parameters
func (u *Passages) GetPassages(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Passages.GetPassages")
	defer span.End()

	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// Build query from parameters
	query := &passagesv2.QueryPassage{
		Language: c.Query("language"),
		Type:     passagesv2.PassageType(parseInt32(c.Query("type"), 0)),
	}

	// Parse pagination parameters
	query.Limit = parseInt32(c.Query("limit"), 120)
	query.Offset = parseInt32(c.Query("offset"), 0)

	// Parse sorting parameters
	query.SortKey = c.DefaultQuery("sort_key", "text")
	query.SortOrder = parseInt32(c.Query("sort_order"), 1)

	passages, err := passage.List(ctx, u.DB, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, passages)
}

func (u *Passages) DeletePassage(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Passages.DeletePassage")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	passageID := c.Param("id")
	if passageID == "" {
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(nil))
		return
	}

	if err := passage.Delete(ctx, u.DB, &passagesv2.Passage{Id: passageID}); err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, web.Deleted)
}

// Helper function to parse int32 with default value
func parseInt32(s string, defaultVal int32) int32 {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultVal
	}

	return int32(val)
}
