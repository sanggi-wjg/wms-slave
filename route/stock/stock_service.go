package stock

import (
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
)

func makeToStockExcel(param map[string]string) string {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	index := f.NewSheet(sheetName)

	stocks, _ := searchMap(param)

	for i := 0; i < len(stocks); i++ {
		log.Println("A"+strconv.Itoa(i+1), stocks[i].Idx, stocks[i].ProductCd)
		f.SetCellValue(sheetName, "A"+strconv.Itoa(i+1), stocks[i].Idx)
	}
	f.SetActiveSheet(index)

	filename := "test.xlsx"
	if err := f.SaveAs(filename); err != nil {
		panic(err)
	}

	return filename
}

func getFilename() string {
	return ""
}
