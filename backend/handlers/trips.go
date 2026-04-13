package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	dbpkg "github.com/pickkinsley/project2/backend/db"
	"github.com/pickkinsley/project2/backend/models"
)

// Handler holds the database connection and query interface.
type Handler struct {
	db *sql.DB
	q  *dbpkg.Queries
}

// NewHandler creates a Handler with an open database connection.
func NewHandler(sqlDB *sql.DB, q *dbpkg.Queries) *Handler {
	return &Handler{db: sqlDB, q: q}
}

// CreateTrip handles POST /api/trips.
// TODO (Lesson 2 — geocoding): Replace DestLat/DestLon "0.000000" with real Open-Meteo geocoding.
// TODO (Lesson 2 — weather): Replace mockWeather() with real Open-Meteo forecast.
// TODO (Lesson 2 — rules): Replace mockItems() with real rule engine output.
func (h *Handler) CreateTrip(c *gin.Context) {
	log.Printf("[INFO] POST /api/trips - Creating new trip")

	var req models.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] POST /api/trips - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	departure, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		log.Printf("[ERROR] POST /api/trips - Invalid departure_date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "departure_date must be in YYYY-MM-DD format.",
		})
		return
	}
	returnDate, err := time.Parse("2006-01-02", req.ReturnDate)
	if err != nil {
		log.Printf("[ERROR] POST /api/trips - Invalid return_date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "return_date must be in YYYY-MM-DD format.",
		})
		return
	}
	if !returnDate.After(departure) {
		log.Printf("[ERROR] POST /api/trips - return_date not after departure_date")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "return_date must be after departure_date.",
		})
		return
	}
	durationDays := int(returnDate.Sub(departure).Hours()/24) + 1

	tripID := uuid.New().String()
	ctx := context.Background()

	log.Printf("[DEBUG] POST /api/trips - Destination: %q, TripType: %q, Duration: %d days", req.Destination, req.TripType, durationDays)

	activitiesJSON, err := json.Marshal(req.Activities)
	if err != nil {
		log.Printf("[ERROR] POST /api/trips - Failed to encode activities: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to encode activities."})
		return
	}

	weather := mockWeather()
	items := mockItems(req.TripType)

	forecastJSON, err := json.Marshal(weather.DailyForecast)
	if err != nil {
		log.Printf("[ERROR] POST /api/trips - Failed to encode forecast: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to encode forecast."})
		return
	}

	// Write trip, weather, and items in a single transaction.
	log.Printf("[DEBUG] POST /api/trips - Beginning transaction for trip %s", tripID)
	tx, err := h.db.Begin()
	if err != nil {
		log.Printf("[ERROR] POST /api/trips - Failed to start transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to start transaction."})
		return
	}
	qtx := h.q.WithTx(tx)

	if err := qtx.InsertTrip(ctx, dbpkg.InsertTripParams{
		ID:            tripID,
		Destination:   req.Destination,
		DestLat:       "0.000000", // TODO: replace with geocoded lat
		DestLon:       "0.000000", // TODO: replace with geocoded lon
		DepartureDate: departure,
		ReturnDate:    returnDate,
		TripType:      req.TripType,
		Companions:    req.Companions,
		Activities:    json.RawMessage(activitiesJSON),
	}); err != nil {
		tx.Rollback()
		log.Printf("[ERROR] POST /api/trips - InsertTrip failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to save trip."})
		return
	}
	log.Printf("[DEBUG] POST /api/trips - Trip row inserted")

	if err := qtx.InsertWeatherSnapshot(ctx, dbpkg.InsertWeatherSnapshotParams{
		TripID:        tripID,
		TempMinF:      int32(weather.TempMinF),
		TempMaxF:      int32(weather.TempMaxF),
		RainDays:      int32(weather.RainDays),
		SnowDays:      int32(weather.SnowDays),
		IsForecast:    weather.IsForecast,
		DailyForecast: json.RawMessage(forecastJSON),
	}); err != nil {
		tx.Rollback()
		log.Printf("[ERROR] POST /api/trips - InsertWeatherSnapshot failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to save weather."})
		return
	}
	log.Printf("[DEBUG] POST /api/trips - Weather snapshot inserted")

	for _, item := range items {
		if err := qtx.InsertPackingItem(ctx, dbpkg.InsertPackingItemParams{
			TripID:      tripID,
			Name:        item.Name,
			Category:    item.Category,
			IsEssential: item.IsEssential,
			Reason:      sql.NullString{String: item.Reason, Valid: item.Reason != ""},
			IsChecked:   false,
			SortOrder:   int32(item.SortOrder),
		}); err != nil {
			tx.Rollback()
			log.Printf("[ERROR] POST /api/trips - InsertPackingItem %q failed: %v", item.Name, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to save packing items."})
			return
		}
	}
	log.Printf("[DEBUG] POST /api/trips - %d packing items inserted", len(items))

	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR] POST /api/trips - Transaction commit failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to commit transaction."})
		return
	}

	// Fetch items back to get their auto-assigned IDs.
	dbItems, err := h.q.GetItemsByTripID(ctx, tripID)
	if err != nil {
		log.Printf("[ERROR] POST /api/trips - GetItemsByTripID failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to retrieve packing items."})
		return
	}

	log.Printf("[INFO] Trip created with ID: %s", tripID)
	log.Printf("[INFO] POST /api/trips - Returning 201")

	c.JSON(http.StatusCreated, models.TripResponse{
		ID:            tripID,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		TripType:      req.TripType,
		Companions:    req.Companions,
		Activities:    req.Activities,
		DurationDays:  durationDays,
		CreatedAt:     time.Now().UTC(),
		Weather:       weather,
		Items:         dbItemsToResponse(dbItems),
	})
}

