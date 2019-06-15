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

type BookLabelCreateRequest struct {
	Label string `json:"label"`
}

func (r *BookLabelCreateRequest) Validate() error {
	if len(r.Label) == 0 {
		return obj.InvalidRequest{
			Msg: "label is required",
		}
	}
	return nil
}

type BookLabelCreateResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	BookID    string `json:"book_id"`
	UserID    string `json:"user_id"`
	Label     string `json:"label"`
}

// BookLabelCreate : Create BookLabel
func (s *Server) BookLabelCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	var req *BookLabelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	bookID := c.Param("book_id")

	labelID := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.BookLabelCreate(ctx, &db.BookLabelCreateInput{
			ID:        labelID,
			CreatedAt: ts,
			UpdatedAt: ts,
			BookID:    bookID,
			UserID:    user.ID,
			Label:     req.Label,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &BookLabelCreateResponse{
		ID:        labelID,
		CreatedAt: ts,
		UpdatedAt: ts,
		BookID:    bookID,
		UserID:    user.ID,
		Label:     req.Label,
	})
}
