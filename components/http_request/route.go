package httprequest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/http_request/", func(c *gin.Context) {
		c.File("./components/http_request/index.html")
	})

	r.POST("/http_request/request", func(c *gin.Context) {
		req, err := newRequest(c)
		if err != nil {
			responseError(c, err)
			return
		}

		resp, err := req.Do()
		if err != nil {
			responseError(c, err)
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"Error": "",
			"Data":  resp,
			"Curl":  req.ToCURL(),
		})
	})
}

func responseError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, map[string]string{
		"Error": err.Error(),
	})
}
