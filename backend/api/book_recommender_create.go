package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type BookRecommenderCreateResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	BookID    string `json:"book_id"`
	UserID    string `json:"user_id"`
}

// BookRecommenderCreate : Create BookRecommender
func (s *Server) BookRecommenderCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	bookID := c.Param("book_id")

	recommenderID := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.BookRecommenderCreate(ctx, &db.BookRecommenderCreateInput{
			ID:        recommenderID,
			CreatedAt: ts,
			UpdatedAt: ts,
			BookID:    bookID,
			UserID:    user.ID,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &BookRecommenderCreateResponse{
		ID:        recommenderID,
		CreatedAt: ts,
		UpdatedAt: ts,
		BookID:    bookID,
		UserID:    user.ID,
	})
}
