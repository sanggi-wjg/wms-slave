package stock

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"wms_slave/model"
	"wms_slave/server/global"
	"wms_slave/server/logger"
	"wms_slave/server/route_util"
)

func ListExcelDownload(context *gin.Context) {
	excel := newStockListExcel(makeParameter(context))
	filename := excel.Make()

	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", "attachment; filename="+filename)
	context.Header("Content-Type", "application/vnd.ms-excel")
	context.Header("Expires", "0")
	context.Header("Connection", "close")
	context.File(filename)
}

/*type parameter struct {
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
}*/

func makeParameter(c *gin.Context) map[string]interface{} {
	p := map[string]interface{}{}
	p["warehouseDomain"] = c.GetHeader("WAREHOUSE_DOMAIN")
	p["warehouseId"] = global.GetWarehouseIdByDomain(p["warehouseDomain"].(string))
	p["fromDate"] = c.DefaultQuery("fromDate", "")
	p["toDate"] = c.DefaultQuery("toDate", "")
	p["partnerId"] = c.DefaultQuery("partnerId", "")
	p["partnerUserType"] = c.DefaultQuery("partnerUserType", "")
	p["productOwner"] = c.DefaultQuery("productOwner", "")
	p["productVendorName"] = c.DefaultQuery("productVendorName", "")
	p["stockType"] = c.DefaultQuery("stockType", "")
	p["partnerUserType"] = c.DefaultQuery("transferCompany", "")
	p["transferCompany"] = c.DefaultQuery("code", "")
	p["code"] = c.DefaultQuery("rackCd", "")
	p["rackCd"] = c.DefaultQuery("productName", "")
	p["productName"] = c.DefaultQuery("productOption", "")
	p["productOption"] = c.DefaultQuery("rackCd", "")
	p["productBrandName"] = c.DefaultQuery("productBrandName", "")
	p["inWaybillNo"] = c.DefaultQuery("inWaybillNo", "")
	p["inOrderCd"] = c.DefaultQuery("inOrderCd", "")

	p["partnerExcelFlag"] = route_util.ConvertGetQueryStringToBoolean(c.DefaultQuery("partnerExcelFlag", ""))
	p["productGroupCLOTH"] = route_util.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_CLOTHING", ""))
	p["productGroupACCESSORY"] = route_util.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_ACCESSORY", ""))
	p["productGroupBAGSHOES"] = route_util.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_BAGSHOES", ""))
	p["productGroupCOSMETIC"] = route_util.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_COSMETICS", ""))
	defer func() {
		logger.Log.Debugf("%+v", p)
	}()

	return p
}

// Service //
type stockListExcel struct {
	file         *excelize.File
	sheetName    string
	filename     string
	requestParam map[string]interface{}
}

func newStockListExcel(param map[string]interface{}) *stockListExcel {
	s := stockListExcel{
		file:         excelize.NewFile(),
		requestParam: param,
		sheetName:    "Sheet1",
	}
	s.setFilename()
	return &s
}

func (s *stockListExcel) setFilename() {
	s.filename = "test.xlsx"
}

func (s *stockListExcel) Make() string {
	sheetIndex := s.file.NewSheet(s.sheetName)
	stocks, _ := model.Search(s.requestParam["warehouseId"].(string), s.requestParam)

	var writeResult []error
	channel := make(chan error)

	for i := 0; i < len(stocks); i++ {
		go writeRow(s.file, s.sheetName, i, stocks[i], channel)
	}
	for i := 0; i < len(stocks); i++ {
		c := <-channel
		if c != nil {
			logger.Log.Error(i, c)
		}
		writeResult = append(writeResult, c)
	}

	s.file.SetActiveSheet(sheetIndex)
	if err := s.file.SaveAs(s.filename); err != nil {
		log.Panic(err)
	}
	return s.filename
}

func writeRow(f *excelize.File, sheetName string, i int, stock model.Stock, c chan error) {
	rowNo := strconv.Itoa(i + 2)
	if err := f.SetCellValue(sheetName, "A"+rowNo, i+1); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "B"+rowNo, stock.RackCd); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "C"+rowNo, fmt.Sprintf("%s(%s)", stock.PartnerName, stock.PartnerId)); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "D"+rowNo, stock.ProductOwner); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "E"+rowNo, fmt.Sprintf("%s(%s)", stock.PartnerUserTypeName, stock.PartnerUserType)); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "F"+rowNo, stock.RegDate); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "G"+rowNo, stock.IntervalDate); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "H"+rowNo, stock.ProductItemCd); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "I"+rowNo, stock.ProductCd); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "J"+rowNo, stock.ProductUnitPrice); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "K"+rowNo, stock.ProductVendorPrice); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "L"+rowNo, stock.ProductVendorName); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "M"+rowNo, stock.ProductName); err != nil {
		c <- err
	}
	if err := f.SetCellValue(sheetName, "N"+rowNo, stock.ProductOption); err != nil {
		c <- err
	}

	if stock.PartnerId == "tb67" {
		if err := f.SetCellValue(sheetName, "O"+rowNo, stock.StockBatchNo); err != nil {
			c <- err
		}
		//if err := f.SetCellValue(sheetName, "N"+rowNo, stock.ProductOption); err != nil {
		//	c <- err
		//}
	}
	c <- nil
}
