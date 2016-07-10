package main

import (
	"fmt"
	"strings"

	"github.com/simpleelegant/devtools/components/index"
	"github.com/simpleelegant/devtools/plugins/conf"
	"github.com/simpleelegant/devtools/plugins/network"

	a "github.com/simpleelegant/devtools/components/documents_service"
	b "github.com/simpleelegant/devtools/components/http_log"
	c "github.com/simpleelegant/devtools/components/http_request"

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
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

func welcome() {
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│              DevTools               │")
	fmt.Println("└─────────────────────────────────────┘")
}
