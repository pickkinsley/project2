# Technical Decisions — PackSmart

## Weather API: Open-Meteo

**Decision:** Use Open-Meteo for weather forecasts and geocoding.

**Reasons:**
- Completely free with no API key required
- Supports up to 16-day forecasts
- Includes a geocoding API to convert destination names to coordinates
- No rate limits for reasonable usage

**Tradeoff:** Forecasts only go 16 days out. Trips beyond that will show a "forecast not yet available" message and fall back to a seasonal estimate based on the departure month and trip type.

---

## List Generation: Rule-Based Logic

**Decision:** Generate packing lists using a hand-crafted rule engine, not an AI API.

**Reasons:**
- No ongoing API costs — the app works for free indefinitely
- Instant results with no network latency on the generation step
- Predictable, testable output
- A well-designed rule engine can still be smart: it considers weather, trip duration, trip type, activities, and travel companions

**Key rules include:**
- Temperature thresholds driving clothing suggestions (e.g., <45°F → heavy coat, thermals)
- Rain days triggering umbrella / rain jacket recommendations
- Trip type changing entire categories (international → passport, visa; staying with friends → host gift, minimal toiletries)
- Duration scaling clothing quantities (capped at 7 to prompt laundry on longer trips)

**Tradeoff:** Less flexible than AI. Edge cases (e.g., "Buddhist temple visit in hot weather") won't be caught. Addressed by the upgrade path below.

---

## Tech Stack: React + Vite

**Decision:** Build the frontend with React and Vite.

**Reasons:**
- React component model maps well to the UI (form → weather card → category sections → checklist items)
- Vite provides fast dev server and simple build pipeline with no configuration overhead
- No backend needed for MVP — all logic runs in the browser
- localStorage handles state persistence without a database

---

## Future Upgrade Path: Claude API "Smart Mode"

**Decision:** Design the rule engine as a standalone module so Claude API can be swapped in later as an optional upgrade.

**Plan:** After MVP ships, add a "Smart Mode" toggle that sends trip details + weather data to Claude and returns a richer, more personalized list with conversational reasoning. The rule-based engine remains the default (free, fast); Smart Mode becomes a premium or opt-in feature.

**Why deferred:** Adds per-request cost, latency, and API key management — none of which are necessary to validate the core value proposition.
