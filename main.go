// 开发工具套件
package main

import (
	"fmt"
	"yujian/devtools/components/index"
	"yujian/devtools/plugins/conf"

	a "yujian/devtools/components/documents_service"
	b "yujian/devtools/components/http_log"
	c "yujian/devtools/components/http_request"

	"github.com/gin-gonic/gin"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	chdir()
	conf.Load("./config.json")
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	{
		index.Route(r)
		a.Route(r)
		b.Route(r)
		c.Route(r)
	}

	addr := fmt.Sprintf("0.0.0.0:%v", conf.Options.Port)
	fmt.Println("Listening and serving HTTP on ", addr)

	if conf.Options.AutoOpenBrowser {
		open.Start("http://" + addr)
	}

	// Start the server
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
