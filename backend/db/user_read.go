package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type UserReadOutput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Name      string
	Company   *string
	Email     string
	IsAdmin   bool
}

func (h *Handler) UserRead(ctx context.Context, opt *Option, out *[]*UserReadOutput) error {

	filter, err := userReadFilter(opt)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	columns := []string{
		"users.id",
		"users.created_at",
		"users.updated_at",
		"users.name",
		"users.company",
		"users.email",
		"users.is_admin",
	}

	builder := psql.
		Select(columns...).
		Distinct().
		From("users").
		Where(filter)

	if opt.Sort != nil {
		if opt.Sort.By != "" {
			var order string
			if opt.Sort.IsAsc {
				order = fmt.Sprintf("%s ASC", opt.Sort.By)
			} else {
				order = fmt.Sprintf("%s DESC", opt.Sort.By)
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

	a := []*UserReadOutput{}
	for cursor.Next() {
		var v UserReadOutput
		if err := cursor.Scan(
			&v.ID,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.Name,
			&v.Company,
			&v.Email,
			&v.IsAdmin,
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

func userReadFilter(opt *Option) (squirrel.And, error) {
	filter := squirrel.And{}

	if opt.Filters != nil {
		for k, v := range opt.Filters {
			switch k {
			case "id":
				filter = append(filter, squirrel.Eq{
					"users.id": v[0],
				})
			case "name":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"users.name": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			case "company":
				f := squirrel.Or{}
				for _, s := range v {
					f = append(f, squirrel.Like{
						"users.company": fmt.Sprintf("%%%s%%", sanitize(s)),
					})
				}
				filter = append(filter, f)
			case "email":
				filter = append(filter, squirrel.Eq{
					"users.email": v[0],
				})
			case "is_admin":
				b, err := strconv.ParseBool(v[0])
				if err != nil {
					return nil, err
				}

				filter = append(filter, squirrel.Eq{
					"users.is_admin": b,
				})
			}
		}
	}

	return filter, nil
}

func (h *Handler) UserCount(ctx context.Context, opt *Option) (uint64, error) {

	filter, err := userReadFilter(opt)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return 0, err
	}

	builder := psql.
		Select("COUNT(*) AS count").
		Distinct().
		From("users").
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
