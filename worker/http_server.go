package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeForever() {
	server := gin.Default()

	if !is_dev {
		gin.SetMode(gin.ReleaseMode)
	}

	server.GET("/classes", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"values": "[\"barbarian\", \"bard\", \"cleric\", \"druid\", \"fighter\", \"monk\", \"paladin\", \"ranger\", \"rogue\", \"sorcerer\", \"warlock\", \"wizard\"]",
		})
	})

	server.GET("/sub-classes", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"values": "[\"berserker\", \"champion\", \"devotion\", \"draconic\", \"evocation\", \"fiend\", \"hunter\", \"land\", \"life\", \"lore\", \"open-hand\", \"thief\"]",
		})
	})

	server.GET("/job-index", func(c *gin.Context) {
		c.String(http.StatusOK, "%03d", new_job_idx())
	})

	log.Fatal(server.Run(fmt.Sprintf("%s:%s", get_env("WORKER_ADDR"), get_env("WORKER_PORT"))))
}
