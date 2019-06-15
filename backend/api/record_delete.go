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

func (s *Server) RecordDelete(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	recordID := c.Param("record_id")

	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		// check record
		opt := &db.Option{
			Filters: map[string][]string{
				"id": []string{recordID},
			},
		}

		var records []*db.RecordReadOutput
		if err := tx.RecordRead(ctx, opt, &records); err != nil {
			return err
		}

		if len(records) == 0 {
			return obj.NotFound{}
		}

		// 在庫解放
		available := true
		if err := tx.StockUpdate(ctx, &db.StockUpdateInput{
			UpdatedAt:   ts,
			IsAvailable: &available,
		}, map[string]string{
			"id": records[0].StockID,
		}); err != nil {
			return err
		}

		// 返却処理
		if err := tx.RecordUpdate(ctx, &db.RecordUpdateInput{
			ReturnedAt: ts,
		}, map[string]string{
			"id":      recordID,
			"user_id": user.ID,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
