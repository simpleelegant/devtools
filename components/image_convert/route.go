package imageconvert

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/image_convert/", func(c *gin.Context) {
		c.File("./components/image_convert/index.html")
	})

	r.POST("/image_convert/image2base64", func(c *gin.Context) {
		f, fh, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{
				"Error": err.Error(),
			})

			return
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{
				"Error": err.Error(),
			})

			return
		}

		out := base64.StdEncoding.EncodeToString(b)
		switch c.PostForm("output") {
		case "data_url":
			tmpl := "data:%s;base64,%s"
			out = fmt.Sprintf(tmpl, fh.Header.Get("Content-Type"), out)
		case "img_tag":
			tmpl := `<img src="data:%s;base64,%s" />`
			out = fmt.Sprintf(tmpl, fh.Header.Get("Content-Type"), out)
		default:
		}

		tmpl, err := ioutil.ReadFile("./components/image_convert/index.html")
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{
				"Error": err.Error(),
			})

			return
		}

		c.Data(http.StatusOK, "text/html", []byte(strings.Replace(string(tmpl), "{{ .output }}", out, -1)))
	})
}
