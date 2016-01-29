package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/index/index.html")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/index/favicon.ico")
	})
	r.GET("/pure-min.css", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/index/pure-min.css")
	})
	r.GET("/jquery.min.js", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/index/jquery.min.js")
	})
}
