package stock

import (
	"gorm.io/gorm"
	"time"
	"wms_slave/server/database"
	"wms_slave/server/e"
)

type Stock struct {
	Idx                  uint64    `gorm:"primaryKey"`
	StockCd              string    `json:"stockCd"`
	StockType            string    `json:"stockType"`
	StockTypeName        string    `json:"stockTypeName"`
	StockVolume          string    `json:"stockVolume"`
	StockVolumeName      string    `json:"stockVolumeName"`
	RackCd               string    `json:"rackCd"`
	RackVolume           string    `json:"rackVolume"`
	RackVolumeName       string    `json:"rackVolumeName"`
	PartnerId            string    `json:"partnerId"`
	PartnerName          string    `json:"partnerName"`
	PackageCd            string    `json:"packageCd"`
	ProductOwner         string    `json:"productOwner"`
	ProductGroup         string    `json:"productGroup"`
	ProductGroupName     string    `json:"productGroupName"`
	ProductGroupNameCh   string    `json:"productGroupNameCh"`
	PartnerProductCd     string    `json:"partnerProductCd"`
	PartnerProductOption string    `json:"partnerProductOption"`
	ProductItemCd        string    `json:"productItemCd"`
	ProductCd            string    `json:"productCd"`
	ProductName          string    `json:"productName"`
	ProductNameCh        string    `json:"productNameCh"`
	ProductNameEn        string    `json:"productNameEn"`
	ProductOption        string    `json:"productOption"`
	ProductOptionKr      string    `json:"productOptionKr"`
	ProductWeight        float32   `json:"productWeight"`
	ProductImageUrl      string    `json:"productImageUrl"`
	ProductSize          string    `json:"productSize"`
	ProductUnitPrice     string    `json:"productUnitPrice"`
	ProductNature        string    `json:"productNature"`
	ProductBrandName     string    `json:"productBrandName"`
	ProductVendorName    string    `json:"productVendorName"`
	ProductVendorPrice   string    `json:"productVendorPrice"`
	ProductHscode        string    `json:"productHscode"`
	ProductQuantity      uint      `json:"productQuantity"`
	ProductFaultFlag     string    `json:"productFaultFlag"`
	ProductBarcode       string    `json:"productBarcode"`
	PartnerUserType      string    `json:"partnerUserType"`
	PartnerUserTypeName  string    `json:"partnerUserTypeName"`
	TransferCompany      string    `json:"transferCompany"`
	TransferCompanyName  string    `json:"transferCompanyName"`
	RegDate              time.Time `json:"regDate"`
	Remark               string    `json:"remark"`
	ExtraData            string    `json:"extraData"`
	IsCosmeticsOrder     string    `json:"isCosmeticsOrder"`
	StockBatchNo         string    `json:"stockBatchNo"`
}

func (Stock) TableName() string {
	return "stock"
}

func findById(id uint64) (Stock, error) {
	var stock Stock
	res := database.DB.Where(Stock{Idx: id}).Take(&stock)
	if res.Error != nil {
		return stock, res.Error
	}
	return stock, nil
}

func findByStockCd(stockCd string) (*Stock, error) {
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

func searchMap(param map[string]string) ([]Stock, error) {
	var stocks []Stock
	res := database.DB.Find(&stocks).Order("RegDate")

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return stocks, nil
}
