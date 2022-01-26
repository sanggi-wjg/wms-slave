package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"wms_slave/route/stock/model"
	"wms_slave/server/logger"
)

type StockListExcel struct {
	file         *excelize.File
	sheetName    string
	filename     string
	requestParam map[string]interface{}
}

func NewStockListExcel(param map[string]interface{}) *StockListExcel {
	s := StockListExcel{
		file:         excelize.NewFile(),
		requestParam: param,
		sheetName:    "Sheet1",
	}
	s.setFilename()
	return &s
}

func (s *StockListExcel) setFilename() {
	s.filename = "test.xlsx"
}

func (s *StockListExcel) GetFilename() string {
	return s.filename
}

func (s *StockListExcel) Make() {
	sheetIndex := s.file.NewSheet(s.sheetName)
	stocks, _ := model.SearchMap(s.requestParam)

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
	c <- nil
}