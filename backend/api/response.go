package api

import (
	"net/http"

	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func boolValue(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func strValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func strArrayValue(a []string) []string {
	if a == nil {
		return []string{}
	}
	return a
}

func errResponse(c *gin.Context, err error) {
	switch err.(type) {
	case obj.InvalidRequest:
		zap.S().Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	case obj.NotFound:
		zap.S().Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case obj.Duplicate:
		zap.S().Warn(err)
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate"})
	case obj.Internal:
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		c.Status(http.StatusInternalServerError)
	default:
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		c.Status(http.StatusInternalServerError)
	}
}