// GetTrip handles GET /api/trips/:uuid.
func (h *Handler) GetTrip(c *gin.Context) {
	tripUUID := c.Param("uuid")
	log.Printf("[INFO] GET /api/trips/%s - Fetching trip", tripUUID)
	ctx := context.Background()

	trip, err := h.q.GetTripByID(ctx, tripUUID)
	if err == sql.ErrNoRows {
		log.Printf("[INFO] GET /api/trips/%s - Trip not found, returning 404", tripUUID)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "trip_not_found",
			"message": "No trip found with that ID.",
		})
		return
	}
	if err != nil {
		log.Printf("[ERROR] GET /api/trips/%s - GetTripByID failed: %v", tripUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to retrieve trip."})
		return
	}
	log.Printf("[DEBUG] GET /api/trips/%s - Trip found: %q", tripUUID, trip.Destination)

	var activities []string
	if err := json.Unmarshal(trip.Activities, &activities); err != nil {
		activities = []string{}
	}

	// Weather is optional — nil if not found (is_forecast: false trips won't have one in the future).
	var weatherResp *models.WeatherResponse
	ws, err := h.q.GetWeatherByTripID(ctx, tripUUID)
	if err == nil {
		log.Printf("[DEBUG] GET /api/trips/%s - Weather snapshot found", tripUUID)
		var forecast []models.DailyForecast
		json.Unmarshal(ws.DailyForecast, &forecast)
		weatherResp = &models.WeatherResponse{
			TempMinF:      int(ws.TempMinF),
			TempMaxF:      int(ws.TempMaxF),
			RainDays:      int(ws.RainDays),
			SnowDays:      int(ws.SnowDays),
			IsForecast:    ws.IsForecast,
			DailyForecast: forecast,
		}
	} else if err != sql.ErrNoRows {
		log.Printf("[ERROR] GET /api/trips/%s - GetWeatherByTripID failed: %v", tripUUID, err)
	}

	dbItems, err := h.q.GetItemsByTripID(ctx, tripUUID)
	if err != nil {
		log.Printf("[ERROR] GET /api/trips/%s - GetItemsByTripID failed: %v", tripUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to retrieve packing items."})
		return
	}
	log.Printf("[DEBUG] GET /api/trips/%s - Fetched %d packing items", tripUUID, len(dbItems))

	durationDays := int(trip.ReturnDate.Sub(trip.DepartureDate).Hours()/24) + 1

	log.Printf("[INFO] GET /api/trips/%s - Returning 200", tripUUID)
	c.JSON(http.StatusOK, models.TripResponse{
		ID:            trip.ID,
		Destination:   trip.Destination,
		DepartureDate: trip.DepartureDate.Format("2006-01-02"),
		ReturnDate:    trip.ReturnDate.Format("2006-01-02"),
		TripType:      trip.TripType,
		Companions:    trip.Companions,
		Activities:    activities,
		DurationDays:  durationDays,
		CreatedAt:     trip.CreatedAt,
		Weather:       weatherResp,
		Items:         dbItemsToResponse(dbItems),
	})
}

// UpdateItemCheckbox handles PATCH /api/trips/:uuid/items/:itemId.
func (h *Handler) UpdateItemCheckbox(c *gin.Context) {
	tripUUID := c.Param("uuid")
	itemIDStr := c.Param("itemId")
	log.Printf("[INFO] PATCH /api/trips/%s/items/%s - Updating item checkbox", tripUUID, itemIDStr)
	ctx := context.Background()

	// Verify trip exists.
	if _, err := h.q.GetTripByID(ctx, tripUUID); err == sql.ErrNoRows {
		log.Printf("[INFO] PATCH /api/trips/%s/items/%s - Trip not found, returning 404", tripUUID, itemIDStr)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "trip_not_found",
			"message": "No trip found with that ID.",
		})
		return
	} else if err != nil {
		log.Printf("[ERROR] PATCH /api/trips/%s/items/%s - GetTripByID failed: %v", tripUUID, itemIDStr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to verify trip."})
		return
	}

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil || itemID < 1 {
		log.Printf("[ERROR] PATCH /api/trips/%s/items/%s - Invalid item ID", tripUUID, itemIDStr)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "item_not_found",
			"message": "No item found with that ID for this trip.",
		})
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] PATCH /api/trips/%s/items/%d - Invalid request body: %v", tripUUID, itemID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Request body must be JSON with an is_checked boolean field.",
		})
		return
	}

	if err := h.q.UpdateItemChecked(ctx, dbpkg.UpdateItemCheckedParams{
		IsChecked: req.IsChecked,
		ID:        int32(itemID),
		TripID:    tripUUID,
	}); err != nil {
		log.Printf("[ERROR] PATCH /api/trips/%s/items/%d - UpdateItemChecked failed: %v", tripUUID, itemID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Failed to update item."})
		return
	}

	log.Printf("[INFO] PATCH /api/trips/%s/items/%d - is_checked set to %v, returning 200", tripUUID, itemID, req.IsChecked)
	c.JSON(http.StatusOK, models.UpdateItemResponse{
		ID:        itemID,
		IsChecked: req.IsChecked,
	})
}

