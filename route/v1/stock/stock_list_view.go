package stock

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"time"
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
	//filepath   string
	filename   string
	sheetName  string
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
	var filename string
	switch s.requestParam["warehouseId"] {
	case "KR01":
		filename += "BuCheon[KR01]"
	case "CN02":
		filename += "SagHai[CN02]"
	}
	ts := time.Now().Format("2006.01.02.15.04.05")
	filename += "_재고리스트_" + ts

	//s.excelDetail.filepath = fmt.Sprintf("media/%s.xlsx", filename)
	s.excelDetail.filename = filename + ".xlsx"
}

func (s *stockListExcel) Make() string {
	sheetIndex := s.file.NewSheet(s.excelDetail.sheetName)

	length := s.write()
	s.writeTitle()
	s.setColWidth()
	s.setStyle(length)

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
	if s.requestParam["partnerId"] != "" {
		switch s.requestParam["partnerId"] {
		case "tb67":
		case "yz03":
		case "zd01":
			if err := s.file.SetColWidth(s.excelDetail.sheetName, "O", "P", 20); err != nil {
				logger.Log.Error("SetColWidth O:P error", err)
			}
		}
	}

}

func (s *stockListExcel) writeTitle() {
	type title struct {
		col  string
		name string
	}
	titles := []title{
		{col: "A1", name: "번호"},
		{col: "B1", name: "랙코드"},
		{col: "C1", name: "파트너 아이디"},
		{col: "D1", name: "재고 소유주"},
		{col: "E1", name: "고객유형"},
		{col: "F1", name: "입고날짜"},
		{col: "G1", name: "적재기간(日)"},
		{col: "H1", name: "아이템코드(PIC)"},
		{col: "I1", name: "상품코드(STC)"},
		{col: "J1", name: "상품단위가격"},
		{col: "K1", name: "상품구매가격"},
		{col: "L1", name: "상품벤더"},
		{col: "M1", name: "상품명"},
		{col: "N1", name: "옵션명"},
	}
	switch s.requestParam["partnerId"] {
	case "tb67":
	case "yz03":
	case "zd01":
		titles = append(titles, title{col: "O1", name: "관리코드"})
		titles = append(titles, title{col: "P1", name: "유호기간"})

	}

	for _, t := range titles {
		if err := s.file.SetCellValue(s.excelDetail.sheetName, t.col, t.name); err != nil {
			logger.Log.Error("writeTitleError", t, err)
		}
	}
}

func (s *stockListExcel) setStyle(size int) {
	// TODO : 스타일 생성자 구조체 만들어 놓고 쓰면 될 듯?
	borderStyle, err := s.file.NewStyle(`{"border":[{"type":"left","color":"000000","style":1},
													{"type":"top","color":"000000","style":1},
													{"type":"bottom","color":"000000","style":1},
													{"type":"right","color":"000000","style":1}]}`)
	if err != nil {
		logger.Log.Error(err)
	}
	fillStyle, err2 := s.file.NewStyle(`{"fill":{"type":"pattern","color":["#dcdcdc"],"pattern":1}}`)
	if err2 != nil {
		logger.Log.Error(err2)
	}

	switch s.requestParam["partnerId"] {
	case "tb67":
	case "yz03":
	case "zd01":
		if err = s.file.SetCellStyle(s.excelDetail.sheetName, "A1", "P1", fillStyle); err != nil {
			logger.Log.Error(err)
		}
		if err = s.file.SetCellStyle(s.excelDetail.sheetName, "A1", "P"+strconv.Itoa(size+1), borderStyle); err != nil {
			logger.Log.Error(err)
		}
	default:
		if err = s.file.SetCellStyle(s.excelDetail.sheetName, "A1", "N1", fillStyle); err != nil {
			logger.Log.Error(err)
		}
		if err = s.file.SetCellStyle(s.excelDetail.sheetName, "A1", "N"+strconv.Itoa(size+1), borderStyle); err != nil {
			logger.Log.Error(err)
		}
	}

}

func (s *stockListExcel) write() int {
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

	return len(stocks)
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
	regDate := stock.RegDate.Format("2006-01-02 15:04:05")
	if err := f.SetCellValue(sheetName, "F"+rowNo, regDate); err != nil {
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
		if err := f.SetCellValue(sheetName, "P"+rowNo, stock.ExpireDate); err != nil {
			c <- err
		}
	}
	c <- nil
}
