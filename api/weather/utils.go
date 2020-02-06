package weather

import (
	"fmt"
)

const endpoint = "https://api.openweathermap.org/data/2.5/weather?zip="

// Builds the URL from the endpoint, adding required args
func buildurl(zipcode int, apikey string) string {
	return fmt.Sprintf("%s%d,us&APPID=%s", endpoint, zipcode, apikey)
}
