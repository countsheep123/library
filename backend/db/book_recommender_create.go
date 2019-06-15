package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type BookRecommenderCreateInput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	BookID    string
	UserID    string
}

func (in *BookRecommenderCreateInput) Validate() error {
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
	return nil
}

// BookRecommenderCreate : Create BookRecommender
func (h *Handler) BookRecommenderCreate(ctx context.Context, in *BookRecommenderCreateInput) error {

	if err := in.Validate(); err != nil {
		return err
	}

	columns := []string{
		"id",
		"created_at",
		"updated_at",
		"book_id",
		"user_id",
	}
	values := []interface{}{
		in.ID,
		in.CreatedAt,
		in.UpdatedAt,
		in.BookID,
		in.UserID,
	}

	if err := h.Insert(ctx, "book_recommenders", columns, values); err != nil {
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
