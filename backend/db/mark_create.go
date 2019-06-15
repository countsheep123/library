package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type MarkCreateInput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	UserID    string
	Name      string
	URL       *string
}

func (in *MarkCreateInput) Validate() error {
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
	if len(in.UserID) != 20 {
		return obj.Internal{
			Msg: "user_id is required",
		}
	}
	if len(in.Name) == 0 {
		return obj.Internal{
			Msg: "name is required",
		}
	}
	return nil
}

// MarkCreate : Create Mark
func (h *Handler) MarkCreate(ctx context.Context, in *MarkCreateInput) error {

	if err := in.Validate(); err != nil {
		return err
	}

	columns := []string{
		"id",
		"created_at",
		"updated_at",
		"user_id",
		"name",
		"url",
	}
	values := []interface{}{
		in.ID,
		in.CreatedAt,
		in.UpdatedAt,
		in.UserID,
		in.Name,
		in.URL,
	}

	if err := h.Insert(ctx, "marks", columns, values); err != nil {
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
