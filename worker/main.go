package main

import (
	"log"
	"net/http"
	"os"

	// rd => redis driver
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	rd "github.com/redis/go-redis/v9"
)

var is_dev bool
var redis *rd.Client

func init() {
    is_dev = len(os.Args) == 2 && os.Args[1] == "-dev"

    if is_dev {
        if err := godotenv.Load("../.env.public"); err != nil {
            log.Fatalf("Failed to load ../.env.public env file: %s\n", err)
        }
    }

    redis_addr_key := "REDIS_ADDR"
    redis_addr, ok := os.LookupEnv(redis_addr_key)
    if !ok {
        log.Panicf("Error: redis addr not found in env. Key: '%s'\n", redis_addr_key)
    }

    redis = rd.NewClient(&rd.Options{
        Addr: redis_addr,
    })
}

func main() {
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
