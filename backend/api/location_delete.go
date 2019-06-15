package api

import (
	"net/http"

	"github.com/countsheep123/library/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LocationDelete : Delete Location
func (s *Server) LocationDelete(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	locationID := c.Param("location_id")

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		filters := map[string]string{
			"id": locationID,
		}
		if err := tx.LocationDelete(ctx, filters); err != nil {
			return err
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
