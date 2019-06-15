package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type StockCreateInput struct {
	ID          string
	CreatedAt   string
	UpdatedAt   string
	IsAvailable bool
	BookID      string
	UserID      string
	MarkID      string
	LocationID  string
}

func (in *StockCreateInput) Validate() error {
	if len(in.ID) != 20 {
		return obj.Internal{
			Msg: "id is required",
		}
	}
	if len(in.CreatedAt) == 0 {
		return obj.Internal{
			Msg: "created_at is required",
		}
	}
	if len(in.UpdatedAt) == 0 {
		return obj.Internal{
			Msg: "updated_at is required",
		}
	}
	if len(in.BookID) != 20 {
		return obj.Internal{
			Msg: "book_id is required",
		}
	}
	if len(in.UserID) != 20 {
		return obj.Internal{
			Msg: "user_id is required",
		}
	}
	if len(in.MarkID) != 20 {
		return obj.Internal{
			Msg: "mark_id is required",
		}
	}
	if len(in.LocationID) != 20 {
		return obj.Internal{
			Msg: "location_id is required",
		}
	}
	return nil
}

// StockCreate : Create Stock
func (h *Handler) StockCreate(ctx context.Context, in *StockCreateInput) error {

	if err := in.Validate(); err != nil {
		return err
	}

	columns := []string{
		"id",
		"created_at",
		"updated_at",
		"is_available",
		"book_id",
		"user_id",
		"mark_id",
		"location_id",
	}
	values := []interface{}{
		in.ID,
		in.CreatedAt,
		in.UpdatedAt,
		in.IsAvailable,
		in.BookID,
		in.UserID,
		in.MarkID,
		in.LocationID,
	}

	if err := h.Insert(ctx, "stocks", columns, values); err != nil {
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
