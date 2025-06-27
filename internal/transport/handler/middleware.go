package handler

import (
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
)

func (h *Handler) TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		guid := uuid.New().String()
		logger.GetLoggerFromCtx(h.Ctx).Info("Request Trace", zap.String("method", c.Request.Method), zap.String("TraceID", guid))
		c.Next()

		status := c.Writer.Status()
		logger.GetLoggerFromCtx(h.Ctx).Info("Response Trace", zap.String("response code", strconv.Itoa(status)), zap.String("TraceID", guid))
	}
}
