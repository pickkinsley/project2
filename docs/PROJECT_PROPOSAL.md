# Project Proposal — PackSmart

## Overview

### Target Audience
Anyone who travels — from weekend visitors to international travelers. No technical skill required; the app should work for anyone who can fill out a form.

### Problem Statement
People don't know what to pack for trips. They either forget important items or overpack out of uncertainty. Generic packing lists don't account for actual weather, trip type, or who they're traveling with.

### Value Proposition
PackSmart generates a personalized packing list in seconds by combining your trip details (destination, dates, activities, companions) with the real weather forecast. Instead of "pack a jacket," it tells you *why* — "Pack layers — Paris will be 45–58°F with 2 rainy days during your trip."

---

## Feature Scope

### Must-Have Features (MVP)

1. **Trip input form** — destination, dates, travel companions, activities, trip category
2. **Automatic weather forecast** via Open-Meteo (free, no API key) for trips within 16 days; seasonal fallback message for trips further out
3. **Rule-based packing list generation** personalized to weather + trip type + duration + activities
4. **Checklist with checkboxes** to mark items as packed, persisted in MySQL
5. **Trip categories** that meaningfully change the list (international vs. beach vs. cold weather vs. staying with friends) — ski trips use `cold_weather` trip type with `skiing` activity

### Deferred Features (Post-MVP)

1. User accounts and saving multiple trips
2. Sharing lists with travel companions
3. Shopping links for missing items
4. Claude API "Smart Mode" for more nuanced, conversational suggestions
5. Mobile app (iOS/Android)
6. Custom list templates the user can create and reuse

---

## Pages

### Page 1: Home — `/`

**One job:** Collect trip details and send them to the API to generate a packing list.

**How users arrive:**
- Direct visit (first time or returning to plan a new trip)
- Clicking "← New Trip" from the packing list page
- Browser back button from `/packing-list/{uuid}`
- Any unrecognized route (redirect here via 404)

**What's on the page:**
- App name and tagline
- Trip input form: destination, departure/return dates, trip type, companions, activities
- "Generate My Packing List" submit button
- Inline validation errors when required fields are missing
- Loading state ("Building your list…") during API call (2–4 seconds)

---

### Page 2: Packing List — `/packing-list/{uuid}`

**One job:** Show the personalized packing checklist and let the user track what they've packed.

**How users arrive:**
- Form submission on `/` after API responds
- Shared link (anyone with the URL)
- Bookmark or direct return visit

**What's on the page:**
- Trip summary header: destination · duration · trip type
- Progress bar: "X of Y items packed"
- Weather forecast card with daily strip
- ⚠️ Essential Items section (visually distinct, always at top)
- Categorized checklist: Documents & Money, Clothing, Toiletries, Health & Safety, Electronics, Activity Specific
- Each item: checkbox, name, reason line
- Celebration message at 100%

**States:** loading, not found (bad UUID), API error, fully packed

---

### Page 3: 404 — `*`

**One job:** Handle wrong URLs gracefully.

**What's on the page:** "Trip not found" message + link back to `/`

---

### Trip ID Format

All trip IDs are UUIDs (e.g., `a3f8c2d1-4b5e-4c3d-8f9a-1b2c3d4e5f6a`). Without user accounts, sequential integers would let anyone enumerate all trips. UUIDs are unguessable — only people with the link can access a trip.

---

## Navigation Flow

