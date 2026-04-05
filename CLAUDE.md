# PackSmart - Implementation Guide for Claude Code

## Project Context

PackSmart is a smart packing list generator. Users enter trip details (destination, dates, trip type, companions, activities), and the app generates a personalized packing list driven by a real weather forecast. Lists are saved to MySQL and shareable via UUID URL.

This is a Module 5 implementation project. All design decisions have been made — do not redesign, add features, or deviate from the documented architecture. Full proposal: `docs/PROJECT_PROPOSAL.md`.

---

## Quick Reference

| Document | Contents |
|---|---|
| `docs/PROJECT_PROPOSAL.md` | Full project overview, all decisions in one place |
| `docs/DATABASE_SCHEMA.md` | 3 tables with CREATE TABLE SQL, relationships, index strategy |
| `docs/API_ENDPOINTS.md` | 3 endpoints with request/response examples and error formats |
| `docs/RULE_ENGINE.md` | Rule groups, execution order, category organization |
| `docs/DECISIONS.md` | Technical decision log with rationale |
| `docs/PAGES.md` | Page inventory with states and arrival paths |
| `docs/COMPONENTS.md` | Component tree and descriptions |

---

## Tech Stack

**Backend**
- Language: Go
- Database: MySQL
- Query layer: sqlc (generate type-safe Go from SQL)
- HTTP router: standard `net/http` or `chi`
- UUID generation: `github.com/google/uuid`

**Frontend**
- Framework: React + Vite
- Routing: React Router v6
- Server state: TanStack Query (React Query)
- HTTP client: fetch (built-in, wrapped in query functions)

**Deployment**
- Server: AWS Lightsail
- Reverse proxy: nginx
- TLS: HTTPS via Let's Encrypt

---

## Implementation Order

Build in this sequence — each layer depends on the one before it.

```
1. Database       schema.sql → MySQL tables
2. sqlc           queries.sql → generated Go types
3. Backend API    Go handlers for 3 endpoints
4. Rule engine    Go module: rules/engine.go
5. Frontend       React app: routing → pages → components → API wiring
6. Deployment     Lightsail → nginx → HTTPS
```

Do not start the frontend until the backend API returns correct responses. Test each endpoint with curl before moving on.

---

## Backend Implementation

### Step 1: Database

Create `backend/schema.sql` with the three tables from `docs/DATABASE_SCHEMA.md`:

```sql
CREATE TABLE trips ( ... );
CREATE TABLE packing_items ( ... );
CREATE TABLE weather_snapshots ( ... );
```

Key details:
- `trips.id` is `CHAR(36)` UUID, not AUTO_INCREMENT
- `packing_items.trip_id` has `ON DELETE CASCADE` and `INDEX idx_trip_id`
- `weather_snapshots.trip_id` is `UNIQUE` (enforces one-to-one)

Run against local MySQL to verify before writing any Go.

### Step 2: sqlc Setup

Create `backend/sqlc.yaml` pointing at `schema.sql` and `queries.sql`. Write queries for:
- `InsertTrip`
- `GetTripByID`
- `InsertPackingItem` (used in a loop)
- `InsertWeatherSnapshot`
- `GetItemsByTripID`
- `GetWeatherByTripID`
- `UpdateItemChecked`

Generate with `sqlc generate`. Do not hand-write database structs.

### Step 3: API Handlers

Implement three handlers. See `docs/API_ENDPOINTS.md` for full request/response shapes.

**POST /api/trips** — the complex one, runs in this order:
1. Decode + validate JSON body (return `400 invalid_request` on bad fields)
2. Call Open-Meteo geocoding API with destination string
   - If zero results returned: return `400 destination_not_found` immediately — do not continue
3. Call Open-Meteo forecast API with lat/lon
   - If trip is >16 days out OR API call fails: set `is_forecast = false`, `weather = null` — continue
4. Call rule engine with trip context + weather (weather may be nil)
5. Write to MySQL: trips row → weather_snapshots row → packing_items rows (in a transaction)
6. Return `201` with complete trip object

**GET /api/trips/{uuid}**
- Look up trip by UUID; return `404 trip_not_found` if missing
- Join weather_snapshots and all packing_items
- Return `200` with same shape as POST response

**PATCH /api/trips/{uuid}/items/{itemId}**
- Verify trip UUID exists; verify item belongs to that trip
- Update `is_checked` only
- Return `200` with `{id, is_checked}`

**CORS:** Add middleware to allow `http://localhost:5173` (dev) and your production domain. Required or the frontend cannot talk to the backend.

**Error format** (all endpoints):
```json
{ "error": "machine_readable_code", "message": "Human readable." }
```

### Step 4: Rule Engine

Create `backend/rules/engine.go`. The engine accepts a context struct and returns `[]PackingItem`.

```go
type TripContext struct {
    TripType     string
    Companions   string
    Activities   []string
    DurationDays int
    Weather      *WeatherContext // nil when is_forecast: false
}
```

Execute rule groups in this order:
1. Base list (toiletries + clothing quantities)
2. Trip type rules (add items; `staying_with_friends` also removes — see removal mechanism below)
3. Weather rules — **skip entirely if `Weather == nil`**
4. Activity rules
5. Companions rules (quantity scaling)
6. Dedupe by name (case-insensitive)
7. Mark essentials (`is_essential = true` on qualifying items)
8. Assign `sort_order` by category

**Removal mechanism for `staying_with_friends`:**
After adding items in step 2, run a filter pass that removes items by name: `["shampoo", "conditioner", "body wash"]`. This is the only trip type that removes items. Implement as a named removal list, not a flag on individual items.

