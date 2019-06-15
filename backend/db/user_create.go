package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type UserCreateInput struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Name      string
	Company   string
	Email     string
	IsAdmin   bool
}

func (in *UserCreateInput) Validate() error {
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
	if len(in.Name) == 0 {
		return obj.Internal{
			Msg: "name is required",
		}
	}
	if len(in.Email) == 0 {
		return obj.Internal{
			Msg: "email is required",
		}
	}
	return nil
}

func (h *Handler) UserCreate(ctx context.Context, in *UserCreateInput) error {

	if err := in.Validate(); err != nil {
		return err
	}

	columns := []string{
		"id",
		"created_at",
		"updated_at",
		"name",
		"company",
		"email",
		"is_admin",
	}
	values := []interface{}{
		in.ID,
		in.CreatedAt,
		in.UpdatedAt,
		in.Name,
		in.Company,
		in.Email,
		in.IsAdmin,
	}

	if err := h.Insert(ctx, "users", columns, values); err != nil {
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
