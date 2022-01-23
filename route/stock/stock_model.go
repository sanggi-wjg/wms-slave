package stock

import (
	"gorm.io/gorm"
	"time"
	"wms_slave/server/database"
	"wms_slave/server/e"
)

type Stock struct {
	idx                  uint64
	stockCd              string
	stockType            string
	stockTypeName        string
	stockVolume          string
	stockVolumeName      string
	rackCd               string
	rackVolume           string
	rackVolumeName       string
	partnerId            string
	partnerName          string
	packageCd            string
	productOwner         string
	productGroup         string
	productGroupName     string
	productGroupNameCh   string
	partnerProductCd     string
	partnerProductOption string
	productItemCd        string
	productCd            string
	productName          string
	productNameCh        string
	productNameEn        string
	productOption        string
	productOptionKr      string
	productWeight        float32
	productImageUrl      string
	productSize          string
	productUnitPrice     string
	productNature        string
	productBrandName     string
	productVendorName    string
	productVendorPrice   string
	productHscode        string
	productQuantity      uint
	productFaultFlag     string
	productBarcode       string
	partnerUserType      string
	partnerUserTypeName  string
	transferCompany      string
	transferCompanyName  string
	regDate              time.Time
	remark               string
	extraData            string
	isCosmeticsOrder     string
	stockBatchNo         string
}

func (Stock) TableName() string {
	return "stock"
}

func findByStockCd(stockCd string) (*Stock, error) {
	var stock Stock
	res := database.DB.Where(Stock{stockCd: stockCd}).Take(&stock)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, e.NewExceptionAddMsg(e.ErrorStockNotFound, stockCd)
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return &stock, nil
}

func findAll() ([]*Stock, error) {
	var stocks []*Stock
	res := database.DB.Find(&stocks).Order("regDate")
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return stocks, nil
}
