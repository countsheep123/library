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

type MarkCreateRequest struct {
	Name string  `json:"name"`
	URL  *string `json:"url"`
}

func (r *MarkCreateRequest) Validate() error {
	if len(r.Name) == 0 {
		return obj.InvalidRequest{
			Msg: "name is required",
		}
	}
	return nil
}

type MarkCreateResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

// MarkCreate : Create Mark
func (s *Server) MarkCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	var req *MarkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	markID := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.MarkCreate(ctx, &db.MarkCreateInput{
			ID:        markID,
			CreatedAt: ts,
			UpdatedAt: ts,
			UserID:    user.ID,
			Name:      req.Name,
			URL:       req.URL,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &MarkCreateResponse{
		ID:        markID,
		CreatedAt: ts,
		UpdatedAt: ts,
		UserID:    user.ID,
		Name:      req.Name,
		URL:       strValue(req.URL),
	})
}
