package weatherUtils

func UpdateWeather() string {
	return "nothing"
}

type coord_weather struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

type weather_weather struct {
	Id          int    `json:"id,omitempty"`
	MainW       string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type main_weather struct {
	Temp      float64 `json:"temp,omitempty"`
	FeelsLike float64 `json:"feels_like,omitempty"`
	TempMin   float64 `json:"temp_min,omitempty"`
	TempMax   float64 `json:"temp_max,omitempty"`
	Pressure  int     `json:"pressure,omitempty"`
	Humidity  int     `json:"humidity,omitempty"`
}

type wind_weather struct {
	Speed float64 `json:"speed,omitempty"`
	Deg   int     `json:"deg,omitempty"`
	Gust  float64 `json:"gust,omitempty"`
}

type clouds_weather struct {
	All int `json:"all,omitempty"`
}

type sys_weather struct {
	TypeField int    `json:"type,omitempty"`
	Id        int    `json:"lid,omitempty"`
	Country   string `json:"country,omitempty"`
	Sunrise   int    `json:"sunrise,omitempty"`
	Sunset    int    `json:"sunset,omitempty"`
}

type Weather_main struct {
	Coord      coord_weather     `json:"coord"`
	Weather    []weather_weather `json:"weather"`
	Base       string            `json:"base"`
	Main       main_weather      `json:"main"`
	Visibility int               `json:"visibility"`
	Wind       wind_weather      `json:"wind"`
	Clouds     clouds_weather    `json:"clouds"`
	Dt         int               `json:"dt"`
	Sys        sys_weather       `json:"sys"`
	Timezone   int               `json:"timezone"`
	Id         int               `json:"id"`
	Name       string            `json:"name"`
	Cod        int               `json:"cod"`
}
