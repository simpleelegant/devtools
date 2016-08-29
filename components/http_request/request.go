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
	RawQuery    string
}

type response struct {
	StatusCode int
	Status     string
	Header     http.Header
	Body       string
}

func newRequest(c *gin.Context) (*request, error) {
	url_, err := prepareURL(c.PostForm("url"))
	if err != nil {
		return nil, err
	}
	r := &request{
		URL:         url_,
		Method:      c.PostForm("method"),
		ContentType: c.PostForm("contentType"),
		Header:      http.Header{},
		RawQuery:    strings.TrimSpace(c.PostForm("query")),
	}

	if r.Method == http.MethodGet && r.ContentType != urlencodedContentType {
		return nil, errors.New("Method: GET 时，Content-Type 必须是 " + urlencodedContentType)
	}

	// headers
	r.Header.Add("Content-Type", r.ContentType)
	err = onMultlineKeyValue(c.PostForm("header"), ":", func(k, v string) error {
		r.Header.Add(k, v)

		return nil
	})
	if err != nil {
		return nil, errors.New("Header " + err.Error())
	}

	return r, nil
}
func (r *request) getQuery() url.Values {
	values := url.Values{}
	onMultlineKeyValue(r.RawQuery, "=", func(k, v string) error {
		values.Add(k, v)

		return nil
	})

	return values
}

func (r *request) Do() (*response, error) {
	u := r.URL
	body := new(bytes.Buffer)

	switch r.Method {
	case http.MethodGet:
		p := r.getQuery().Encode()
		if p == "" {
			break
		}

		if !strings.Contains(u, "?") {
			u += "?" + p
		} else if strings.HasSuffix(u, "?") || strings.HasSuffix(u, "&") {
			u += p
		} else {
			u += "&" + p
		}
	default:
		switch r.ContentType {
		case urlencodedContentType:
			body.WriteString(r.getQuery().Encode())
		case multipartContentType:
			w := multipart.NewWriter(body)
			defer w.Close()

			err := onMultlineKeyValue(r.RawQuery, "=", func(k, v string) error {
				fw, err := w.CreateFormField(k)
				if err != nil {
					return err
				}
				fw.Write([]byte(v))

				return nil
			})
			if err != nil {
				return nil, err
			}
		default:
			body.WriteString(r.RawQuery)
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Header:     resp.Header,
		Body:       string(result),
	}

	// if result is json string
	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		var b bytes.Buffer
		if err := json.Indent(&b, result, "", "    "); err == nil {
			out.Body = b.String()
		}
	}

	return out, nil
}

func onMultlineKeyValue(s, sep string, do func(k, v string) error) error {
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

func prepareURL(s string) (string, error) {
	s = strings.TrimSpace(s)

	if s == "" {
		return "", errors.New("URL empty not allow")
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", errors.New("URL invalid")
	}
	if u.Scheme == "" {
		s = "http://" + s

		// re-parse
		u, err = url.Parse(s)
		if err != nil {
			return "", errors.New("URL invalid")
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New(`URL must be start with "http://" or "https://"`)
	}

	return s, nil
}
