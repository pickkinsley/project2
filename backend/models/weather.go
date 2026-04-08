package models

import "time"

// WeatherSnapshot maps to the weather_snapshots table.
type WeatherSnapshot struct {
	ID            int             `json:"id"`
	TripID        string          `json:"-"`
	TempMinF      int             `json:"temp_min_f"`
	TempMaxF      int             `json:"temp_max_f"`
	RainDays      int             `json:"rain_days"`
	SnowDays      int             `json:"snow_days"`
	IsForecast    bool            `json:"is_forecast"`
	DailyForecast []DailyForecast `json:"daily_forecast"`
	FetchedAt     time.Time       `json:"-"`
}

// WeatherResponse is the weather object returned inside TripResponse.
// It is a pointer in TripResponse so it serializes as null when unavailable.
type WeatherResponse struct {
	TempMinF      int             `json:"temp_min_f"`
	TempMaxF      int             `json:"temp_max_f"`
	RainDays      int             `json:"rain_days"`
	SnowDays      int             `json:"snow_days"`
	IsForecast    bool            `json:"is_forecast"`
	DailyForecast []DailyForecast `json:"daily_forecast"`
}

// DailyForecast represents one day in the weather strip.
// Stored as a JSON array in the daily_forecast column of weather_snapshots.
type DailyForecast struct {
	Date string `json:"date"`   // YYYY-MM-DD
	Icon string `json:"icon"`   // sunny | partly_cloudy | cloudy | rainy | snowy | stormy
	MinF int    `json:"min_f"`
	MaxF int    `json:"max_f"`
}
