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
