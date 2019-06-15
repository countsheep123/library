package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

func (h *Handler) UserDelete(ctx context.Context, filters map[string]string) error {

	if err := h.Delete(ctx, "users", filters); err != nil {
		switch err.(type) {
		case obj.NotFound:
			zap.S().Warn(err)
		default:
			raven.CaptureError(err, nil)
			zap.S().Error(err)
		}
		return err
	}

	return nil
}
