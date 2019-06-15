package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type LocationCreateRequest struct {
	Name string `json:"name"`
}

func (r *LocationCreateRequest) Validate() error {
	if len(r.Name) == 0 {
		return obj.InvalidRequest{
			Msg: "name is required",
		}
	}
	return nil
}

type LocationCreateResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
}

// LocationCreate : Create Location
func (s *Server) LocationCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	var req *LocationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	locationID := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.LocationCreate(ctx, &db.LocationCreateInput{
			ID:        locationID,
			CreatedAt: ts,
			UpdatedAt: ts,
			Name:      req.Name,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &LocationCreateResponse{
		ID:        locationID,
		CreatedAt: ts,
		UpdatedAt: ts,
		Name:      req.Name,
	})
}
