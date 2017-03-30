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

	r.GET("/kubernetes_dashboard/list_jobs", wrap(listJobs))
	r.GET("/kubernetes_dashboard/describe_job", wrap(describeJob))
	r.DELETE("/kubernetes_dashboard/delete_job", wrap(deleteJob))
	r.GET("/kubernetes_dashboard/list_pods", wrap(listPods))
	r.GET("/kubernetes_dashboard/describe_pod", wrap(describePod))
	r.DELETE("/kubernetes_dashboard/delete_pod", wrap(deletePod))
}

type apiFunc func(*gin.Context) (int, interface{})

func wrap(f apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, res := f(c)

		switch z := res.(type) {
		case error:
			c.JSON(status, map[string]string{"error": z.Error()})
		default:
			c.JSON(status, res)
		}
	}
}
