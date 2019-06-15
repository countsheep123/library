package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StockUpdateRequest struct {
	IsAvailable *bool   `json:"is_available"`
	MarkID      *string `json:"mark_id"`
	LocationID  *string `json:"location_id"`
}

func (r *StockUpdateRequest) Validate() error {
	if r.MarkID != nil && len(*r.MarkID) != 20 {
		return obj.InvalidRequest{
			Msg: "mark_id is required",
		}
	}
	if r.LocationID != nil && len(*r.LocationID) != 20 {
		return obj.InvalidRequest{
			Msg: "location_id is required",
		}
	}
	return nil
}

type StockUpdateResponse struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsAvailable bool   `json:"is_available"`
	BookID      string `json:"book_id"`
	UserID      string `json:"user_id"`
	MarkID      string `json:"mark_id"`
	LocationID  string `json:"location_id"`
}

// StockUpdate : Update Stock
func (s *Server) StockUpdate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	bookID := c.Param("book_id")
	stockID := c.Param("stock_id")

	var req *StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	ts := fmt.Sprint(time.Now().UnixNano())

	res := &StockUpdateResponse{
		ID:     stockID,
		BookID: bookID,
		UserID: user.ID,
	}

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		opt := &db.Option{
			Filters: map[string][]string{
				"id": []string{stockID},
			},
		}
		var stocks []*db.StockReadOutput
		if err := tx.StockRead(ctx, opt, &stocks); err != nil {
			return err
		}

		if len(stocks) == 0 {
			return obj.NotFound{}
		}

		res.CreatedAt = stocks[0].CreatedAt
		res.UpdatedAt = ts

		filters := map[string]string{
			"id": stockID,
		}

		if err := tx.StockUpdate(ctx, &db.StockUpdateInput{
			UpdatedAt:   ts,
			IsAvailable: req.IsAvailable,
			MarkID:      req.MarkID,
			LocationID:  req.LocationID,
		}, filters); err != nil {
			return err
		}

		if req.IsAvailable != nil {
			res.IsAvailable = *req.IsAvailable
		} else {
			res.IsAvailable = stocks[0].IsAvailable
		}

		if req.MarkID != nil {
			res.MarkID = *req.MarkID
		} else {
			res.MarkID = stocks[0].MarkID
		}

		if req.LocationID != nil {
			res.LocationID = *req.LocationID
		} else {
			res.LocationID = stocks[0].LocationID
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
