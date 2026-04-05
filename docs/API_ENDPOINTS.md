# API Endpoints — PackSmart

## Overview

PackSmart has 3 RESTful endpoints. All requests and responses use JSON. The base path for all endpoints is `/api`.

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/trips` | Create trip, generate weather + packing list |
| `GET` | `/api/trips/{uuid}` | Fetch complete trip with weather and items |
| `PATCH` | `/api/trips/{uuid}/items/{itemId}` | Check or uncheck a packing item |

---

## Error Response Format

All error responses use this consistent structure:

```json
{
  "error": "machine_readable_code",
  "message": "Human readable explanation."
}
```

**HTTP status codes:**

| Code | When |
|---|---|
| `201` | Trip created successfully |
| `200` | Successful GET or PATCH |
| `400` | Invalid request body (missing fields, bad date format, return before departure) |
| `404` | Trip UUID or item ID not found |
| `500` | Unexpected server error |

---

## POST /api/trips

**Purpose:** Create a trip. The backend geocodes the destination, fetches weather, runs the rule engine, persists everything to MySQL, and returns the complete trip in one response.

**Request body:**
```json
{
  "destination": "Paris, France",
  "departure_date": "2026-04-10",
  "return_date": "2026-04-15",
  "trip_type": "international",
  "companions": "couple",
  "activities": ["sightseeing", "fine_dining"]
}
```

**Required fields:** all fields above
**Date format:** `YYYY-MM-DD`
**`trip_type` valid values:** `international`, `beach`, `cold_weather`, `weekend_getaway`, `staying_with_friends`
**`companions` valid values:** `solo`, `couple`, `family`, `friends`
**`activities` valid values:** `sightseeing`, `fine_dining`, `hiking`, `business`, `swimming`, `skiing`

**Backend execution order:**
1. Validate request body
2. Geocode destination via Open-Meteo geocoding API → get lat/lon
3. Fetch weather forecast via Open-Meteo (or generate seasonal estimate if >16 days out or API fails)
4. Run rule engine with trip details + weather → generate packing items
5. Write `trips` row to MySQL
6. Write `weather_snapshots` row to MySQL
7. Write all `packing_items` rows to MySQL
8. Return complete trip

**If Open-Meteo fails:** Generate seasonal estimate from destination + departure month. Set `is_forecast: false`. Continue with trip creation — do not fail the request.

**Response `201 Created`:**
```json
{
  "id": "a3f8c2d1-4b5e-4c3d-8f9a-1b2c3d4e5f6a",
  "destination": "Paris, France",
  "departure_date": "2026-04-10",
  "return_date": "2026-04-15",
  "trip_type": "international",
  "companions": "couple",
  "activities": ["sightseeing", "fine_dining"],
  "duration_days": 5,
  "created_at": "2026-04-05T14:23:00Z",
  "weather": {
    "temp_min_f": 52,
    "temp_max_f": 63,
    "rain_days": 2,
    "snow_days": 0,
    "is_forecast": true,
    "daily_forecast": [
      { "date": "2026-04-10", "icon": "partly_cloudy", "min_f": 50, "max_f": 61 },
      { "date": "2026-04-11", "icon": "rainy",          "min_f": 48, "max_f": 55 },
      { "date": "2026-04-12", "icon": "sunny",          "min_f": 52, "max_f": 63 },
      { "date": "2026-04-13", "icon": "partly_cloudy",  "min_f": 51, "max_f": 60 },
      { "date": "2026-04-14", "icon": "rainy",          "min_f": 49, "max_f": 57 }
    ]
  },
  "items": [
    {
      "id": 1,
      "name": "Passport",
      "category": "Essential Items",
      "is_essential": true,
      "reason": "Required for international travel",
      "is_checked": false,
      "sort_order": 1
    },
    {
      "id": 2,
      "name": "Prescriptions",
      "category": "Essential Items",
      "is_essential": true,
      "reason": "Never travel without your medications",
      "is_checked": false,
      "sort_order": 2
    },
    {
      "id": 7,
      "name": "Light jacket",
      "category": "Clothing",
      "is_essential": false,
      "reason": "Paris will be 52–63°F during your trip",
      "is_checked": false,
      "sort_order": 22
    },
    {
      "id": 8,
      "name": "Umbrella",
      "category": "Clothing",
      "is_essential": false,
      "reason": "2 rainy days expected during your trip",
      "is_checked": false,
      "sort_order": 23
    }
  ]
}
```

**Response `400 Bad Request`:**
```json
{
  "error": "invalid_request",
  "message": "return_date must be after departure_date."
}
```

**Tables written:** `trips`, `weather_snapshots`, `packing_items`

---

## GET /api/trips/{uuid}

**Purpose:** Fetch a complete existing trip — trip details, weather snapshot, and all packing items — for the packing list page.

**Request body:** None

**Path parameter:** `uuid` — the trip's UUID

**Backend execution:**
1. Look up trip by UUID in `trips`
2. If not found, return 404
3. Join `weather_snapshots` and all `packing_items` for that trip
4. Return complete trip

**Response `200 OK`:** Same structure as `POST /api/trips` response above, with current `is_checked` values reflecting any updates the user has made.

**Response `404 Not Found`:**
```json
{
  "error": "trip_not_found",
  "message": "No trip found with that ID."
}
```

**Tables read:** `trips`, `weather_snapshots`, `packing_items`

---

## PATCH /api/trips/{uuid}/items/{itemId}

**Purpose:** Check or uncheck a single packing item. Called every time the user taps a checkbox.

**Path parameters:**
- `uuid` — the trip's UUID
- `itemId` — the packing item's integer ID

**Request body:**
```json
{
  "is_checked": true
}
```

**Backend execution:**
1. Verify trip UUID exists in `trips`
2. Verify item exists in `packing_items` AND belongs to that trip (prevents cross-trip tampering)
3. Update `is_checked` on the item
4. Return the updated item

**Response `200 OK`:**
```json
{
  "id": 7,
  "is_checked": true
}
```

**Response `404 Not Found`:**
```json
{
  "error": "item_not_found",
  "message": "No item found with that ID for this trip."
}
```

**Tables read/written:** `trips` (verify), `packing_items` (update)

---

## Table Operations Summary

| Endpoint | trips | weather_snapshots | packing_items |
|---|---|---|---|
| `POST /api/trips` | Write | Write | Write (bulk) |
| `GET /api/trips/{uuid}` | Read | Read | Read (all for trip) |
| `PATCH /api/trips/{uuid}/items/{itemId}` | Read (verify) | — | Write (one row) |

---

## Implementation Notes

**POST /api/trips is the slow endpoint.** It makes 2 external API calls (geocoding + weather) and performs multiple DB writes. Expected response time: 2–4 seconds. The frontend shows a loading state for this duration. Keep it synchronous for MVP — async job queues are not worth the complexity.

**PATCH is the high-frequency endpoint.** Called on every checkbox interaction. Keep it fast: single row read to verify ownership, single row update.

**CORS:** The Go backend needs CORS headers configured to accept requests from the React frontend's origin during development (`http://localhost:5173`) and production domain.
