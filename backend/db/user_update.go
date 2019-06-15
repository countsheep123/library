package db

import (
	"context"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

type UserUpdateInput struct {
	UpdatedAt string
	Name      *string
	Email     *string
	Company   *string
	IsAdmin   *bool
}

func (in *UserUpdateInput) Validate() error {
	if len(in.UpdatedAt) == 0 {
		return obj.Internal{
			Msg: "updated_at is required",
		}
	}
	if in.Name != nil && len(*in.Name) == 0 {
		return obj.Internal{
			Msg: "name is required",
		}
	}
	if in.Email != nil && len(*in.Email) == 0 {
		return obj.Internal{
			Msg: "email is required",
		}
	}
	if in.Company != nil && len(*in.Company) == 0 {
		return obj.Internal{
			Msg: "company is required",
		}
	}
	return nil
}

func (h *Handler) UserUpdate(ctx context.Context, in *UserUpdateInput, filters map[string]string) error {

	if err := in.Validate(); err != nil {
		return err
	}

	kv := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}

	if in.Name != nil {
		kv["name"] = *in.Name
	}
	if in.Email != nil {
		kv["email"] = *in.Email
	}
	if in.Company != nil {
		kv["company"] = *in.Company
	}
	if in.IsAdmin != nil {
		kv["is_admin"] = *in.IsAdmin
	}

	if err := h.Update(ctx, "users", kv, filters); err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	return nil
}
