package main

import (
	"yujian/devtools/components/index"

	"yujian/devtools/components/http.log"
	"yujian/devtools/components/http.request"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// register routes
	{
		index.Route(r)
		httplog.Route(r)
		httprequest.Route(r)
	}

	addr := "0.0.0.0:9900"

	println("Listening and serving HTTP on " + addr)

	// Start the server
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
