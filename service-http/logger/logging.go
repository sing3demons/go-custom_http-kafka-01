package logger

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logrus.SetLevel(logLevel)
	log.SetOutput(logger.Writer())
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return &Logger{logger}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()
		// Processing request
		ctx.Next()
		// End Time request
		endTime := time.Now()
		// Request method
		reqMethod := ctx.Request.Method
		// Request route
		path := ctx.Request.RequestURI
		// status code
		statusCode := ctx.Writer.Status()
		// Request IP
		clientIP := ctx.ClientIP()
		// Request host
		host := ctx.Request.Host
		// Request user agent
		userID, exists := ctx.Get("userId")
		if exists {
			userID = userID.(string)
		} else {
			userID = ""
		}

		reqId := ctx.Writer.Header().Get("X-Request-Id")
		if reqId == "" {
			reqId = "-"
		}
		body_size := ctx.Writer.Size()
		// execution time
		latencyTime := endTime.Sub(startTime)
		logrus.WithFields(logrus.Fields{
			"method":        reqMethod,
			"status":        statusCode,
			"latency":       latencyTime,
			"client_ip":     clientIP,
			"request_id":    reqId,
			"remote_ip":     ctx.Request.RemoteAddr,
			"user_id":       userID,
			"user_agent":    ctx.Request.UserAgent(),
			"error":         ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			"request":       ctx.Request.PostForm.Encode(),
			"body_size":     body_size,
			"host":          host,
			"protocol":      ctx.Request.Proto,
			"path":          path,
			"query":         ctx.Request.URL.RawQuery,
			"response_size": ctx.Writer.Size(),
			"ContentType":   ctx.ContentType(),
			"ContentLength": ctx.Request.ContentLength,
			"timezone":      time.Now().Location().String(),
			"ISOTime":       startTime,
			"UnixTime":      startTime.UnixNano(),
		}).Debug("HTTP::REQUEST")
		ctx.Next()
	}
}
