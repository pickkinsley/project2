package models

import "time"

// Trip is the database model — maps directly to the trips table.
type Trip struct {
	ID            string    `json:"id"`
	Destination   string    `json:"destination"`
	DestLat       float64   `json:"dest_lat"`
	DestLon       float64   `json:"dest_lon"`
	DepartureDate time.Time `json:"departure_date"`
	ReturnDate    time.Time `json:"return_date"`
	TripType      string    `json:"trip_type"`
	Companions    string    `json:"companions"`
	Activities    []string  `json:"activities"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateTripRequest is the JSON body received by POST /api/trips.
// Dates are strings in YYYY-MM-DD format; parsed to time.Time in the handler.
type CreateTripRequest struct {
	Destination   string   `json:"destination"    binding:"required"`
	DepartureDate string   `json:"departure_date" binding:"required"`
	ReturnDate    string   `json:"return_date"    binding:"required"`
	TripType      string   `json:"trip_type"      binding:"required"`
	Companions    string   `json:"companions"     binding:"required"`
	Activities    []string `json:"activities"`
}

// TripResponse is the JSON returned by POST /api/trips and GET /api/trips/{uuid}.
// Dates are formatted as YYYY-MM-DD strings for readability.
// Weather is a pointer so it serializes as null when unavailable.
type TripResponse struct {
	ID            string           `json:"id"`
	Destination   string           `json:"destination"`
	DepartureDate string           `json:"departure_date"`
	ReturnDate    string           `json:"return_date"`
	TripType      string           `json:"trip_type"`
	Companions    string           `json:"companions"`
	Activities    []string         `json:"activities"`
	DurationDays  int              `json:"duration_days"`
	CreatedAt     time.Time        `json:"created_at"`
	Weather       *WeatherResponse `json:"weather"`
	Items         []PackingItem    `json:"items"`
}
