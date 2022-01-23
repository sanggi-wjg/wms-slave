package stock

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"wms_slave/server/e"
)

func ListToExcel(context *gin.Context) {
	warehouseDomain := context.GetHeader("WAREHOUSE_DOMAIN")
	switch warehouseDomain {
	case "kr01.warehouse.pickby.us":
	case "kr02.warehouse.pickby.us":
	case "api.wms.warehouse.pickby.us":
		break
	default:
		context.String(http.StatusBadRequest, e.NewException(e.InvalidRequestHeader).Error())
		return
	}

	partnerId := context.DefaultQuery("partnerId", "")
	log.Println("partnerId:" + partnerId)

	filename := makeToStockExcel()
	//context.Header("Content-Description", "File Transfer")
	//context.Header("Content-Transfer-Encoding", "binary")
	//context.Header("Content-Disposition", "attachment; filename="+filename)
	//context.Header("Content-Type", "application/vnd.ms-excel")
	//context.Header("Expires", "0")
	//context.Header("Connection", "close")
	//context.File(filename)
	context.String(http.StatusOK, filename)
}
