package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type LocationReadOutput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Name      string
}

// LocationRead : Read Location
func (h *Handler) LocationRead(ctx context.Context, opt *Option, out *[]*LocationReadOutput) error {

	filter := locationReadFilter(opt)

	columns := []string{
		"locations.id",
		"locations.created_at",
		"locations.updated_at",
		"locations.name",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("locations").
		Where(filter)

	if opt.Sort != nil {
		by := ""
		switch opt.Sort.By {
		case "created_at":
			by = "locations.created_at"
		case "name":
			by = "locations.name"
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

	a := []*LocationReadOutput{}
	for cursor.Next() {
		var v LocationReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.Name,
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

func locationReadFilter(opt *Option) squirrel.And {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"locations.id": v[0],
				})
			case "name":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"locations.name": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			}
		}
	}

	return filter
}
