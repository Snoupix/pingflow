package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeForever(is_dev bool) {
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

    log.Fatal(server.Run())
}
