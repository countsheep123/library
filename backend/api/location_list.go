package api

import (
	"net/http"
	"strings"

	"github.com/countsheep123/library/db"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type LocationListResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
}

// LocationList : List Location
func (s *Server) LocationList(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	opt := locationListOpt(c)

	res := []*LocationListResponse{}

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		var locations []*db.LocationReadOutput
		if err := tx.LocationRead(ctx, opt, &locations); err != nil {
			return err
		}

		for _, l := range locations {
			res = append(res, &LocationListResponse{
				ID:        l.ID,
				CreatedAt: l.CreatedAt,
				UpdatedAt: l.UpdatedAt,
				Name:      l.Name,
			})
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func locationListOpt(c *gin.Context) *db.Option {
	opt := &db.Option{
		Filters: map[string][]string{},
		Sort:    &db.Sort{},
		// Range:   &db.Range{},
	}

	if vs, ok := c.GetQueryArray("name"); ok {
		opt.Filters["name"] = vs
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
