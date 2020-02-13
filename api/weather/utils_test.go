package weather

import "testing"

func Test_buildurl(t *testing.T) {
	url := buildurl(90000, "API_KEY")
	if url != "https://api.openweathermap.org/data/2.5/weather?zip=90000,us&APPID=API_KEY" {
		t.Error("Generated URL is not correct: " + url)
	}
}

func Test_kelvinToF(t *testing.T) {
	if (kelvinToF(0) - -459.67) > 1 || (kelvinToF(0) - -459.67) < -1 {
		t.Error("conversion error")
		t.Error(kelvinToF(0))
		t.Error(-459.67)
	}
}
