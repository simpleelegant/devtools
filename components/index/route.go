package index

import "github.com/gin-gonic/gin"

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.File("./components/index/index.html")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./components/index/favicon.ico")
	})
	r.GET("/pure-min.css", func(c *gin.Context) {
		c.File("./components/index/pure-min.css")
	})
	r.GET("/jquery.min.js", func(c *gin.Context) {
		c.File("./components/index/jquery.min.js")
	})
}
