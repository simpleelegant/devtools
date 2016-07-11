package httprequest

import (
	"bytes"
	"encoding/json"
	"errors"
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

type request struct {
	URL         string
	Method      string
	ContentType string
	Header      http.Header
	Body        *bytes.Buffer
}

type response struct {
	StatusCode int
	Status     string
	Header     http.Header
	Body       string
}

func newRequest(c *gin.Context) (*request, error) {
	body := c.PostForm("query")

	r := &request{
		URL:         c.PostForm("url"),
		Method:      c.PostForm("method"),
		ContentType: c.PostForm("contentType"),
		Header:      make(http.Header),
		Body:        &bytes.Buffer{},
	}

	if strings.TrimSpace(r.URL) == "" {
		return nil, errors.New("URL invalid")
	}

	if _, err := url.Parse(r.URL); err != nil {
		return nil, errors.New("URL 格式错误")
	}

	if r.Method == http.MethodGet {
		if r.ContentType != urlencodedContentType {
			return nil, errors.New("Method: GET 时，Content-Type 必须是 " + urlencodedContentType)
		}

		b, err := parseToQueryString(body)
		if err != nil {
			return nil, err
		}

		if !strings.Contains(r.URL, "?") {
			r.URL += "?" + b
		} else if strings.HasSuffix(r.URL, "?") {
			r.URL += b
		} else {
			r.URL += "&" + b
		}
	}

	// headers
	r.Header.Add("Content-Type", r.ContentType)
	err := onMultlineKV(c.PostForm("header"), ":", func(k, v string) error {
		r.Header.Add(k, v)

		return nil
	})
	if err != nil {
		return nil, errors.New("Header " + err.Error())
	}

	// body
	if r.Method != http.MethodGet {
		if err := parseToBody(r.ContentType, body, r.Body); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *request) Do() (*response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Header:     resp.Header,
		Body:       string(body),
	}

	// if body is json string
	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		var b bytes.Buffer
		if err := json.Indent(&b, body, "", "    "); err == nil {
			out.Body = b.String()
		}
	}

	return out, nil
}

func onMultlineKV(s, sep string, do func(k, v string) error) error {
	for _, r := range strings.Split(strings.TrimSpace(s), "\n") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		kv := strings.SplitN(r, sep, 2)
		if len(kv) != 2 {
			return errors.New("格式错误")
		}
		do(kv[0], kv[1])
	}

	return nil
}

func parseToQueryString(s string) (string, error) {
	values := url.Values{}
	err := onMultlineKV(s, "=", func(k, v string) error {
		values.Add(k, v)

		return nil
	})
	if err != nil {
		return "", err
	}

	return values.Encode(), nil
}

func parseToBody(contentType, s string, b *bytes.Buffer) error {
	switch contentType {
	case urlencodedContentType:
		q, err := parseToQueryString(s)
		if err != nil {
			return err
		}
		b.WriteString(q)
	case multipartContentType:
		w := multipart.NewWriter(b)
		defer w.Close()

		err := onMultlineKV(s, "=", func(k, v string) error {
			fw, err := w.CreateFormField(k)
			if err != nil {
				return err
			}
			fw.Write([]byte(v))

			return nil
		})
		if err != nil {
			return err
		}
	default:
		b.WriteString(s)
	}

	return nil
}

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/http_request/", func(c *gin.Context) {
		c.File("./components/http_request/index.html")
	})

	r.POST("/http_request/request", func(c *gin.Context) {
		req, err := newRequest(c)
		if err != nil {
			responseError(c, err.Error())
			return
		}

		resp, err := req.Do()
		if err != nil {
			responseError(c, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"Error": "",
			"Data":  resp,
		})
	})
}

func responseError(c *gin.Context, err string) {
	c.JSON(http.StatusOK, map[string]string{
		"Error": err,
	})
}
