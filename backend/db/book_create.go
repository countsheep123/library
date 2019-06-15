package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type BookCreateInput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Title     string
	ISBN      *string
	Publisher *string
	Pubdate   *string
	Authors   []string
	Cover     *string
}

func (in *BookCreateInput) Validate() error {
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
	if len(in.Title) == 0 {
		return obj.Internal{
			Msg: "title is required",
		}
	}
	if in.Authors == nil {
		return obj.Internal{
			Msg: "authors must not be nil",
		}
	}
	return nil
}

// BookCreate : Create Book
func (h *Handler) BookCreate(ctx context.Context, in *BookCreateInput) error {

	if err := in.Validate(); err != nil {
		return err
	}

	columns := []string{
		"id",
		"created_at",
		"updated_at",
		"title",
		"isbn",
		"publisher",
		"pubdate",
		"authors",
		"cover_url",
	}
	values := []interface{}{
		in.ID,
		in.CreatedAt,
		in.UpdatedAt,
		in.Title,
		in.ISBN,
		in.Publisher,
		in.Pubdate,
		pq.Array(in.Authors),
		in.Cover,
	}

	if err := h.Insert(ctx, "books", columns, values); err != nil {
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
