package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type RecordUpdateInput struct {
	// LentAt string
	ReturnedAt string
	// UserID     string
	// StockID    string
}

func (in *RecordUpdateInput) Validate() error {
	if len(in.ReturnedAt) == 0 {
		return obj.Internal{
			Msg: "returned_at is required",
		}
	}
	return nil
}

// RecordUpdate : Update Record
func (h *Handler) RecordUpdate(ctx context.Context, in *RecordUpdateInput, filters map[string]string) error {

	if err := in.Validate(); err != nil {
		return err
	}

	kv := map[string]interface{}{
		"returned_at": in.ReturnedAt,
	}

	if err := h.Update(ctx, "records", kv, filters); err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	return nil
}
