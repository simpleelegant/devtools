package httprequest

import "github.com/gin-gonic/gin"

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/http_request", func(c *gin.Context) {
		c.File("./components/http_request/index.html")
	})
}
