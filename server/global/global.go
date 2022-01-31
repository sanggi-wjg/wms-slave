package global

func GetWarehouseIdByDomain(warehouseDomain string) string {
	switch warehouseDomain {
	case "kr01.warehouse.pickby.us":
		return "KR01"
	case "cn02.warehouse.pickby.us":
		return "CN02"
	}
	return ""
}

func ConvertGetQueryStringToBoolean(queryValue string) bool {
	if queryValue == "Y" || queryValue == "T" {
		return true
	}
	return false
}
