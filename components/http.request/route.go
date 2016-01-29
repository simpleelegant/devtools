package httprequest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/http.request", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/http.request/index.html")
	})
}
