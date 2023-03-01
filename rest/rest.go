package rest

import (
	"sync"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const port = ":4306"

func Start(wg *sync.WaitGroup) {
	r := gin.Default()
	gin.ForceConsoleColor()

	//CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		MaxAge:       12 * time.Hour,
	}))
	r.Use(gin.Recovery())
	memoryStore := persist.NewMemoryStore(1 * time.Minute)
	r.Use(cache.CacheByRequestPath(memoryStore, 10*time.Second))

	r.GET("/", Root)
	r.GET("/uatr", UatrRank)
	r.GET("/uhar", UHarRank)
	r.GET("/ucor", UCorRank)
	r.GET("/uord", UOrdRank)
	r.GET("/scor", SCorRank)
	r.GET("/sord", SOrdRank)
	r.GET("/shar", SHarRank)
	r.GET("/satr", SAtrRank)
	r.GET("/account/:address", UserReward)
	r.Run(port)
	defer wg.Done()
}
