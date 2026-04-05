# Database Schema — PackSmart

## Overview

PackSmart uses 3 MySQL tables. Trips are the central resource — weather snapshots and packing items both belong to a trip.

```
trips ──< packing_items     (one-to-many)
trips ──  weather_snapshots (one-to-one)
```

---

## Table 1: `trips`

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

**Column notes:**

| Column | Notes |
|---|---|
| `id` | UUID stored as CHAR(36), e.g. `"a3f8c2d1-4b5e-4c3d-8f9a-1b2c3d4e5f6a"` |
| `dest_lat` / `dest_lon` | Geocoded once at creation, reused on all subsequent loads |
| `trip_type` | One of: `international`, `beach`, `cold_weather`, `weekend_getaway`, `staying_with_friends` |
| `companions` | One of: `solo`, `couple`, `family`, `friends` |
| `activities` | JSON array of strings — see structure below |

**Activities JSON structure:**
```json
["sightseeing", "fine_dining", "hiking", "business", "swimming", "skiing"]
```
Valid values: `sightseeing`, `fine_dining`, `hiking`, `business`, `swimming`, `skiing`

---

## Table 2: `packing_items`

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

**Column notes:**

| Column | Notes |
|---|---|
| `id` | Auto-increment INT — items are only accessed within a known trip, never by URL |
| `category` | VARCHAR(100) — rule engine enforces valid values in code, not the DB |
| `is_essential` | TRUE only for passport, prescriptions, power adapter, travel insurance, sunscreen (beach) |
| `reason` | Human-readable explanation, e.g. "Paris will be 52–63°F with 2 rainy days" |
| `sort_order` | Controls display order — frontend renders items in this order, no client sorting needed |

**Category values and sort_order ranges:**

| Category | sort_order range |
|---|---|
| Essential Items | 1–9 |
| Documents & Money | 10–19 |
| Clothing | 20–39 |
| Toiletries | 40–54 |
| Health & Safety | 55–64 |
| Electronics | 65–74 |
| Activity Specific | 75–99 |

**`ON DELETE CASCADE`:** Deleting a trip automatically deletes all its packing items.

---

## Table 3: `weather_snapshots`

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

**Column notes:**

| Column | Notes |
|---|---|
| `trip_id` | UNIQUE enforces one-to-one relationship at the DB level |
| `is_forecast` | TRUE = real forecast from Open-Meteo; FALSE = seasonal estimate |
| `daily_forecast` | JSON array, null if seasonal estimate |

**Daily forecast JSON structure:**
```json
[
  { "date": "2026-04-10", "icon": "partly_cloudy", "min_f": 50, "max_f": 61 },
  { "date": "2026-04-11", "icon": "rainy",          "min_f": 48, "max_f": 55 }
]
```

**Icon values:** `sunny`, `partly_cloudy`, `cloudy`, `rainy`, `snowy`, `stormy`

---

## Relationships Diagram

```
trips
  id (PK, CHAR(36), UUID)
  destination
  dest_lat / dest_lon
  departure_date / return_date
  trip_type
  companions
  activities (JSON)
  created_at
       │
       ├──────────────────────────────< packing_items
       │                                 id (PK, INT, auto-increment)
       │                                 trip_id (FK) ──────────────→ trips.id
       │                                 name
       │                                 category
       │                                 is_essential
       │                                 reason
       │                                 is_checked
       │                                 sort_order
       │
       └────────────────────────────── weather_snapshots
                                         id (PK, INT, auto-increment)
                                         trip_id (FK, UNIQUE) ───────→ trips.id
                                         temp_min_f / temp_max_f
                                         rain_days / snow_days
                                         is_forecast
                                         daily_forecast (JSON)
```

---

## Index Strategy

| Table | Index | Reason |
|---|---|---|
| `trips` | PRIMARY KEY on `id` | UUID lookup by packing list URL |
| `packing_items` | PRIMARY KEY on `id` | Item lookup for PATCH requests |
| `packing_items` | `idx_trip_id` on `trip_id` | Fetch all items for a trip (most common query) |
| `weather_snapshots` | UNIQUE on `trip_id` | Enforces one-to-one + fast join |

**Not indexed:** `trips.created_at` — no query filters or sorts by this in MVP. Add when trip history is built.

---

## Design Decisions

**UUID for trip ID, INT for item ID**
Trip UUIDs appear in the URL (`/packing-list/{uuid}`). Without user accounts, sequential integers would let anyone enumerate all trips. UUIDs are unguessable. Item IDs never appear in the URL — they're only used internally in PATCH requests within a known trip, so auto-increment INT is simpler and faster.

**JSON for activities**
Activities are a fixed-vocabulary multi-select. A junction table (`trip_activities`) would be normalized but adds a join for every query. JSON in a single column is simpler for MVP and still queryable via MySQL's JSON functions if needed later.

**Store weather as a snapshot**
Subsequent page loads don't re-call Open-Meteo. The forecast won't change between when you generate the list and when you're packing. Also: if Open-Meteo is down, existing trips still load correctly.

**VARCHAR for category, not ENUM**
The rule engine controls what category values get written. Consistency is enforced in application code. VARCHAR avoids a schema migration every time a category is added or renamed.

**Geocode once, store lat/lon**
The frontend sends a destination string ("Paris, France"). The backend geocodes it once at trip creation and stores the coordinates. All subsequent weather lookups use the stored coordinates — no repeated geocoding.
