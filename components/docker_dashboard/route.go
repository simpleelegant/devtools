package docker_dashboard

import "github.com/gin-gonic/gin"

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/docker_dashboard/", func(c *gin.Context) {
		c.File("./components/docker_dashboard/index.html")
	})
	r.GET("/docker_dashboard/index.js", func(c *gin.Context) {
		c.File("./components/docker_dashboard/index.js")
	})

	r.POST("/docker_dashboard/get_images", getImages)
	r.POST("/docker_dashboard/delete_image", deleteImage)
	r.POST("/docker_dashboard/inspect_image", inspectImage)
}
