package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type BookReadOutput struct {
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

// BookRead : Read Book
func (h *Handler) BookRead(ctx context.Context, opt *Option, out *[]*BookReadOutput) error {

	filter := bookReadFilter(opt)

	columns := []string{
		"books.id",
		"books.created_at",
		"books.updated_at",
		"books.title",
		"books.isbn",
		"books.publisher",
		"books.pubdate",
		"books.authors",
		"books.cover_url",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("books").
		LeftJoin("book_labels on book_labels.book_id = books.id").
		LeftJoin("stocks on stocks.book_id = books.id").
		LeftJoin("users on users.id = stocks.user_id").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "created_at":
			by = "books.created_at"
		case "title":
			by = "books.title"
		case "isbn":
			by = "books.isbn"
		case "publisher":
			by = "books.publisher"
		case "pubdate":
			by = "books.pubdate"
		case "authors":
			by = "books.authors"
		case "owner":
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

	a := []*BookReadOutput{}
	for cursor.Next() {
		var v BookReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.Title,
			&v.ISBN,
			&v.Publisher,
			&v.Pubdate,
			pq.Array(&v.Authors),
			&v.Cover,
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

func bookReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
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
			case "owner_id":
				filter = append(filter, squirrel.Eq{
					"users.id": v[0],
				})
			case "owner":
				filter = append(filter, squirrel.Eq{
					"users.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			}
		}
	}

	return filter
}

func (h *Handler) BookCount(ctx context.Context, opt *Option) (uint64, error) {

	filter := bookReadFilter(opt)

	builder := psql.
		Select("COUNT(*) AS count").
		Distinct().
		From("books").
		LeftJoin("book_labels on book_labels.book_id = books.id").
		LeftJoin("stocks on stocks.book_id = books.id").
		LeftJoin("users on users.id = stocks.user_id").
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
