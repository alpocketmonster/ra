package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// get prometheus handler for server
func (m *metrics) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()

		m.monitor.GetMetric("requests_total").Inc([]string{method, path})
		m.monitor.GetMetric("response_status").Inc([]string{strconv.Itoa(statusCode)})
	}
}
