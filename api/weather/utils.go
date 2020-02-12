package weather

import (
	"fmt"
)

const endpoint = "https://api.openweathermap.org/data/2.5/weather"

// Builds the URL from the endpoint, adding required args
func buildurl(zipcode int, apikey string) string {
	url := fmt.Sprintf("%s?zip=%d,us&APPID=%s", endpoint, zipcode, apikey)
	return url
}

func kelvinToF(k float64) float64 {
	return (k-273.15)*(9/5) + 32
}
