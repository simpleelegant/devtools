package httprequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/data_convert/", func(c *gin.Context) {
		c.File("./components/data_convert/index.html")
	})

	r.POST("/data_convert/convert", func(c *gin.Context) {
		function := c.PostForm("function")
		input := c.PostForm("input")

		var output string
		var err error

		switch function {
		case "jsonBeautify":
			output, err = jsonBeautify(input)
		default:
			err = errors.New("function not supported")
		}

		if err != nil {
			c.JSON(http.StatusOK, map[string]string{
				"Error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"Error":   "",
			"Content": output,
		})
	})
}

func jsonBeautify(input string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(input), "", "    ")

	return out.String(), err
}
