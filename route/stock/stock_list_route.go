package stock

import (
	"github.com/gin-gonic/gin"
	"wms_slave/server/logger"
)

func ListToExcel(context *gin.Context) {
	param := makeParamMap(context)
	filename := makeExcelFile(param)

	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", "attachment; filename="+filename)
	context.Header("Content-Type", "application/vnd.ms-excel")
	context.Header("Expires", "0")
	context.Header("Connection", "close")
	context.File(filename)
}

func makeParamMap(c *gin.Context) map[string]interface{} {
	p := map[string]interface{}{}
	defer logger.Log.Debug(p)

	// http://localhost:9000/v1/excel/stock/list?partnerId=jamy&fromDate=2020-01-01&toDate=2022-01-01
	//p["warehouseDomain"] = c.DefaultQuery("warehouseDomain", c.GetHeader("WAREHOUSE_DOMAIN"))
	p["fromDate"] = c.DefaultQuery("fromDate", "")
	p["toDate"] = c.DefaultQuery("toDate", "")
	p["partnerId"] = c.DefaultQuery("partnerId", "")
	p["partnerUserType"] = c.DefaultQuery("partnerUserType", "")
	p["transferCompany"] = c.DefaultQuery("transferCompany", "")
	p["code"] = c.DefaultQuery("code", "")
	p["rackCd"] = c.DefaultQuery("rackCd", "")
	p["productName"] = c.DefaultQuery("productName", "")
	p["productOption"] = c.DefaultQuery("productOption", "")
	p["rackCd"] = c.DefaultQuery("rackCd", "")
	p["productBrandName"] = c.DefaultQuery("productBrandName", "")
	p["inWaybillNo"] = c.DefaultQuery("inWaybillNo", "")
	p["inOrderCd"] = c.DefaultQuery("inOrderCd", "")

	return p
}
