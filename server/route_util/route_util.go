package route_util

func ConvertGetQueryStringToBoolean(queryValue string) bool {
	if queryValue == "Y" || queryValue == "T" {
		return true
	}
	return false
}
