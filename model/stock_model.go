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

	// Additional
	IntervalDate uint   `gorm:"column:intervalDate"`
	ExpireDate   string `gorm:"column:expireDate"`
}

func (Stock) TableName() string {
	return "stock"
}

func FindById(warehouseId string, id uint64) (Stock, error) {
	var stock Stock
	res := database.DB[warehouseId].Where(Stock{Idx: id}).Take(&stock)
	if res.Error != nil {
		return stock, res.Error
	}
	return stock, nil
}

func FindByStockCd(warehouseId string, stockCd string) (*Stock, error) {
	var stock Stock
	res := database.DB[warehouseId].Where(Stock{StockCd: stockCd}).Take(&stock)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, e.NewExceptionAddMsg(e.ErrorStockNotFound, stockCd)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &stock, nil
}

func Search(warehouseId string, param map[string]interface{}) ([]Stock, error) {
	var stocks []Stock
	db := database.DB[warehouseId]

	db = db.Select("rackCd, partnerId, partnerName, productOwner, partnerUserType, partnerUserTypeName, regDate, productItemCd, " +
		"productCd, productUnitPrice, productVendorPrice, productVendorName, productName, productOption, stockBatchNo, " +
		"TO_DAYS(CURRENT_DATE())-TO_DAYS(DATE_FORMAT(regDate,'%Y-%m-%d')) AS intervalDate, " +
		"STR_TO_DATE(RIGHT(stockBatchNo, 6), '%y%m%d') AS expireDate ")

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
	if param["code"] != "" {
		db = db.Where("( stockCd = ? OR packageCd = ? productItemCd = ? OR productCd = ? OR rackCd = ?)",
			param["code"], param["code"], param["code"], param["code"], param["code"])
	}
	if param["inWaybillNo"] != "" {
		inLogSubQuery := database.DB[warehouseId].Table("in_log").Select("productCd")
		db = db.Where("productCd IN (?)", inLogSubQuery.Where("inWaybillNo = ?", param["inWaybillNo"]))
	}
	if param["inOrderCd"] != "" {
		inLogSubQuery := database.DB[warehouseId].Table("in_log").Select("productCd")
		db = db.Where("productCd IN (?)", inLogSubQuery.Where("inOrderCd = ?", param["inOrderCd"]))
	}
	if param["productGroupCLOTH"].(bool) {
		db = db.Where("productGroup IN ('TS','BL','HD','KN','DD','BX','CT','CS','JK','CA','VT','YR','XZ','PD','OP','LS','MS','SS','LP','MP','SP','JN','LG','JS','CL','ST','CE','SE','SC','MF','GL','NY','SY','YJ','BK','LD','QT') ")
	}
	if param["productGroupACCESSORY"].(bool) {
		db = db.Where("productGroup IN ('ER','XL','RG','BR','JL','BC','HA','ET','BE','GS','HT','GJ','KC','PC','PJ','AT') ")
	}
	if param["productGroupBAGSHOES"].(bool) {
		db = db.Where("productGroup IN ('HB','CB','SB','DB','MB','BB','WA','CP','SH','BT','SD','HS','SL','HH','FS','PX') ")
	}
	if param["productGroupCOSMETIC"].(bool) {
		db = db.Where("productGroup IN ('CM') ")
	}

	db = database.IfEqual(db, "partnerId", param["partnerId"])
	db = database.IfEqual(db, "productOwner", param["productOwner"])
	db = database.IfEqual(db, "productVendorName", param["productVendorName"])
	db = database.IfEqual(db, "stockType", param["stockType"])
	db = database.IfEqual(db, "partnerUserType", param["partnerUserType"])
	db = database.IfEqual(db, "transferCompany", param["transferCompany"])
	db = database.IfEqual(db, "rackCd", param["rackCd"])

	db = database.IfContains(db, "productName", param["productName"].(string))
	db = database.IfContains(db, "productOption", param["productOption"].(string))
	db = database.IfContains(db, "productBrandName", param["productBrandName"].(string))

	res := db.Order("regDate").Find(&stocks)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return stocks, nil
}
