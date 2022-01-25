package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wms_slave/server/e"
)

func ExcelHeader() gin.HandlerFunc {
	return func(context *gin.Context) {
		domain := map[string]bool{
			"kr01.warehouse.pickby.us":    true,
			"kr02.warehouse.pickby.us":    true,
			"api.wms.warehouse.pickby.us": true,
		}
		// 허용 도메인 인가
		if domain[context.GetHeader("WAREHOUSE_DOMAIN")] {
			context.Next()
		} else {
			//context.AbortWithStatus(http.StatusBadRequest)
			context.AbortWithError(http.StatusBadRequest, e.NewException(e.InvalidRequestHeader))
			return
		}
	}
}
