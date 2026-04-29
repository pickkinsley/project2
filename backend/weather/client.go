// Package weather provides geocoding and forecast retrieval via the Open-Meteo API.
// Open-Meteo is free and requires no API key.
package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/pickkinsley/project2/backend/models"
)

const (
	geocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/search"
	forecastBaseURL  = "https://api.open-meteo.com/v1/forecast"
	// ForecastMaxDays is the maximum number of days ahead the free Open-Meteo
	// forecast API supports. Trips departing beyond this window get nil weather.
	ForecastMaxDays = 16
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GeocodedLocation holds the resolved coordinates for a destination string.
type GeocodedLocation struct {
	Latitude  float64
	Longitude float64
}

// ── internal response shapes ─────────────────────────────────────────────────

type geocodingResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type forecastResponse struct {
	Daily struct {
		Time                 []string  `json:"time"`
		WeatherCode          []int     `json:"weather_code"`
		TemperatureMax       []float64 `json:"temperature_2m_max"`
		TemperatureMin       []float64 `json:"temperature_2m_min"`
		PrecipProbabilityMax []int     `json:"precipitation_probability_max"`
		WindSpeedMax         []float64 `json:"wind_speed_10m_max"`
	} `json:"daily"`
}

// ── public API ───────────────────────────────────────────────────────────────

// GeocodeLocation resolves a destination string to lat/lon using Open-Meteo
// geocoding. Returns (nil, nil) when the destination is not found.
func GeocodeLocation(destination string) (*GeocodedLocation, error) {
	params := url.Values{}
	params.Set("name", destination)
	params.Set("count", "1")
	params.Set("language", "en")
	params.Set("format", "json")

	resp, err := httpClient.Get(geocodingBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var result geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("geocoding decode failed: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, nil // destination not found — caller should return 400
	}

	return &GeocodedLocation{
		Latitude:  result.Results[0].Latitude,
		Longitude: result.Results[0].Longitude,
	}, nil
}

// FetchForecast retrieves daily weather data for the given coordinates and date
// range from Open-Meteo. Returns nil on API failure so callers can fall back
// gracefully rather than blocking trip creation.
func FetchForecast(lat, lon float64, startDate, endDate string) (*models.WeatherResponse, error) {
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", lat))
	params.Set("longitude", fmt.Sprintf("%f", lon))
	params.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min,precipitation_probability_max,wind_speed_10m_max")
	params.Set("temperature_unit", "fahrenheit")
	params.Set("wind_speed_unit", "mph")
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	params.Set("timezone", "auto")

	resp, err := httpClient.Get(forecastBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("forecast request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	var result forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("forecast decode failed: %w", err)
	}

	return buildWeatherResponse(result), nil
}

// ── internal helpers ─────────────────────────────────────────────────────────

func buildWeatherResponse(f forecastResponse) *models.WeatherResponse {
	days := f.Daily.Time
	if len(days) == 0 {
		return nil
	}

	overallMin := int(math.Round(safeFloat(f.Daily.TemperatureMin, 0)))
	overallMax := int(math.Round(safeFloat(f.Daily.TemperatureMax, 0)))
	var rainDays, snowDays int
	forecasts := make([]models.DailyForecast, 0, len(days))

	for i, date := range days {
		code := safeInt(f.Daily.WeatherCode, i)
		icon := wmoToIcon(code)
		minF := int(math.Round(safeFloat(f.Daily.TemperatureMin, i)))
		maxF := int(math.Round(safeFloat(f.Daily.TemperatureMax, i)))
		precip := safeInt(f.Daily.PrecipProbabilityMax, i)
		wind := math.Round(safeFloat(f.Daily.WindSpeedMax, i)*10) / 10 // 1 decimal place

		if minF < overallMin {
			overallMin = minF
		}
		if maxF > overallMax {
			overallMax = maxF
		}

		switch icon {
		case "rainy", "stormy":
			rainDays++
		case "snowy":
			snowDays++
		}

		forecasts = append(forecasts, models.DailyForecast{
			Date:              date,
			Icon:              icon,
			MinF:              minF,
			MaxF:              maxF,
			PrecipProbability: precip,
			WindSpeedMph:      wind,
		})
	}

	return &models.WeatherResponse{
		TempMinF:      overallMin,
		TempMaxF:      overallMax,
		RainDays:      rainDays,
		SnowDays:      snowDays,
		IsForecast:    true,
		DailyForecast: forecasts,
	}
}

// wmoToIcon maps WMO Weather Interpretation Codes to the app's icon vocabulary.
// Reference: https://open-meteo.com/en/docs (scroll to "WMO Weather Code")
func wmoToIcon(code int) string {
	switch {
	case code == 0:
		return "sunny" // Clear sky
	case code <= 2:
		return "partly_cloudy" // Mainly clear, partly cloudy
	case code <= 48:
		return "cloudy" // Overcast, fog, depositing rime fog
	case code <= 55:
		return "rainy" // Drizzle: light, moderate, dense
	case code <= 57:
		return "snowy" // Freezing drizzle
	case code <= 65:
		return "rainy" // Rain: slight, moderate, heavy
	case code <= 67:
		return "snowy" // Freezing rain
	case code <= 77:
		return "snowy" // Snow fall, snow grains
	case code <= 82:
		return "rainy" // Rain showers
	case code <= 86:
		return "snowy" // Snow showers
	default:
		return "stormy" // Thunderstorm (95, 96, 99)
	}
}

func safeInt(s []int, i int) int {
	if i < len(s) {
		return s[i]
	}
	return 0
}

func safeFloat(s []float64, i int) float64 {
	if i < len(s) {
		return s[i]
	}
	return 0
}
