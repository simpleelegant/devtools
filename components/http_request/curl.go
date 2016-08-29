package httprequest

import (
	"fmt"
	"strings"
)

// ToCURL convert to CURL figure
func (r *request) ToCURL() string {
	var (
		header string
		data   string
	)

	for k, h := range r.Header {
		for _, v := range h {
			header += fmt.Sprintf("  -H \"%s: %s\" \\\n", k, v)
		}
	}

	switch r.ContentType {
	case urlencodedContentType, multipartContentType:
		for k, vl := range r.getQuery() {
			for _, v := range vl {
				data += fmt.Sprintf("  -d \"%s=%s\" \\\n", k, v)
			}

		}
	default:
		if r.RawQuery != "" {
			data += fmt.Sprintf("  -d \"%s\" \\\n", strings.Replace(r.RawQuery, "\"", "\\\"", -1))
		}
	}

	return fmt.Sprintf("curl -i -X %s \\\n%s%s  \"%s\"", r.Method, header, data, r.URL)
}
