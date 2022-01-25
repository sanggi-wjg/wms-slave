package stock

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
)

func makeExcelFile(param map[string]string) string {
	filename := getFilename(param)
	sheetName := "Sheet1"

	f := excelize.NewFile()
	sheetIndex := f.NewSheet(sheetName)

	var writeResult []bool
	channel := make(chan bool)

	stocks, _ := searchMap(param)
	for i := 0; i < len(stocks); i++ {
		go writeRow(f, sheetName, i, stocks[i], channel)
	}
	for i := 0; i < len(stocks); i++ {
		writeResult = append(writeResult, <-channel)
	}
	log.Println("writeResult", writeResult)

	f.SetActiveSheet(sheetIndex)
	if err := f.SaveAs(filename); err != nil {
		log.Panic(err)
	}
	return filename
}

func getFilename(param map[string]string) string {
	return "text.xlsx"
}

func writeRow(f *excelize.File, sheetName string, i int, stock Stock, c chan bool) {
	rowNo := strconv.Itoa(i + 2)
	if err := f.SetCellValue(sheetName, "A"+rowNo, i+1); err != nil {
		c <- false
	}
	if err := f.SetCellValue(sheetName, "B"+rowNo, stock.RackCd); err != nil {
		c <- false
	}
	if err := f.SetCellValue(sheetName, "C"+rowNo, fmt.Sprintf("%s(%s)", stock.PartnerName, stock.PartnerId)); err != nil {
		c <- false
	}
	c <- true
}
