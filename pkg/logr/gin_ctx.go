package logr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type CtxRequestMeta struct {
	Token          string
	Method         string
	URI            string
	QueryString    string
	Body           json.RawMessage
	AppVersion     string
	AcceptLanguage string
}

func ExtractReqMeta(c *gin.Context) *CtxRequestMeta {
	meta := CtxRequestMeta{}
	meta.Token = c.GetHeader("Bearer")
	meta.AppVersion = c.GetHeader("App-Version")
	meta.AcceptLanguage = c.GetHeader("Accept-Language")
	// when you read the the body buffer, the buffer will gone
	bodyBytes, _ := c.GetRawData()
	if len(bodyBytes) > 0 {
		meta.Body = bodyBytes
	} else {
		meta.Body, _ = json.Marshal(map[string]string{})
	}

	// restore the io.ReadCloser
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	query := c.Request.URL.Query()

	// handle the query string map to one row string
	b := new(bytes.Buffer)

	for qsK, q := range query {
		for _, qsV := range q {
			_, _ = fmt.Fprintf(b, "%s=\"%s\" ", qsK, qsV)
		}
	}
	// handle done...

	meta.Method = c.Request.Method
	meta.URI = c.Request.URL.EscapedPath()
	meta.QueryString = b.String()

	return &meta
}
