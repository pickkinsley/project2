package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pickkinsley/project2/backend/models"
)

// mockTripID is the UUID returned by the mock handler.
// Replace with real database lookups in Lesson 2.
const mockTripID = "a3f8c2d1-4b5e-4c3d-8f9a-1b2c3d4e5f6a"

var mockTrip = models.TripResponse{
	ID:            mockTripID,
	Destination:   "Paris, France",
	DepartureDate: "2026-04-10",
	ReturnDate:    "2026-04-15",
	TripType:      "international",
	Companions:    "couple",
	Activities:    []string{"sightseeing", "fine_dining"},
	DurationDays:  5,
	Weather: &models.WeatherResponse{
		TempMinF:   52,
		TempMaxF:   63,
		RainDays:   2,
		SnowDays:   0,
		IsForecast: true,
		DailyForecast: []models.DailyForecast{
			{Date: "2026-04-10", Icon: "partly_cloudy", MinF: 50, MaxF: 61},
			{Date: "2026-04-11", Icon: "rainy", MinF: 48, MaxF: 55},
			{Date: "2026-04-12", Icon: "sunny", MinF: 52, MaxF: 63},
			{Date: "2026-04-13", Icon: "partly_cloudy", MinF: 51, MaxF: 60},
			{Date: "2026-04-14", Icon: "rainy", MinF: 49, MaxF: 57},
		},
	},
	Items: []models.PackingItem{
		{ID: 1, Name: "Passport", Category: "Essential Items", IsEssential: true, Reason: "Required for international travel", IsChecked: false, SortOrder: 1},
		{ID: 2, Name: "Prescriptions", Category: "Essential Items", IsEssential: true, Reason: "Never travel without your medications", IsChecked: false, SortOrder: 2},
		{ID: 3, Name: "Light jacket", Category: "Clothing", IsEssential: false, Reason: "Paris will be 52–63°F during your trip", IsChecked: false, SortOrder: 22},
		{ID: 4, Name: "Umbrella", Category: "Clothing", IsEssential: false, Reason: "2 rainy days expected during your trip", IsChecked: false, SortOrder: 23},
	},
}

// GetTrip handles GET /api/trips/:uuid
// TODO (Lesson 2): Replace mock data with a real database lookup.
func GetTrip(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid != mockTripID {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "trip_not_found",
			"message": "No trip found with that ID.",
		})
		return
	}

	c.JSON(http.StatusOK, mockTrip)
}
