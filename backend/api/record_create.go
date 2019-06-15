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

type RecordCreateRequest struct {
	StockID string `json:"stock_id"`
}

func (r *RecordCreateRequest) Validate() error {
	if len(r.StockID) != 20 {
		return obj.InvalidRequest{
			Msg: "stock_id is required",
		}
	}
	return nil
}

type RecordCreateResponse struct {
	ID      string `json:"id"`
	LentAt  string `json:"lent_at"`
	UserID  string `json:"user_id"`
	StockID string `json:"stock_id"`
}

// RecordCreate : Create Record
func (s *Server) RecordCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	var req *RecordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	id := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		// 在庫チェック
		opt := &db.Option{
			Filters: map[string][]string{
				"id": []string{req.StockID},
			},
		}

		var stocks []*db.StockReadOutput
		if err := tx.StockRead(ctx, opt, &stocks); err != nil {
			return err
		}

		if len(stocks) == 0 {
			return obj.InvalidRequest{}
		}

		if !stocks[0].IsAvailable {
			return obj.InvalidRequest{
				Msg: "not available",
			}
		}

		// 在庫おさえる
		available := false
		if err := tx.StockUpdate(ctx, &db.StockUpdateInput{
			UpdatedAt:   ts,
			IsAvailable: &available,
		}, map[string]string{
			"id": req.StockID,
		}); err != nil {
			return err
		}

		// 貸出処理
		if err := tx.RecordCreate(ctx, &db.RecordCreateInput{
			ID:      id,
			LentAt:  ts,
			UserID:  user.ID,
			StockID: req.StockID,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &RecordCreateResponse{
		ID:      id,
		LentAt:  ts,
		UserID:  user.ID,
		StockID: req.StockID,
	})
}
