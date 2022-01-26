package stock

import (
	"github.com/gin-gonic/gin"
	"wms_slave/server/global"
	"wms_slave/server/logger"
)

func ListExcelDownload(context *gin.Context) {
	param := makeParameter(context)
	excel := NewStockListExcel(param)
	filename := excel.Make()

	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", "attachment; filename="+filename)
	context.Header("Content-Type", "application/vnd.ms-excel")
	context.Header("Expires", "0")
	context.Header("Connection", "close")
	context.File(filename)
}

type parameter struct {
	WarehouseDomain   string
	WarehouseId       string
	FromDate          string
	ToDate            string
	PartnerId         string
	ProductOwner      string
	ProductVendorName string
	StockType         string
	PartnerUserType   string
	TransferCompany   string
	Code              string
	RackCd            string
	ProductName       string
	ProductOption     string
	ProductBrandName  string
	InWaybillNo       string
	InOrderCd         string

	ProductGroupCLOTH     bool
	ProductGroupACCESSORY bool
	ProductGroupBAGSHOES  bool
	ProductGroupCOSMETIC  bool
	PartnerExcelFlag      bool
}

func makeParameter(c *gin.Context) parameter {
	// http://localhost:9000/v1/excel/stock/list?partnerId=jamy&fromDate=2020-01-01&toDate=2022-01-01
	//p := map[string]interface{}{}
	//p["warehouseDomain"] = c.GetHeader("WAREHOUSE_DOMAIN")
	//p["warehouseId"] = global.GetWarehouseIdByDomain(c.GetHeader("WAREHOUSE_DOMAIN"))
	//p["fromDate"] = c.DefaultQuery("fromDate", "")
	//p["toDate"] = c.DefaultQuery("toDate", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")
	//p["partnerId"] = c.DefaultQuery("partnerId", "")

	p := parameter{}
	p.WarehouseDomain = c.GetHeader("WAREHOUSE_DOMAIN")
	p.WarehouseId = global.GetWarehouseIdByDomain(p.WarehouseDomain)
	p.FromDate = c.DefaultQuery("fromDate", "")
	p.ToDate = c.DefaultQuery("toDate", "")
	p.PartnerId = c.DefaultQuery("partnerId", "")
	p.PartnerUserType = c.DefaultQuery("partnerUserType", "")
	p.ProductOwner = c.DefaultQuery("productOwner", "")
	p.ProductVendorName = c.DefaultQuery("productVendorName", "")
	p.StockType = c.DefaultQuery("stockType", "")
	p.PartnerUserType = c.DefaultQuery("transferCompany", "")
	p.TransferCompany = c.DefaultQuery("code", "")
	p.Code = c.DefaultQuery("rackCd", "")
	p.RackCd = c.DefaultQuery("productName", "")
	p.ProductName = c.DefaultQuery("productOption", "")
	p.ProductOption = c.DefaultQuery("rackCd", "")
	p.ProductBrandName = c.DefaultQuery("productBrandName", "")
	p.InWaybillNo = c.DefaultQuery("inWaybillNo", "")
	p.InOrderCd = c.DefaultQuery("inOrderCd", "")

	p.PartnerExcelFlag = convertGetQueryStringValueToBoolean(c.DefaultQuery("partnerExcelFlag", ""))
	p.ProductGroupCLOTH = convertGetQueryStringValueToBoolean(c.DefaultQuery("productGroup_CLOTHING", ""))
	p.ProductGroupACCESSORY = convertGetQueryStringValueToBoolean(c.DefaultQuery("productGroup_ACCESSORY", ""))
	p.ProductGroupBAGSHOES = convertGetQueryStringValueToBoolean(c.DefaultQuery("productGroup_BAGSHOES", ""))
	p.ProductGroupCOSMETIC = convertGetQueryStringValueToBoolean(c.DefaultQuery("productGroup_COSMETICS", ""))
	defer func() {
		logger.Log.Debugf("%+v\n", p)
	}()

	return p
}

func convertGetQueryStringValueToBoolean(queryValue string) bool {
	if queryValue == "Y" || queryValue == "T" {
		return true
	}
	return false
}
