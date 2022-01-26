package route

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"wms_slave/route/home"
	"wms_slave/route/stock"
	"wms_slave/server"
	"wms_slave/server/middleware"
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
	excel := router.Group("/v1/excel")
	excel.Use(middleware.ExcelHeader())
	{
		excel.GET("/stock/list", stock.ListExcelDownload)
	}

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
