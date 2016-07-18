package dataconvert

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/data_convert/", func(c *gin.Context) {
		c.File("./components/data_convert/index.html")
	})

	r.POST("/data_convert/convert", func(c *gin.Context) {
		ability := c.PostForm("ability")
		input := c.PostForm("input")

		var output string
		var err error

		switch ability {
		case "jsonBeautify":
			output, err = jsonBeautify(input)
		case "base64Encode":
			output, err = base64Encode(input)
		default:
			err = errors.New("Not supported")
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

func base64Encode(input string) (string, error) {
	input = strings.TrimSpace(input)
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	return encoded, nil
}
