package model

import (
	"gorm.io/gorm"
	"time"
	"wms_slave/server/database"
	"wms_slave/server/e"
)

type Stock struct {
	Idx                  uint64    `gorm:"primaryKey"`
	StockCd              string    `gorm:"column:stockCd"`
	StockType            string    `gorm:"column:stockType"`
	StockTypeName        string    `gorm:"column:stockTypeName"`
	StockVolume          string    `gorm:"column:stockVolume"`
	StockVolumeName      string    `gorm:"column:stockVolumeName"`
	RackCd               string    `gorm:"column:rackCd"`
	RackVolume           string    `gorm:"column:rackVolume"`
	RackVolumeName       string    `gorm:"column:rackVolumeName"`
	PartnerId            string    `gorm:"column:partnerId"`
	PartnerName          string    `gorm:"column:partnerName"`
	PackageCd            string    `gorm:"column:packageCd"`
	ProductOwner         string    `gorm:"column:productOwner"`
	ProductGroup         string    `gorm:"column:productGroup"`
	ProductGroupName     string    `gorm:"column:productGroupName"`
	ProductGroupNameCh   string    `gorm:"column:productGroupNameCh"`
	PartnerProductCd     string    `gorm:"column:partnerProductCd"`
	PartnerProductOption string    `gorm:"column:partnerProductOption"`
	ProductItemCd        string    `gorm:"column:productItemCd"`
	ProductCd            string    `gorm:"column:productCd"`
	ProductName          string    `gorm:"column:productName"`
	ProductNameCh        string    `gorm:"column:productNameCh"`
	ProductNameEn        string    `gorm:"column:productNameEn"`
	ProductOption        string    `gorm:"column:productOption"`
	ProductOptionKr      string    `gorm:"column:productOptionKr"`
	ProductWeight        float32   `gorm:"column:productWeight"`
	ProductImageUrl      string    `gorm:"column:productImageUrl"`
	ProductSize          string    `gorm:"column:productSize"`
	ProductUnitPrice     string    `gorm:"column:productUnitPrice"`
	ProductNature        string    `gorm:"column:productNature"`
	ProductBrandName     string    `gorm:"column:productBrandName"`
	ProductVendorName    string    `gorm:"column:productVendorName"`
	ProductVendorPrice   string    `gorm:"column:productVendorPrice"`
	ProductHscode        string    `gorm:"column:productHscode"`
	ProductQuantity      uint      `gorm:"column:productQuantity"`
	ProductFaultFlag     string    `gorm:"column:productFaultFlag"`
	ProductBarcode       string    `gorm:"column:productBarcode"`
	PartnerUserType      string    `gorm:"column:partnerUserType"`
	PartnerUserTypeName  string    `gorm:"column:partnerUserTypeName"`
	TransferCompany      string    `gorm:"column:transferCompany"`
	TransferCompanyName  string    `gorm:"column:transferCompanyName"`
	RegDate              time.Time `gorm:"column:regDate"`
	Remark               string    `gorm:"column:remark"`
	ExtraData            string    `gorm:"column:extraData"`
	IsCosmeticsOrder     string    `gorm:"column:isCosmeticsOrder"`
	StockBatchNo         string    `gorm:"column:stockBatchNo"`
}

func (Stock) TableName() string {
	return "stock"
}

func FindById(id uint64) (Stock, error) {
	var stock Stock
	res := database.DB.Where(Stock{Idx: id}).Take(&stock)
	if res.Error != nil {
		return stock, res.Error
	}
	return stock, nil
}

func FindByStockCd(stockCd string) (*Stock, error) {
	var stock Stock
	res := database.DB.Where(Stock{StockCd: stockCd}).Take(&stock)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, e.NewExceptionAddMsg(e.ErrorStockNotFound, stockCd)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &stock, nil
}

func SearchMap(param map[string]interface{}) ([]Stock, error) {
	var stocks []Stock
	//res := database.DB.Where(param).Order("regDate").Find(&stocks)
	db := database.DB

	if param["fromDate"] != "" {
		t, _ := time.Parse("2006-01-02", param["fromDate"].(string))
		ts := t.Format("2006-01-02 15:04:05")
		db = db.Where("regDate >= ?", ts)
	}
	if param["toDate"] != "" {
		t, _ := time.Parse("2006-01-02", param["toDate"].(string))
		ts := t.Format("2006-01-02") + " 23:59:59"
		db = db.Where("regDate <= ?", ts)
	}
	if param["partnerId"] != "" {
		db = db.Where("partnerId = ?", param["partnerId"])
	}
	if param["partnerUserType"] != "" {
		db = db.Where("partnerUserType = ?", param["partnerUserType"])
	}
	if param["transferCompany"] != "" {
		db = db.Where("transferCompany = ?", param["transferCompany"])
	}

	res := db.Order("regDate").Find(&stocks)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return stocks, nil
}
