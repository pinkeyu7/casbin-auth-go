package logr

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewCtxLogger(c *gin.Context) *zap.Logger {
	logger := NewLogger()

	// TBD body is gone
	// when you read the the body buffer, the buffer will gone
	//bodyBytes, _ := c.GetRawData()

	// restore the io.ReadCloser
	//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	query := c.Request.URL.Query()

	// handle the query string map to one row string
	b := new(bytes.Buffer)
	for qsK, q := range query {
		for _, qsV := range q {
			_, _ = fmt.Fprintf(b, "%s=\"%s\" ", qsK, qsV)
		}
	}
	// handle done...

	method := c.Request.Method
	URI := c.Request.URL.EscapedPath()
	qs := b.String()

	logger = logger.With(
		zap.String("method", method),
		zap.String("URI", URI),
		zap.String("qs", qs),
	)

	return logger
}
