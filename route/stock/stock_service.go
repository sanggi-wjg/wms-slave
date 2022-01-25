package stock

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"wms_slave/server/logger"
)

func makeExcelFile(param map[string]interface{}) string {
	filename := getFilename(param)
	sheetName := "Sheet1"

	f := excelize.NewFile()
	sheetIndex := f.NewSheet(sheetName)

	var writeResult []error
	channel := make(chan error)

	stocks, _ := searchMap(param)
	for i := 0; i < len(stocks); i++ {
		go writeRow(f, sheetName, i, stocks[i], channel)
	}
	for i := 0; i < len(stocks); i++ {
		c := <-channel
		if c != nil {
			logger.Log.Error(i, c)
		}
		writeResult = append(writeResult, c)
	}

	f.SetActiveSheet(sheetIndex)
	if err := f.SaveAs(filename); err != nil {
		log.Panic(err)
	}
	return filename
}

func getFilename(param map[string]interface{}) string {
	return "text.xlsx"
}

func writeRow(f *excelize.File, sheetName string, i int, stock Stock, c chan error) {
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
