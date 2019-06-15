package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/countsheep123/library/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserListResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	IsAdmin   bool   `json:"is_admin"`
}

// UserList : List User
func (s *Server) UserList(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	opt := userListOpt(c, user)

	res := []*UserListResponse{}
	var total uint64

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		var users []*db.UserReadOutput
		if err := tx.UserRead(ctx, opt, &users); err != nil {
			return err
		}

		count, err := tx.UserCount(ctx, opt)
		if err != nil {
			return err
		}
		total = count

		for _, u := range users {
			res = append(res, &UserListResponse{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
				Email:     u.Email,
				Name:      u.Name,
				Company:   strValue(u.Company),
				IsAdmin:   u.IsAdmin,
			})
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"users": res,
	})
}

func userListOpt(c *gin.Context, user *User) *db.Option {
	opt := &db.Option{
		Filters: map[string][]string{},
		Sort:    &db.Sort{},
		Range:   &db.Range{},
	}

	if vs, ok := c.GetQueryArray("name"); ok {
		opt.Filters["name"] = vs
	}
	if vs, ok := c.GetQueryArray("email"); ok {
		opt.Filters["email"] = vs
	}

	sort := c.DefaultQuery("sort", "name:asc")
	splited := strings.Split(sort, ":")
	if len(splited) != 2 {
		opt.Sort.By = "name"
		opt.Sort.IsAsc = true
	} else {
		switch splited[0] {
		case "created_at", "name", "email":
			opt.Sort.By = splited[0]
		default:
			opt.Sort.By = "name"
		}

		switch splited[1] {
		case "ASC", "asc":
			opt.Sort.IsAsc = true
		case "DESC", "desc":
			opt.Sort.IsAsc = false
		default:
			opt.Sort.IsAsc = true
		}
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err == nil {
		opt.Range.Limit = limit
	}
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err == nil {
		opt.Range.Offset = offset
	}

	return opt
}
