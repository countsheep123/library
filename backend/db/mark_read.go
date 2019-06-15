package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type MarkReadOutput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Name      string
	URL       *string
	UserID    string
	UserName  string
}

// MarkRead : Read Mark
func (h *Handler) MarkRead(ctx context.Context, opt *Option, out *[]*MarkReadOutput) error {

	filter := markReadFilter(opt)

	columns := []string{
		"marks.id",
		"marks.created_at",
		"marks.updated_at",
		"marks.name",
		"marks.url",
		"users.id",
		"users.name",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("marks").
		LeftJoin("users on users.id = marks.user_id").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "created_at":
			by = "marks.created_at"
		case "name":
			by = "marks.name"
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

	a := []*MarkReadOutput{}
	for cursor.Next() {
		var v MarkReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.Name,
			&v.URL,
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

func markReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"marks.id": v[0],
				})
			case "owner_id":
				filter = append(filter, squirrel.Eq{
					"users.id": v[0],
				})
			case "owner":
				filter = append(filter, squirrel.Eq{
					"users.name": fmt.Sprintf("%%%s%%", sanitize(v[0])),
				})
			case "name":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"marks.name": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			}
		}
	}

	return filter
}
