package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type RecordCreateInput struct {
	ID     string
	LentAt string
	// ReturnedAt string
	UserID  string
	StockID string
}

func (in *RecordCreateInput) Validate() error {
	if len(in.ID) != 20 {
		return obj.Internal{
			Msg: "id is required",
		}
	}
	if len(in.LentAt) == 0 {
		return obj.Internal{
			Msg: "lent_at is required",
		}
	}
	if len(in.UserID) != 20 {
		return obj.Internal{
			Msg: "user_id is required",
		}
	}
	if len(in.StockID) != 20 {
		return obj.Internal{
			Msg: "stock_id is required",
		}
	}
	return nil
}

// RecordCreate : Create Record
func (h *Handler) RecordCreate(ctx context.Context, in *RecordCreateInput) error {

	if err := in.Validate(); err != nil {
		return err
	}

	columns := []string{
		"id",
		"lent_at",
		"user_id",
		"stock_id",
	}
	values := []interface{}{
		in.ID,
		in.LentAt,
		in.UserID,
		in.StockID,
	}

	if err := h.Insert(ctx, "records", columns, values); err != nil {
		switch err.(type) {
		case obj.Duplicate:
			zap.S().Warn(err)
		default:
			raven.CaptureError(err, nil)
			zap.S().Error(err)
		}
		return err
	}

	return nil
}
