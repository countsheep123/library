package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type StockReadOutput struct {
	ID           string
	CreatedAt    string
	UpdatedAt    string
	IsAvailable  bool
	BookID       string
	UserID       string
	UserName     string
	MarkID       string
	MarkName     string
	MarkURL      string
	LocationID   string
	LocationName string
}

// StockRead : Read Stock
func (h *Handler) StockRead(ctx context.Context, opt *Option, out *[]*StockReadOutput) error {

	filter := stockReadFilter(opt)

	columns := []string{
		"stocks.id",
		"stocks.created_at",
		"stocks.updated_at",
		"stocks.is_available",
		"stocks.book_id",
		"users.id",
		"users.name",
		"marks.id",
		"marks.name",
		"marks.url",
		"locations.id",
		"locations.name",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("stocks").
		LeftJoin("users on users.id = stocks.user_id").
		LeftJoin("marks on marks.id = stocks.mark_id").
		LeftJoin("locations on locations.id = stocks.location_id").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "created_at":
			by = "stocks.created_at"
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

	a := []*StockReadOutput{}
	for cursor.Next() {
		var v StockReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.IsAvailable,
			&v.BookID,
			&v.UserID,
			&v.UserName,
			&v.MarkID,
			&v.MarkName,
			&v.MarkURL,
			&v.LocationID,
			&v.LocationName,
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

func stockReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"stocks.id": v[0],
				})
			case "is_available":
				filter = append(filter, squirrel.Eq{
					"stocks.is_available": v[0],
				})
			case "book_id":
				filter = append(filter, squirrel.Eq{
					"stocks.book_id": v[0],
				})
			case "user_id":
				filter = append(filter, squirrel.Eq{
					"users.id": v[0],
				})
			case "user_name":
				filter = append(filter, squirrel.Eq{
					"users.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			case "mark_id":
				filter = append(filter, squirrel.Eq{
					"marks.id": v[0],
				})
			case "mark_name":
				filter = append(filter, squirrel.Eq{
					"marks.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			case "location_id":
				filter = append(filter, squirrel.Eq{
					"locations.id": v[0],
				})
			case "location_name":
				filter = append(filter, squirrel.Eq{
					"locations.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			}
		}
	}

	return filter
}
