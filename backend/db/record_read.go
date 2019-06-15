package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type RecordReadOutput struct {
	ID         string
	LentAt     string
	ReturnedAt *string
	UserID     string
	StockID    string
}

// RecordRead : Read Record
func (h *Handler) RecordRead(ctx context.Context, opt *Option, out *[]*RecordReadOutput) error {

	filter := recordReadFilter(opt)

	columns := []string{
		"records.id",
		"records.lent_at",
		"records.returned_at",
		"records.user_id",
		"records.stock_id",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("records").
		LeftJoin("stocks on stocks.id = records.stock_id").
		LeftJoin("users on users.id = stocks.user_id").
		LeftJoin("books on books.id = stocks.book_id").
		LeftJoin("book_labels on book_labels.book_id = books.id").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "lent_at":
			by = "records.lent_at"
		case "returned_at":
			by = "records.returned_at"
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

	a := []*RecordReadOutput{}
	for cursor.Next() {
		var v RecordReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.LentAt,
			&v.ReturnedAt,
			&v.UserID,
			&v.StockID,
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

func recordReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"records.id": v[0],
				})
			case "user_Id":
				filter = append(filter, squirrel.Eq{
					"records.user_id": v[0],
				})
			case "stock_id":
				filter = append(filter, squirrel.Eq{
					"stocks.id": v[0],
				})
			case "owner_id":
				filter = append(filter, squirrel.Eq{
					"users.id": v[0],
				})
			case "owner":
				filter = append(filter, squirrel.Eq{
					"users.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			case "book_id":
				filter = append(filter, squirrel.Eq{
					"books.id": v[0],
				})
			case "title":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"books.title": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			case "isbn":
				filter = append(filter, squirrel.Eq{
					"books.isbn": v,
				})
			case "publisher":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"books.publisher": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			case "authors":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"authors.author": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			case "labels":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"labels.label": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			case "returned":
				switch v[0] {
				case "true":
					filter = append(filter, squirrel.NotEq{
						"records.returned_at": nil,
					})
				case "false":
					filter = append(filter, squirrel.Eq{
						"records.returned_at": nil,
					})
				}
			}
		}
	}

	return filter
}

func (h *Handler) RecordCount(ctx context.Context, opt *Option) (uint64, error) {

	filter := recordReadFilter(opt)

	builder := psql.
		Select("COUNT(*) AS count").
		Distinct().
		From("records").
		LeftJoin("stocks on stocks.id = records.stock_id").
		LeftJoin("users on users.id = stocks.user_id").
		LeftJoin("books on books.id = stocks.book_id").
		LeftJoin("book_labels on book_labels.book_id = books.id").
		Where(filter)

	query, args, err := builder.ToSql()
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return 0, err
	}

	zap.S().Debug(query, args)

	cursor, err := h.readDB.QueryContext(ctx, query, args...)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return 0, err
	}
	defer cursor.Close()

	var count uint64
	if !cursor.Next() {
		return 0, nil
	}
	err = cursor.Scan(
		&count,
	)

	return count, nil
}
