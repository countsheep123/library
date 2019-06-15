package api

import (
	"net/http"

	"github.com/countsheep123/library/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BookDelete : Delete Book
func (s *Server) BookDelete(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	id := c.Param("book_id")

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		filters := map[string]string{
			"id": id,
		}
		if err := tx.BookDelete(ctx, filters); err != nil {
			return err
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
