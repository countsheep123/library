package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type BookLabelReadOutput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	BookID    string
	UserID    string
	Label     string
}

// BookLabelRead : Read BookLabel
func (h *Handler) BookLabelRead(ctx context.Context, opt *Option, out *[]*BookLabelReadOutput) error {

	filter := bookLabelReadFilter(opt)

	columns := []string{
		"book_labels.id",
		"book_labels.created_at",
		"book_labels.updated_at",
		"book_labels.book_id",
		"book_labels.user_id",
		"book_labels.label",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("book_labels").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "created_at":
			by = "book_labels.created_at"
		case "label":
			by = "book_labels.label"
		}
		if by != "" {
			var order string
			if opt.Sort.IsAsc {
				order = fmt.Sprintf("%s ASC", by)
			} else {
				order = fmt.Sprintf("%s DESC", by)
			}

			builder = builder.OrderBy(order)
		}
	}

	if opt.Range != nil {
		builder = builder.Limit(opt.Range.Limit).Offset(opt.Range.Offset)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	zap.S().Debug(query, args)

	cursor, err := h.readDB.QueryContext(ctx, query, args...)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}
	defer cursor.Close()

	a := []*BookLabelReadOutput{}
	for cursor.Next() {
		var v BookLabelReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.BookID,
			&v.UserID,
			&v.Label,
		); err != nil {
			raven.CaptureError(err, nil)
			zap.S().Error(err)
			return err
		}
		a = append(a, &v)
	}

	*out = a

	return nil
}

func bookLabelReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"book_labels.id": v[0],
				})
			case "book_id":
				filter = append(filter, squirrel.Eq{
					"book_labels.book_id": v[0],
				})
			case "user_id":
				filter = append(filter, squirrel.Eq{
					"book_labels.user_id": v[0],
				})
			case "label":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"book_labels.label": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			}
		}
	}

	return filter
}
