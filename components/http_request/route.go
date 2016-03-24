package httprequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// supported content-type
const (
	urlencodedContentType = "application/x-www-form-urlencoded"
	multipartContentType  = "multipart/form-data"
	jsonContentType       = "application/json"
	xmlContentType        = "application/xml"
)

// supported methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/http_request/", func(c *gin.Context) {
		c.File("./components/http_request/index.html")
	})

	r.POST("/http_request/request", func(c *gin.Context) {
		ur := c.PostForm("url")
		method := c.PostForm("method")
		contentType := c.PostForm("contentType")
		query := strings.TrimSpace(c.PostForm("query"))
		buffer := &bytes.Buffer{}

		if _, err := url.Parse(ur); err != nil {
			responseJSON(c, "错误：请输入正确的URL", "")
			return
		}

		// precess query
		switch method {
		case GET:
			if contentType != urlencodedContentType {
				responseJSON(c, "错误：Method: GET 时，Content-Type 必须是 "+urlencodedContentType, "")
				return
			}

			if query != "" {
				if _, err := url.ParseQuery(query); err != nil {
					responseJSON(c, "Query 格式错误", "")
					return
				}

				ur += "?" + query
			}
		case POST, PUT, DELETE:
			switch contentType {
			case urlencodedContentType:
				buffer.WriteString(query)
			case multipartContentType:
				writeMultipartForm(query, buffer)
			case jsonContentType, xmlContentType:
				buffer.WriteString(query)
			default:
				responseJSON(c, "错误：不支持的 Content-Type", "")
				return
			}
		default:
			responseJSON(c, "错误：不支持的 Method", "")
			return
		}

		// request
		{
			req, err := http.NewRequest(method, ur, buffer)
			if err != nil {
				responseJSON(c, "错误：创建请求失败", "")
				return
			}
			req.Header.Set("Content-Type", contentType)
			writeHeaders(c.PostForm("headers"), req)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				responseJSON(c, "错误：无法访问你指定的URL", "")
				return
			}

			// parse response
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				responseJSON(c, resp.Status, "（读取响应内容失败）")
				return
			}

			responseJSON(c, resp.Status, string(body))
		}
	})
}

func writeHeaders(headers string, req *http.Request) {
	var vc map[string]interface{}

	json.Unmarshal([]byte(headers), &vc)
	for k, v := range vc {
		req.Header.Add(k, fmt.Sprintf("%v", v))
	}
}

func writeMultipartForm(data string, buffer *bytes.Buffer) error {
	w := multipart.NewWriter(buffer)
	defer w.Close()

	var vc map[string]interface{}
	json.Unmarshal([]byte(data), &vc)
	for k, v := range vc {
		fw, err := w.CreateFormField(k)
		if err != nil {
			return err
		}

		fw.Write([]byte(fmt.Sprintf("%v", v)))
	}

	return nil
}

func responseJSON(c *gin.Context, err, content string) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Error":   err,
		"Content": content,
	})
}
