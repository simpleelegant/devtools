package kubernetes_dashboard

import "github.com/gin-gonic/gin"

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/kubernetes_dashboard/", func(c *gin.Context) {
		c.File("./components/kubernetes_dashboard/index.html")
	})
	r.GET("/kubernetes_dashboard/index.js", func(c *gin.Context) {
		c.File("./components/kubernetes_dashboard/index.js")
	})

	r.POST("/kubernetes_dashboard/get_jobs", getJobs)
	r.POST("/kubernetes_dashboard/describe_job", describeJob)
	r.POST("/kubernetes_dashboard/delete_job", deleteJob)
	r.POST("/kubernetes_dashboard/describe_pod", describePod)
	r.POST("/kubernetes_dashboard/delete_pod", deletePod)
}
