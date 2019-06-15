package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type StockUpdateInput struct {
	UpdatedAt   string
	IsAvailable *bool
	MarkID      *string
	LocationID  *string
}

func (in *StockUpdateInput) Validate() error {
	if len(in.UpdatedAt) == 0 {
		return obj.Internal{
			Msg: "updated_at is required",
		}
	}
	if in.MarkID != nil && len(*in.MarkID) != 20 {
		return obj.Internal{
			Msg: "mark_id is required",
		}
	}
	if in.LocationID != nil && len(*in.LocationID) != 20 {
		return obj.Internal{
			Msg: "location_id is required",
		}
	}
	return nil
}

// StockUpdate : Update Stock
func (h *Handler) StockUpdate(ctx context.Context, in *StockUpdateInput, filters map[string]string) error {

	if err := in.Validate(); err != nil {
		return err
	}

	kv := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}

	if in.IsAvailable != nil {
		kv["is_available"] = *in.IsAvailable
	}
	if in.MarkID != nil {
		kv["mark_id"] = *in.MarkID
	}
	if in.LocationID != nil {
		kv["location_id"] = *in.LocationID
	}

	if err := h.Update(ctx, "stocks", kv, filters); err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	return nil
}