**sort_order ranges by category:**
- Essential Items: 1–9
- Documents & Money: 10–19
- Clothing: 20–39
- Toiletries: 40–54
- Health & Safety: 55–64
- Electronics: 65–74
- Activity Specific: 75–99

**Ski trips:** `trip_type = "cold_weather"` + `activities = ["skiing"]`. The cold_weather block adds base cold gear; the skiing activity block adds ski-specific items on top.

---

## Frontend Implementation

### Step 1: Project Setup

```bash
npm create vite@latest frontend -- --template react
cd frontend
npm install react-router-dom @tanstack/react-query
```

### Step 2: Routing

In `src/main.jsx`, wrap the app in `QueryClientProvider` and `BrowserRouter`. Define routes in `src/App.jsx`:

```jsx
<Routes>
  <Route path="/" element={<HomePage />} />
  <Route path="/packing-list/:tripId" element={<PackingListPage />} />
  <Route path="*" element={<NotFoundPage />} />
</Routes>
```

### Step 3: API Layer

Create `src/api/trips.js` with three functions that wrap fetch calls:

```js
createTrip(formData)         // POST /api/trips
getTrip(tripId)              // GET /api/trips/{uuid}
updateItemChecked(tripId, itemId, isChecked)  // PATCH /api/trips/{uuid}/items/{itemId}
```

These are passed to TanStack Query hooks — keep them pure fetch functions with no React in them.

### Step 4: Home Page (`/`)

`TripForm` component handles all form state. On submit:
1. Validate required fields — show inline errors, do not submit
2. Fire `useMutation` → `createTrip(formData)`
3. Show `<LoadingState message="Building your list…" />` during mutation
4. On success: `navigate(\`/packing-list/${data.id}\`)`
5. On error: show `<ErrorMessage>` banner, stay on page

Form fields: destination (text), departure date, return date, trip type (card grid, single select), companions (card grid, single select), activities (tag multi-select).

### Step 5: Packing List Page (`/packing-list/:tripId`)

On mount, call `useQuery` → `getTrip(tripId)`.

States to handle:
- **Loading:** show `<LoadingState />`
- **Error / not found:** show `<ErrorMessage>` with link to `/`
- **Success:** render the list

Render order:
1. Trip summary header + "← New Trip" + Print buttons
2. Progress bar (checked count / total)
3. `<WeatherCard />` — or weather unavailable banner if `weather === null`
4. `<EssentialItems />` — items where `is_essential === true`
5. `<PackingCategory />` for each remaining category, in `sort_order` order

Checkbox interaction: call `useMutation` → `updateItemChecked(...)` on every toggle. Update local query cache optimistically so the UI responds instantly without waiting for the PATCH to complete.

### Step 6: Components

Build in dependency order (leaf nodes first):

```
PackingItem → PackingCategory → PackingList page
PackingItem → EssentialItems → PackingList page
WeatherCard (standalone)
LoadingState (standalone)
ErrorMessage (standalone)
Layout (wraps all pages)
```

---

## Key Architectural Patterns

**UUIDs for shareable trips**
Generate UUID in Go before inserting: `id := uuid.New().String()`. Store as `CHAR(36)`. The frontend navigates to `/packing-list/{uuid}` after POST — anyone with the URL can view and check off the list.

**Weather fallback: skip rules, show banner**
When `is_forecast: false`, the API returns `"weather": null`. The rule engine skips Rule Group 3 (weather rules) entirely. The frontend checks for `weather === null` and shows: *"Weather forecast unavailable — showing general recommendations for [trip type]."* Do not attempt to infer weather from other data.

**Rule engine is server-side only**
The frontend never generates packing items. It renders what the API returns, in the order the API returns it (`sort_order` ascending). No client-side sorting, filtering, or rule logic.

**Checkbox state via PATCH**
Every checkbox toggle fires `PATCH /api/trips/{uuid}/items/{itemId}`. State lives in MySQL — not React state, not localStorage. Use TanStack Query's `invalidateQueries` or optimistic updates to keep the UI in sync.

**Geocoding failure is a 400, not a 500**
If Open-Meteo geocoding returns zero results, return `400 destination_not_found` before touching the database. The frontend should surface this as an inline form error near the destination field, not a generic error banner.

---

## Common Pitfalls

**Database**
- Forgetting `ON DELETE CASCADE` on `packing_items.trip_id` — orphaned items will accumulate
- Using `INT AUTO_INCREMENT` for `trips.id` — must be `CHAR(36)` UUID
- Not adding `INDEX idx_trip_id` on `packing_items` — every GET will do a full table scan

**Backend**
- Geocoding before validation — validate the request body first, then geocode
- Not wrapping the 3 DB writes in a transaction — a crash mid-write leaves partial data
- Returning a 500 when Open-Meteo is down — fall back to `is_forecast: false` instead
- Forgetting CORS middleware — the frontend will get blocked with no useful error message

**Rule Engine**
- Running weather rules when `Weather == nil` — always guard with a nil check
- Forgetting the removal pass for `staying_with_friends` — base toiletries will appear
- Duplicate items from overlapping rules — dedupe step must run before sort_order assignment
- `sort_order` not set — frontend renders in insertion order, which is unpredictable

**Frontend**
- Navigating to `/packing-list/{uuid}` before the POST response arrives — wait for `onSuccess`
- Not handling the `weather === null` case in `<WeatherCard />` — will throw a runtime error
- Fetching trip data on every checkbox click — use optimistic updates, not refetch
- Forgetting to install React Router — routing won't work without it
