package dataconvert

import (
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
		input := strings.TrimSpace(c.PostForm("input"))

		var (
			output string
			err    error
		)

		switch ability {
		case "jsonIndent":
			output, err = jsonIndent(input)
		case "jsonCompact":
			output, err = jsonCompact(input)
		case "jsonToGoStruct":
			output, err = jsonToGoStruct(input)
		case "jsonToYAML":
			output, err = jsonToYAML(input)
		case "keyValueToJSON":
			output, err = keyValueToJSON(input)
		case "keyValueToQueryString":
			output, err = keyValueToQueryString(input)
		case "queryStringToKeyValue":
			output, err = queryStringToKeyValue(input)
		case "base64URLEncode":
			output, err = base64URLEncode(input)
		case "base64URLDecode":
			output, err = base64URLDecode(input)
		case "md5Checksum":
			output, err = md5Checksum(input)
		case "markdownToHTML":
			output, err = markdownToHTML(input)
		case "escapeNewline":
			output, err = escapeNewline(input)
		case "captureNewline":
			output, err = captureNewline(input)
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
