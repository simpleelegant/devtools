package dataconvert

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/simpleelegant/devtools/components/data_convert/json.syntax.analyser"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/data_convert/", func(c *gin.Context) {
		c.File("./components/data_convert/index.html")
	})

	r.POST("/data_convert/convert", func(c *gin.Context) {
		ability := c.PostForm("ability")
		input := strings.TrimSpace(c.PostForm("input"))

		var output string
		var err error

		switch ability {
		case "jsonIndent":
			output, err = jsonIndent(input)
		case "jsonCompact":
			output, err = jsonCompact(input)
		case "base64Encode":
			output, err = base64Encode(input)
		case "base64Decode":
			output, err = base64Decode(input)
		case "md5Checksum":
			output, err = md5Checksum(input)
		case "jsonToGoStruct":
			output, err = jsonToGoStruct(input)
		case "keyValueToJSON":
			output, err = keyValueToJSON(input)
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

func jsonIndent(input string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(input), "", "    ")

	return out.String(), err
}

func jsonCompact(input string) (string, error) {
	var out bytes.Buffer
	err := json.Compact(&out, []byte(input))

	return out.String(), err
}

func base64Encode(input string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	return encoded, nil
}

func base64Decode(input string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)

	return string(decoded), err
}

func md5Checksum(input string) (string, error) {
	return fmt.Sprintf("%x", md5.Sum([]byte(input))), nil
}

func jsonToGoStruct(input string) (string, error) {
	a := &analyser.LexemeList{}
	if _, err := a.Write([]byte(input)); err != nil {
		return "", err
	}

	t := analyser.SyntaxTree{}
	if err := t.Write(a); err != nil {
		return "", err
	}

	return string(t.EncodeToStruct()), nil
}

func keyValueToJSON(input string) (string, error) {
	e := errors.New("Bad input")
	ctn := make(map[string]string, 10)

	for _, line := range strings.Split(input, "\n") {
		// skip empty line
		if strings.TrimSpace(line) == "" {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			return "", e
		}

		kv[0] = strings.TrimSpace(kv[0])
		if kv[0] == "" {
			return "", e
		}

		ctn[kv[0]] = strings.TrimSpace(kv[1])
	}

	b, err := json.Marshal(ctn)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
