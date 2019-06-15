package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserUpdateRequest struct {
	Company *string `json:"company"`
	IsAdmin *bool   `json:"is_admin"`
}

func (r *UserUpdateRequest) Validate() error {
	return nil
}

type UserUpdateResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	IsAdmin   bool   `json:"is_admin"`
}

// UserUpdate : Update User
func (s *Server) UserUpdate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	id := c.Param("user_id")

	var req *UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	ts := fmt.Sprint(time.Now().UnixNano())

	res := &UserUpdateResponse{
		ID: id,
	}

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		opt := &db.Option{
			Filters: map[string][]string{
				"id": []string{id},
			},
		}
		var users []*db.UserReadOutput
		if err := tx.UserRead(ctx, opt, &users); err != nil {
			return err
		}

		if len(users) == 0 {
			return obj.NotFound{}
		}

		res.CreatedAt = users[0].CreatedAt
		res.UpdatedAt = ts

		filters := map[string]string{
			"id": id,
		}

		if err := tx.UserUpdate(ctx, &db.UserUpdateInput{
			UpdatedAt: ts,
			Company:   req.Company,
			IsAdmin:   req.IsAdmin,
		}, filters); err != nil {
			return err
		}

		if req.Company != nil {
			res.Company = *req.Company
		} else {
			res.Company = strValue(users[0].Company)
		}

		if req.IsAdmin != nil {
			res.IsAdmin = *req.IsAdmin
		} else {
			res.IsAdmin = users[0].IsAdmin
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
