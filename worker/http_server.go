package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"worker/utils"
)

func ServeForever() {
	server := gin.Default()

	if !is_dev {
		gin.SetMode(gin.ReleaseMode)
	}

	server.GET("/job-index", func(c *gin.Context) {
		c.String(http.StatusOK, "%03d", utils.NewWorkIdx())
	})

	log.Fatal(server.Run(fmt.Sprintf("%s:%s", utils.GetEnv("WORKER_ADDR"), utils.GetEnv("WORKER_PORT"))))
}
