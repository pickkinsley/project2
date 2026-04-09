package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// CreateTrip handles POST /api/trips.
// TODO (Lesson 2): Replace mock data with geocoding, weather fetch, rule engine, and DB writes.
func CreateTrip(c *gin.Context) {
	var req models.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	// Calculate duration
	departure, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "departure_date must be in YYYY-MM-DD format.",
		})
		return
	}
	returnDate, err := time.Parse("2006-01-02", req.ReturnDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "return_date must be in YYYY-MM-DD format.",
		})
		return
	}
	if !returnDate.After(departure) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "return_date must be after departure_date.",
		})
		return
	}
	durationDays := int(returnDate.Sub(departure).Hours()/24) + 1

	// TODO (Lesson 2): Geocode req.Destination → lat/lon via Open-Meteo geocoding API.
	// TODO (Lesson 2): Fetch weather forecast via Open-Meteo forecast API.
	// TODO (Lesson 2): Run rule engine with trip context + weather.
	// TODO (Lesson 2): Write trip, weather, and items to MySQL in a transaction.

	tripID := uuid.New().String()

	response := models.TripResponse{
		ID:            tripID,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		TripType:      req.TripType,
		Companions:    req.Companions,
		Activities:    req.Activities,
		DurationDays:  durationDays,
		CreatedAt:     time.Now().UTC(),
		Weather:       mockWeather(),
		Items:         mockItems(req.TripType),
	}

	c.JSON(http.StatusCreated, response)
}

// mockWeather returns a fixed weather snapshot for Lesson 1 testing.
// TODO (Lesson 2): Replace with real Open-Meteo forecast data.
func mockWeather() *models.WeatherResponse {
	return &models.WeatherResponse{
		TempMinF:   52,
		TempMaxF:   68,
		RainDays:   2,
		SnowDays:   0,
		IsForecast: true,
		DailyForecast: []models.DailyForecast{
			{Date: "2026-05-01", Icon: "sunny", MinF: 55, MaxF: 68},
			{Date: "2026-05-02", Icon: "partly_cloudy", MinF: 52, MaxF: 65},
			{Date: "2026-05-03", Icon: "rainy", MinF: 50, MaxF: 60},
		},
	}
}

// mockItems returns a small packing list based on trip type for Lesson 1 testing.
// TODO (Lesson 2): Replace with real rule engine output.
func mockItems(tripType string) []models.PackingItem {
	items := []models.PackingItem{
		{ID: 1, Name: "Prescriptions", Category: "Essential Items", IsEssential: true, Reason: "Never travel without your medications", IsChecked: false, SortOrder: 1},
		{ID: 2, Name: "Phone charger", Category: "Electronics", IsEssential: false, Reason: "Keep your phone powered", IsChecked: false, SortOrder: 65},
		{ID: 3, Name: "Toothbrush + toothpaste", Category: "Toiletries", IsEssential: false, Reason: "Daily essential", IsChecked: false, SortOrder: 40},
	}

	switch tripType {
	case "international":
		items = append([]models.PackingItem{
			{ID: 4, Name: "Passport", Category: "Essential Items", IsEssential: true, Reason: "Required for international travel", IsChecked: false, SortOrder: 0},
			{ID: 5, Name: "Power adapter", Category: "Essential Items", IsEssential: true, Reason: "Different outlets abroad", IsChecked: false, SortOrder: 2},
		}, items...)
	case "beach":
		items = append(items,
			models.PackingItem{ID: 4, Name: "Swimsuit", Category: "Clothing", IsEssential: false, Reason: "Beach trip essential", IsChecked: false, SortOrder: 20},
			models.PackingItem{ID: 5, Name: "Sunscreen (SPF 30+)", Category: "Essential Items", IsEssential: true, Reason: "Protect your skin at the beach", IsChecked: false, SortOrder: 3},
		)
	case "cold_weather":
		items = append(items,
			models.PackingItem{ID: 4, Name: "Heavy coat", Category: "Clothing", IsEssential: false, Reason: "Cold weather essential", IsChecked: false, SortOrder: 20},
			models.PackingItem{ID: 5, Name: "Gloves", Category: "Clothing", IsEssential: false, Reason: "Keep hands warm", IsChecked: false, SortOrder: 21},
		)
	}

	return items
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

// UpdateItemCheckbox handles PATCH /api/trips/:uuid/items/:itemId
// TODO (Lesson 2): Verify item belongs to trip and update is_checked in MySQL.
func UpdateItemCheckbox(c *gin.Context) {
	tripUUID := c.Param("uuid")
	itemIDStr := c.Param("itemId")

	// Verify trip exists (mock: only the hardcoded UUID is valid)
	if tripUUID != mockTripID {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "trip_not_found",
			"message": "No trip found with that ID.",
		})
		return
	}

	// Parse item ID
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil || itemID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "item_not_found",
			"message": "No item found with that ID for this trip.",
		})
		return
	}

	// Validate request body
	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Request body must be JSON with an is_checked boolean field.",
		})
		return
	}

	// TODO (Lesson 2): Verify item ID exists in packing_items and belongs to this trip.
	// TODO (Lesson 2): UPDATE packing_items SET is_checked = ? WHERE id = ? AND trip_id = ?

	c.JSON(http.StatusOK, models.UpdateItemResponse{
		ID:        itemID,
		IsChecked: req.IsChecked,
	})
}
