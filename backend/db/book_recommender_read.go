package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type BookRecommenderReadOutput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	BookID    string
	UserID    string
	UserName  string
}

// BookRecommenderRead : Read BookRecommender
func (h *Handler) BookRecommenderRead(ctx context.Context, opt *Option, out *[]*BookRecommenderReadOutput) error {

	filter := bookRecommenderReadFilter(opt)

	columns := []string{
		"book_recommenders.id",
		"book_recommenders.created_at",
		"book_recommenders.updated_at",
		"book_recommenders.book_id",
		"users.id",
		"users.name",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("book_recommenders").
		LeftJoin("users on users.id = book_recommenders.user_id").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "created_at":
			by = "book_recommenders.created_at"
		case "recommender":
			by = "users.name"
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

	a := []*BookRecommenderReadOutput{}
	for cursor.Next() {
		var v BookRecommenderReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.BookID,
			&v.UserID,
			&v.UserName,
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

func bookRecommenderReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"book_recommenders.id": v[0],
				})
			case "book_id":
				filter = append(filter, squirrel.Eq{
					"book_recommenders.book_id": v[0],
				})
			case "recommender_id":
				filter = append(filter, squirrel.Eq{
					"users.id": v[0],
				})
			case "recommender_name":
				filter = append(filter, squirrel.Eq{
					"users.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			}
		}
	}

	return filter
}
