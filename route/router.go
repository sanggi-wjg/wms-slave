package route

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"wms_slave/route/home"
	"wms_slave/route/stock"
	"wms_slave/server"
)

func Init() *gin.Engine {
	router := gin.New()

	// Initialize
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Load Balancer 설정
	router.RemoteIPHeaders = []string{"X-Forwarded-For"}
	if err := router.SetTrustedProxies(server.Config.TrustProxies); err != nil {
		panic(err)
	}

	// Default
	router.GET("/", home.Hello)
	router.GET("/ping", home.Ping)
	router.GET("/panic", home.PanicTest)

	// Excel
	router.GET("/stock_list", stock.ListToExcel)

	router.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)
			// note that you are using the copied context "cCp", IMPORTANT
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	router.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	return router
}
