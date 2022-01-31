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
	p["warehouseDomain"] = c.DefaultQuery("WAREHOUSE_DOMAIN", "")
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

	p["partnerExcelFlag"] = global.ConvertGetQueryStringToBoolean(c.DefaultQuery("partnerExcelFlag", ""))
	p["productGroupCLOTH"] = global.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_CLOTHING", ""))
	p["productGroupACCESSORY"] = global.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_ACCESSORY", ""))
	p["productGroupBAGSHOES"] = global.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_BAGSHOES", ""))
	p["productGroupCOSMETIC"] = global.ConvertGetQueryStringToBoolean(c.DefaultQuery("productGroup_COSMETICS", ""))
	defer func() {
		logger.Log.Debugf("%+v", p)
	}()

	return p
}

// Service //
type stockListExcel struct {
	file         *excelize.File
	excelDetail  *excelDetail
	requestParam map[string]interface{}
}

type excelDetail struct {
	sheetName  string
	filename   string
	sheetIndex int
}

func newStockListExcel(param map[string]interface{}) *stockListExcel {
	s := stockListExcel{
		file:         excelize.NewFile(),
		requestParam: param,
		excelDetail: &excelDetail{
			sheetName: "Sheet1",
		},
	}
	s.setFilename()
	return &s
}

func (s *stockListExcel) setFilename() {
	s.excelDetail.filename = "test.xlsx"
}

func (s *stockListExcel) Make() string {
	sheetIndex := s.file.NewSheet(s.excelDetail.sheetName)
	s.setColWidth()

	var writeResult []error
	channel := make(chan error)
	stocks, _ := model.Search(s.requestParam["warehouseId"].(string), s.requestParam)

	for i := 0; i < len(stocks); i++ {
		go writeRow(s.file, s.excelDetail.sheetName, i, stocks[i], channel)
	}
	for i := 0; i < len(stocks); i++ {
		c := <-channel
		if c != nil {
			logger.Log.Error(i, c)
		}
		writeResult = append(writeResult, c)
	}

	s.setStyle(len(stocks))
	s.file.SetActiveSheet(sheetIndex)
	if err := s.file.SaveAs(s.excelDetail.filename); err != nil {
		log.Panic(err)
	}
	return s.excelDetail.filename
}

func (s *stockListExcel) setColWidth() {
	if err := s.file.SetColWidth(s.excelDetail.sheetName, "A", "A", 10); err != nil {
		logger.Log.Error("SetColWidth A error", err)
	}
	if err := s.file.SetColWidth(s.excelDetail.sheetName, "B", "L", 20); err != nil {
		logger.Log.Error("SetColWidth B:L error", err)
	}
	if err := s.file.SetColWidth(s.excelDetail.sheetName, "M", "N", 100); err != nil {
		logger.Log.Error("SetColWidth M:N error", err)
	}
}

func (s *stockListExcel) setStyle(size int) {
	style, err := s.file.NewStyle(`{"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]}`)
	if err != nil {
		fmt.Println(err)
	}
	err = s.file.SetCellStyle(s.excelDetail.sheetName, "A1", "N"+strconv.Itoa(size+1), style)
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
