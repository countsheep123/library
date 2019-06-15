package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type StockCreateRequest struct {
	MarkID     string `json:"mark_id"`
	LocationID string `json:"location_id"`
}

func (r *StockCreateRequest) Validate() error {
	if len(r.MarkID) != 20 {
		return obj.InvalidRequest{
			Msg: "mark_id is required",
		}
	}
	if len(r.LocationID) != 20 {
		return obj.InvalidRequest{
			Msg: "location_id is required",
		}
	}
	return nil
}

type StockCreateResponse struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsAvailable bool   `json:"is_available"`
	BookID      string `json:"book_id"`
	UserID      string `json:"user_id"`
	MarkID      string `json:"mark_id"`
	LocationID  string `json:"location_id"`
}

// StockCreate : Create Stock
func (s *Server) StockCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	var req *StockCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	bookID := c.Param("book_id")

	stockID := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())
	availability := true

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.StockCreate(ctx, &db.StockCreateInput{
			ID:          stockID,
			CreatedAt:   ts,
			UpdatedAt:   ts,
			IsAvailable: availability,
			BookID:      bookID,
			UserID:      user.ID,
			MarkID:      req.MarkID,
			LocationID:  req.LocationID,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &StockCreateResponse{
		ID:          stockID,
		CreatedAt:   ts,
		UpdatedAt:   ts,
		IsAvailable: availability,
		BookID:      bookID,
		UserID:      user.ID,
		MarkID:      req.MarkID,
		LocationID:  req.LocationID,
	})
}
