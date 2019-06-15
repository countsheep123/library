package db

import (
	"context"

	"github.com/getsentry/raven-go"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/countsheep123/library/obj"
)

type BookUpdateInput struct {
	UpdatedAt string
	Title     *string
	ISBN      *string
	Publisher *string
	Pubdate   *string
	Authors   []string
	Cover     *string
}

func (in *BookUpdateInput) Validate() error {
	if len(in.UpdatedAt) == 0 {
		return obj.Internal{
			Msg: "updated_at is required",
		}
	}
	if in.Title != nil && len(*in.Title) == 0 {
		return obj.Internal{
			Msg: "title is required",
		}
	}
	return nil
}

// BookUpdate : Update Book
func (h *Handler) BookUpdate(ctx context.Context, in *BookUpdateInput, filters map[string]string) error {

	if err := in.Validate(); err != nil {
		return err
	}

	kv := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}

	if in.Title != nil {
		kv["title"] = *in.Title
	}
	if in.ISBN != nil {
		kv["isbn"] = *in.ISBN
	}
	if in.Publisher != nil {
		kv["publisher"] = *in.Publisher
	}
	if in.Pubdate != nil {
		kv["pubdate"] = *in.Pubdate
	}
	if in.Authors != nil {
		kv["authors"] = pq.Array(in.Authors)
	}
	if in.Cover != nil {
		kv["cover_url"] = *in.Cover
	}

	if err := h.Update(ctx, "books", kv, filters); err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	return nil
}
