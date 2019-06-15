package api

import (
	"net/http"

	"github.com/countsheep123/library/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BookLabelDelete : Delete BookLabel
func (s *Server) BookLabelDelete(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	bookID := c.Param("book_id")
	labelID := c.Param("label_id")

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		filters := map[string]string{
			"id":      labelID,
			"book_id": bookID,
		}
		if err := tx.BookLabelDelete(ctx, filters); err != nil {
			return err
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
