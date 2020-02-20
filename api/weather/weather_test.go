package weather

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/pisign/pisign-backend/mocks"
	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

func fakeWeatherResponse() OpenWeatherResponse {
	output := OpenWeatherResponse{
		Coord: coord{
			Lon: -122.3,
			Lat: 47.66,
		},
		Weather: []weather{
			{
				ID:          501,
				Main:        "Rain",
				Description: "moderate rain",
				Icon:        "10n",
			},
			{
				ID:          701,
				Main:        "Mist",
				Description: "mist",
				Icon:        "50n",
			},
		},
		Base: "stations",
		Main: main{
			Temp:      281.71,
			FeelsLike: 280.02,
			TempMin:   279.26,
			TempMax:   283.71,
			Pressure:  1017,
			Humidity:  100,
		},
		Visibility: 6437,
		Wind: wind{
			Speed: 1.94,
			Deg:   207,
		},
		Rain: rain{
			OneHR: 1.43,
		},
		Clouds: clouds{
			All: 90,
		},
		DT: 1580953406,
		Sys: sys{
			Type:    1,
			ID:      2674,
			Country: "US",
			Sunrise: 1580916664,
			Sunset:  1580951718,
		},
		Timezone: -28800,
		ID:       0,
		Name:     "Seattle",
		COD:      200,
	}
	return output
}

// Example GoMock test
func TestGetNameWithGoMocks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockAPI := mocks.NewMockAPI(mockCtrl)

	defer mockCtrl.Finish()

	mockAPI.EXPECT().GetName().Return("fake")
	if mockAPI.GetName() != "fake" {
		t.Error("error with go mock")
	}
}

func TestConfigure(t *testing.T) {
	a := NewAPI(nil, nil, uuid.New())
	if err := a.Configure(types.ClientMessage{}); err == nil {
		t.Error("somehow we didn't return an error after configuring junk")
	}
}

func TestGetName(t *testing.T) {
	a := NewAPI(nil, nil, uuid.New())
	if a.GetName() != "weather" {
		t.Error("somehow we didn't return the name")
	}
}

func TestNewAPI(t *testing.T) {
	testAPI := NewAPI(nil, nil, uuid.New())
	if testAPI.Name != "weather" || testAPI.ValidCache {
		t.Error("weather.NewAPI does not set name or valid cache correctly")
	}
}

func TestBuildUrl(t *testing.T) {
	url := buildurl(90000, "API_KEY")
	if url != "https://api.openweathermap.org/data/2.5/weather?zip=90000,us&APPID=API_KEY" {
		t.Error("Generated URL is not correct: " + url)
	}
}

func TestKelvinToF(t *testing.T) {
	if (kelvinToF(0) - -459.67) > 1 || (kelvinToF(0) - -459.67) < -1 {
		t.Error("conversion error")
		t.Error(kelvinToF(0))
		t.Error(-459.67)
	}
}

func TestParseJSON(t *testing.T) {
	JSONdata := []byte("{\"coord\":{\"lon\":-122.3,\"lat\":47.66},\"weather\":[{\"id\":501,\"main\":\"Rain\",\"description\":\"moderate rain\",\"icon\":\"10n\"},{\"id\":701,\"main\":\"Mist\",\"description\":\"mist\",\"icon\":\"50n\"}],\"base\":\"stations\",\"main\":{\"temp\":281.71,\"feels_like\":280.02,\"temp_min\":279.26,\"temp_max\":283.71,\"pressure\":1017,\"humidity\":100},\"visibility\":6437,\"wind\":{\"speed\":1.94,\"deg\":207},\"rain\":{\"1h\":1.43},\"clouds\":{\"all\":90},\"dt\":1580953406,\"sys\":{\"type\":1,\"id\":2674,\"country\":\"US\",\"sunrise\":1580916664,\"sunset\":1580951718},\"timezone\":-28800,\"id\":0,\"name\":\"Seattle\",\"cod\":200}")

	exampledata := fakeWeatherResponse()

	var parsedJSON OpenWeatherResponse
	utils.ParseJSON(JSONdata, &parsedJSON)

	if !reflect.DeepEqual(exampledata, parsedJSON) {
		t.Error("Error in parsed json! Does not match excpected struct")
		t.Error(parsedJSON)
		t.Error(exampledata)
	}
}

func TestTransform(t *testing.T) {
	name := "Seattle"
	openWeatherResponse := fakeWeatherResponse()

	resp := openWeatherResponse.Transform()
	weatherResponse := *resp.(*types.WeatherResponse)

	if weatherResponse.Name != name || weatherResponse.Main.Temp < 0 ||
		weatherResponse.Main.TempMin < 0 || weatherResponse.Main.TempMax < 0 {
		t.Error("transformation error")
		t.Error(openWeatherResponse)
		t.Error(weatherResponse)
	}
}
