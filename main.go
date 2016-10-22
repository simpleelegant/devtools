package main

import (
	"fmt"
	"strings"

	"github.com/simpleelegant/devtools/plugins/conf"
	"github.com/simpleelegant/devtools/plugins/network"

	c5 "github.com/simpleelegant/devtools/components/data_convert"
	c2 "github.com/simpleelegant/devtools/components/documents_service"
	c3 "github.com/simpleelegant/devtools/components/http_log"
	c4 "github.com/simpleelegant/devtools/components/http_request"
	c6 "github.com/simpleelegant/devtools/components/image_convert"
	c1 "github.com/simpleelegant/devtools/components/index"
	c7 "github.com/simpleelegant/devtools/components/kubernetes_dashboard"

	"github.com/gin-gonic/gin"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	chdir()

	conf.Load("./config.json")

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	{
		c1.Route(r)
		c2.Route(r)
		c3.Route(r)
		c4.Route(r)
		c5.Route(r)
		c6.Route(r)
		c7.Route(r)
	}

	addr := fmt.Sprintf("0.0.0.0:%v", conf.Options.Port)

	// bootstrap
	{
		welcome()

		if ips, err := network.GetLocalIPs(); err == nil {
			fmt.Printf("\nServer IP: %s\n\n", strings.Join(ips, " , "))
		}

		fmt.Printf("Listening and serving HTTP on %s\n", addr)

		if conf.Options.AutoOpenBrowser {
			open.Start(fmt.Sprintf("http://127.0.0.1:%v", conf.Options.Port))
		}
	}

	// Start the server
	if err := run(r, addr); err != nil {
		panic(err)
	}
}

func welcome() {
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│              DevTools               │")
	fmt.Println("└─────────────────────────────────────┘")
}
