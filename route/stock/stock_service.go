package stock

import (
	"github.com/xuri/excelize/v2"
	"strconv"
)

func makeToStockExcel() string {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	index := f.NewSheet(sheetName)

	stocks, _ := findAll()
	for i := 0; i < len(stocks); i++ {
		f.SetCellValue(sheetName, "A"+strconv.Itoa(i+1), stocks[i].stockCd)
	}
	f.SetActiveSheet(index)

	filename := "test.xlsx"
	if err := f.SaveAs(filename); err != nil {
		panic(err)
	}

	return filename
}
