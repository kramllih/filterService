package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/kramllih/filterService/internal/logger"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logger.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		reqId, _ := c.Get("requestId")

		logger.WithFields(logrus.Fields{"statusCode": statusCode,
			"latencyTime":   latencyTime,
			"clientIP":      clientIP,
			"requestMethod": reqMethod,
			"RequestUri":    reqUri,
			"RequestID":     reqId}).Infof("| %3d | %13v | %15s | %s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			reqId,
		)
	}
}

func ErrorHandler(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		reqId, _ := c.Get("requestId")
		for _, ginErr := range c.Errors {
			logger.WithField("RequestID", reqId).Error(ginErr)
		}
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestID, err := uuid.NewV4()
		if err != nil {
			c.Error(fmt.Errorf("error generating new requestID: %w", err))
			c.Next()
		}
		c.Set("requestId", xRequestID)
		c.Next()
	}
}
