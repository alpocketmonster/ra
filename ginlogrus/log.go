package ginlogrus

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// 2016-09-27 09:38:21.541541811 +0200 CEST
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700]
// "GET /apache_pb.gif HTTP/1.0" 200 2326
// "http://www.example.com/start.html"
// "Mozilla/4.08 [en] (Win98; I ;Nav)"

var timeFormat = "02/Jan/2006:15:04:05 -0700"

// Logger is the logrus logger handler
func Logger(notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {

		//log := zerolog.New(os.Stderr).With().Timestamp().Logger()

		log.Debug().Str("foo", "bar").Msg("Hello World")

		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := float64(stop.Nanoseconds()) / 1000000.0
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		// entry := log.WithFields(log.Fields{
		// 	"hostname": hostname,
		// 	"statusCode": statusCode,
		// 	"latency":    latency, // time to process
		// 	"clientIP":   clientIP,
		// 	"method": c.Request.Method,
		// 	"path":       path,
		// 	"referer":    referer,
		// 	"dataLength": dataLength,
		// 	"userAgent":  clientUserAgent,
		// })

		entry := log.With().
			Str("component", "foo").
			Str("hostname", hostname).
			Int("statusCode", statusCode).
			Float64("latency", latency). // time to process
			Str("clientIP", clientIP).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("referer", referer).
			Int("dataLength", dataLength).
			Str("userAgent", clientUserAgent).
			Logger()

		if len(c.Errors) > 0 {
			entry.Error().Msg(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			//msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				entry.Error().Send()
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn().Send()
			} else {
				entry.Info().Send()
			}
		}
	}
}
