package api

import (
	"net/http"
	"strings"

	"github.com/countsheep123/library/db"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type MarkListResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

// MarkList : List Mark
func (s *Server) MarkList(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	opt := markListOpt(c, user)

	res := []*MarkListResponse{}

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		var marks []*db.MarkReadOutput
		if err := tx.MarkRead(ctx, opt, &marks); err != nil {
			return err
		}

		for _, m := range marks {
			res = append(res, &MarkListResponse{
				ID:        m.ID,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
				OwnerID:   m.UserID,
				OwnerName: m.UserName,
				Name:      m.Name,
				URL:       strValue(m.URL),
			})
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func markListOpt(c *gin.Context, user *User) *db.Option {
	opt := &db.Option{
		Filters: map[string][]string{},
		Sort:    &db.Sort{},
		// Range:   &db.Range{},
	}

	if vs, ok := c.GetQueryArray("name"); ok {
		opt.Filters["name"] = vs
	}

	own := c.DefaultQuery("own", "false")
	switch own {
	case "true":
		opt.Filters["owner_id"] = []string{user.ID}
	case "false":
		// do nothing
	default:
		// do nothing
	}

	sort := c.DefaultQuery("sort", "name:asc")
	splited := strings.Split(sort, ":")
	if len(splited) != 2 {
		opt.Sort.By = "name"
		opt.Sort.IsAsc = true
	} else {
		switch splited[0] {
		case "created_at", "name":
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

	// limitStr := c.DefaultQuery("limit", "100")
	// limit, err := strconv.ParseUint(limitStr, 10, 64)
	// if err == nil {
	// 	opt.Range.Limit = limit
	// }
	// offsetStr := c.DefaultQuery("offset", "0")
	// offset, err := strconv.ParseUint(offsetStr, 10, 64)
	// if err == nil {
	// 	opt.Range.Offset = offset
	// }

	return opt
}
