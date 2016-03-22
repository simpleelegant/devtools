package httplog

import (
	"encoding/json"
	"net/http"
	"yujian/devtools/plugins/websocket"

	"github.com/gin-gonic/gin"
)

// Route register routes
func Route(r *gin.Engine) {
	r.GET("/http_log/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/http_log/index.html")
	})

	r.GET("/http_log/help.html", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./components/http_log/help.html")
	})

	r.GET("/http_log/ws", func(c *gin.Context) {
		websocket.ConnectHandler(c.Writer, c.Request)
	})

	r.POST("/http_log/log", func(c *gin.Context) {
		f := c.Request

		client := f.FormValue("client")
		if client == "" {
			client = "Unkown"
		}

		data, _ := json.Marshal(map[string]string{
			"client":       client,
			"timeStamp":    f.FormValue("timeStamp"),
			"method":       f.FormValue("method"),
			"url":          f.FormValue("url"),
			"params":       f.FormValue("params"),
			"statusCode":   f.FormValue("statusCode"),
			"responseBody": f.FormValue("responseBody"),
		})

		// broadcast
		websocket.Server.Broadcast(data)

		c.Writer.Write([]byte("success"))
	})
}