```
                        ┌─────────────────────────────────┐
                        │         ENTRY POINTS            │
                        └─────────────────────────────────┘
                               │              │
                    Direct visit /      Shared link
                    or unknown URL    /packing-list/{uuid}
                               │              │
                               ▼              ▼
┌──────────────────────────────────────────────────────────────────┐
│                         HOME  /                                  │
│                                                                  │
│  [ Destination        ] [ Departure ] [ Return ]                 │
│  [ Trip Type          ] [ Companions          ]                  │
│  [ Activities (multi) ]                                          │
│                                                                  │
│          [ Generate My Packing List ]                            │
└──────────────────────────────────────────────────────────────────┘
         │                          │
    Validation                 All fields
    fails                      valid
         │                          │
         ▼                          ▼
   Inline errors          Loading state (2-4s)
   stay on page           "Building your list…"
                          POST /api/trips
                               │
                    ┌──────────┴──────────┐
                  API                   API
                  error                 success
                    │                     │
                    ▼                     ▼
             Error banner         Navigate to
             stays on /       /packing-list/{uuid}
             user can retry
                                          │
                                          ▼
┌──────────────────────────────────────────────────────────────────┐
│              PACKING LIST  /packing-list/{uuid}                  │
│                                                                  │
│  Paris, France · 5 days · International          [← New Trip]   │
│  ████████████░░░░░░░░  32 of 47 items packed     [🖨 Print]     │
│                                                                  │
│  ⚠️  ESSENTIAL ITEMS                                             │
│     ☐ Passport   ☐ Prescriptions   ☐ Power Adapter              │
│                                                                  │
│  📋 Documents & Money  ·  👕 Clothing  ·  🧴 Toiletries  ...    │
│                                                                  │
│         [ ☑ checkbox ] → PATCH /api/trips/{uuid}/items/{id}     │
└──────────────────────────────────────────────────────────────────┘
         │                               │
    [← New Trip]                   All items checked
    or back button                        │
         │                               ▼
      HOME /                   "🎉 You're all packed!"
```

### User Journey Summary

| Trigger | From | To |
|---|---|---|
| First visit | — | `/` |
| Form submit (valid) | `/` | Loading → `/packing-list/{uuid}` |
| Form submit (invalid) | `/` | Stay on `/`, inline errors |
| API error on submit | `/` | Stay on `/`, error banner |
| "← New Trip" / back button | `/packing-list/{uuid}` | `/` |
| Shared UUID link | — | `/packing-list/{uuid}` |
| Bad UUID | — | 404 → `/` |
| Check a checkbox | `/packing-list/{uuid}` | Stay on page, PATCH to API |

---

## Database Schema

### Tech Stack
- **Frontend:** React + Vite
- **Backend:** Go API
- **Database:** MySQL

### Tables

**`trips`**
```sql
CREATE TABLE trips (
  id             CHAR(36)     PRIMARY KEY,
  destination    VARCHAR(255) NOT NULL,
  dest_lat       DECIMAL(9,6) NOT NULL,
  dest_lon       DECIMAL(9,6) NOT NULL,
  departure_date DATE         NOT NULL,
  return_date    DATE         NOT NULL,
  trip_type      VARCHAR(50)  NOT NULL,
  companions     VARCHAR(50)  NOT NULL,
  activities     JSON         NOT NULL,
  created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**`packing_items`**
```sql
CREATE TABLE packing_items (
  id           INT          PRIMARY KEY AUTO_INCREMENT,
  trip_id      CHAR(36)     NOT NULL,
  name         VARCHAR(255) NOT NULL,
  category     VARCHAR(100) NOT NULL,
  is_essential BOOLEAN      NOT NULL DEFAULT FALSE,
  reason       VARCHAR(500),
  is_checked   BOOLEAN      NOT NULL DEFAULT FALSE,
  sort_order   INT          NOT NULL DEFAULT 0,
  created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
  INDEX idx_trip_id (trip_id)
);
```

**`weather_snapshots`**
```sql
CREATE TABLE weather_snapshots (
  id             INT          PRIMARY KEY AUTO_INCREMENT,
  trip_id        CHAR(36)     NOT NULL UNIQUE,
  temp_min_f     INT          NOT NULL,
  temp_max_f     INT          NOT NULL,
  rain_days      INT          NOT NULL DEFAULT 0,
  snow_days      INT          NOT NULL DEFAULT 0,
  is_forecast    BOOLEAN      NOT NULL,
  daily_forecast JSON,
  fetched_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
);
```

### Relationships

```
trips (UUID PK)
  ├──< packing_items  (one-to-many, FK: trip_id)
  └──  weather_snapshots (one-to-one, FK: trip_id UNIQUE)
```

### Key Design Decisions

- **UUID for trip ID** — unguessable, prevents enumeration without user accounts
- **INT for item ID** — items are never accessed by URL, auto-increment is simpler
- **JSON for activities** — fixed-vocabulary multi-select; simpler than a junction table for MVP
- **Store weather as a snapshot** — fast loads, no Open-Meteo dependency on return visits
- **VARCHAR for category** — rule engine enforces consistency in code; avoids schema migrations to add categories
- **Geocode once, store lat/lon** — no repeated geocoding API calls on subsequent loads

---

## API Endpoints

### Endpoint Summary

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/trips` | Create trip, generate weather + packing list |
| `GET` | `/api/trips/{uuid}` | Fetch complete trip with weather and items |
| `PATCH` | `/api/trips/{uuid}/items/{itemId}` | Check or uncheck a packing item |

