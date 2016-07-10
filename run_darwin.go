package main

import "github.com/gin-gonic/gin"

func run(e *gin.Engine, address string) error {
	return e.Run(address)
}