// dbItemsToResponse converts sqlc PackingItem rows to the API response type.
func dbItemsToResponse(rows []dbpkg.PackingItem) []models.PackingItem {
	out := make([]models.PackingItem, len(rows))
	for i, r := range rows {
		out[i] = models.PackingItem{
			ID:          int(r.ID),
			Name:        r.Name,
			Category:    r.Category,
			IsEssential: r.IsEssential,
			Reason:      r.Reason.String,
			IsChecked:   r.IsChecked,
			SortOrder:   int(r.SortOrder),
		}
	}
	return out
}

// mockWeather returns a fixed weather snapshot until Open-Meteo integration is added.
// TODO (Lesson 2 — weather): Replace with real Open-Meteo forecast data.
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

// mockItems returns a small packing list based on trip type until the rule engine is added.
// TODO (Lesson 2 — rules): Replace with real rule engine output.
func mockItems(tripType string) []models.PackingItem {
	items := []models.PackingItem{
		{Name: "Prescriptions", Category: "Essential Items", IsEssential: true, Reason: "Never travel without your medications", SortOrder: 1},
		{Name: "Phone charger", Category: "Electronics", IsEssential: false, Reason: "Keep your phone powered", SortOrder: 65},
		{Name: "Toothbrush + toothpaste", Category: "Toiletries", IsEssential: false, Reason: "Daily essential", SortOrder: 40},
	}

	switch tripType {
	case "international":
		items = append([]models.PackingItem{
			{Name: "Passport", Category: "Essential Items", IsEssential: true, Reason: "Required for international travel", SortOrder: 0},
			{Name: "Power adapter", Category: "Essential Items", IsEssential: true, Reason: "Different outlets abroad", SortOrder: 2},
		}, items...)
	case "beach":
		items = append(items,
			models.PackingItem{Name: "Swimsuit", Category: "Clothing", IsEssential: false, Reason: "Beach trip essential", SortOrder: 20},
			models.PackingItem{Name: "Sunscreen (SPF 30+)", Category: "Essential Items", IsEssential: true, Reason: "Protect your skin at the beach", SortOrder: 3},
		)
	case "cold_weather":
		items = append(items,
			models.PackingItem{Name: "Heavy coat", Category: "Clothing", IsEssential: false, Reason: "Cold weather essential", SortOrder: 20},
			models.PackingItem{Name: "Gloves", Category: "Clothing", IsEssential: false, Reason: "Keep hands warm", SortOrder: 21},
		)
	}

	return items
}
