package entity

type Temperature struct {
	TempC string `json:"temp_C"`
	TempF string `json:"temp_F"`
	TempK string `json:"temp_K"`
}

type WeatherData struct {
	CurrentCondition []struct {
		TempC string `json:"temp_C"`
		TempF string `json:"temp_F"`
	} `json:"current_condition"`
}
