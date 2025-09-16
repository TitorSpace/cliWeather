package weatherapi

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		TempC            float64 `json:"temp_c"`
		Condition        struct {
			Text string `json:"text"`
		} `json:"condition"`
		WindKph  float64 `json:"wind_kph"`
		Humidity int     `json:"humidity"`
	}

	Forecast struct {
		Forecastday []struct {
			Day struct {
				MaxtempC          float64 `json:"maxtemp_c"`
				MintempC          float64 `json:"mintemp_c"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DailyWillItRain   int     `json:"daily_will_it_rain"`
			} `json:"day"`
			Astro struct {
				Sunrise string `json:"sunrise"`
				Sunset  string `json:"sunset"`
			} `json:"astro"`
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			}
		} `json:"forecastday"`
	} `json:"forecast"`
}