---

### POST /api/trips

Creates a trip. Backend geocodes destination, fetches weather, runs rule engine, writes to MySQL, returns complete trip.

**Request:**
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

**Response `201`:** Full trip object with `weather` and `items` arrays.

**If Open-Meteo fails or trip is >16 days out:** Sets `is_forecast: false`, skips weather rules, returns `"weather": null`. Frontend shows: "Weather forecast unavailable — showing general recommendations for [trip type]."

**Tables written:** `trips`, `weather_snapshots`, `packing_items`

---

### GET /api/trips/{uuid}

Fetches a complete existing trip for the packing list page.

**Response `200`:** Same structure as POST response, with current `is_checked` values.

**Response `404`:**
```json
{ "error": "trip_not_found", "message": "No trip found with that ID." }
```

**Tables read:** `trips`, `weather_snapshots`, `packing_items`

---

### PATCH /api/trips/{uuid}/items/{itemId}

Updates the checked state of one packing item.

**Request:** `{ "is_checked": true }`

**Response `200`:** `{ "id": 7, "is_checked": true }`

Verifies item belongs to the trip before updating (prevents cross-trip tampering).

**Tables:** `trips` (verify), `packing_items` (update)

---

### Error Format (all endpoints)

```json
{ "error": "machine_readable_code", "message": "Human readable explanation." }
```

Status codes: `201` created · `200` ok · `400` bad request · `404` not found · `500` server error

---

## Rule Engine

### How It Works

The rule engine runs on the backend during `POST /api/trips`. It takes a context object (trip details + weather) and outputs a deduplicated, sorted list of packing items. The frontend only renders — it never runs rules.

### Execution Order

```
1. Base list       → always-included items (toiletries, clothing quantities)
2. Trip type       → add/remove items by trip type
3. Weather         → clothing + rain/snow gear from forecast
4. Activities      → activity-specific items
5. Companions      → scale quantities
6. Dedupe          → merge duplicate item names
7. Mark essentials → flag is_essential: true on qualifying items
8. Sort order      → sequence items within category
```

### Key Rule Patterns

**Weather — temperature thresholds (uses `temp_min_f`):**

| temp_min_f | Items added |
|---|---|
| < 32°F | Heavy coat, thermals, gloves, wool socks, hat, scarf |
| 32–45°F | Warm jacket, long pants, scarf |
| 46–60°F | Light jacket, mix of pants |
| > 60°F | Light clothing, shorts ok |

**Weather — precipitation:**

| Condition | Items added |
|---|---|
| rain_days ≥ 1 | Umbrella, packable rain jacket |
| rain_days ≥ 3 | Waterproof shoes |
| snow_days ≥ 1 | Waterproof boots, heavy layers |

**Trip type blocks:**

| Trip type | Key additions | Removals |
|---|---|---|
| `international` | Passport ⚠️, power adapter ⚠️, travel insurance ⚠️, foreign currency | — |
| `beach` | Swimsuit ×2, sunscreen ⚠️, flip flops, beach towel | — |
| `cold_weather` | Heavy coat, thermals, gloves, hat, hand warmers | — |
| `staying_with_friends` | Host gift, thank-you note | Shampoo, conditioner, body wash |

**Activity additions:** sightseeing → walking shoes + daypack · fine_dining → dressy outfit · hiking → boots + water bottle · business → business attire + laptop · skiing → ski layers + goggles

**Essential items** (is_essential: true — forgetting ruins the trip):
Passport (international), Prescriptions (always), Power adapter (international), Travel insurance docs (international), Sunscreen (beach)

### Category Order

| # | Category | sort_order |
|---|---|---|
| 1 | Essential Items | 1–9 |
| 2 | Documents & Money | 10–19 |
| 3 | Clothing | 20–39 |
| 4 | Toiletries | 40–54 |
| 5 | Health & Safety | 55–64 |
| 6 | Electronics | 65–74 |
| 7 | Activity Specific | 75–99 |
