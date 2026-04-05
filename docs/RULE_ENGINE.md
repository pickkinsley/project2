# Rule Engine — PackSmart

## Overview

The rule engine takes a context object (trip details + weather) and outputs a list of packing items. Rules are additive — items from multiple rule groups stack together and are deduped by name before being returned.

The rule engine runs entirely on the backend (Go) during `POST /api/trips`. The frontend never runs rules — it only renders what the backend returns.

---

## Input: Context Object

```
context = {
  trip_type:    "international" | "beach" | "cold_weather" | "weekend_getaway" | "staying_with_friends"
  companions:   "solo" | "couple" | "family" | "friends"
  activities:   ["sightseeing", "fine_dining", ...]
  duration_days: 5
  weather: {
    temp_min_f: 52
    temp_max_f: 63
    rain_days:  2
    snow_days:  0
    is_forecast: true
  }
}
```

---

## Execution Order

Rules run in this sequence. Each step can add items; step 6 removes duplicates.

```
1. Base list          → always-included items (toiletries, basic clothing)
2. Trip type rules    → add (or remove) items based on trip type
3. Weather rules      → add clothing and rain/snow gear
4. Activity rules     → add activity-specific items
5. Companions rules   → scale quantities
6. Dedupe             → merge items with the same name, keep one
7. Mark essentials    → flag qualifying items as is_essential: true
8. Assign sort_order  → sequence items within their category
```

---

## Rule Group 1: Base List

Always included regardless of any other inputs. These are the items every trip needs.

**Toiletries (always):** Toothbrush, toothpaste, deodorant, shampoo, conditioner, body wash, moisturizer, razor

**Clothing (always — quantity scaled by duration):**
```
underwear: MIN(duration_days, 7)
socks:     MIN(duration_days, 7)
shirts:    MIN(duration_days, 7)
```
Capped at 7 — trips longer than a week assume laundry access.

**Health (always):** Prescriptions ⚠️, pain relievers, bandages

---

## Rule Group 2: Trip Type Rules

Each trip type adds a distinct block of items. `staying_with_friends` is the only type that also *removes* items.

**`international`**
- Adds: Passport ⚠️, power adapter ⚠️, travel insurance docs ⚠️, foreign currency, visa reminder, copies of documents, notify-your-bank reminder
- Removes: nothing

**`beach`**
- Adds: Swimsuit (2×), sunscreen ⚠️, flip flops, beach towel, after-sun lotion, rash guard
- Removes: nothing

**`cold_weather`**
- Adds: Heavy coat, thermals (top + bottom), gloves, hat, scarf, wool socks, hand warmers, lip balm
- Removes: nothing
- Note: Stacks with weather rules — cold weather trip in a blizzard gets both

**`weekend_getaway`**
- Adds: nothing beyond base list
- Trims clothing count (shorter trip = fewer changes)

**`staying_with_friends`**
- Adds: Host gift reminder, thank-you note reminder
- Removes: shampoo, conditioner, body wash (host has these) — trims toiletries to: toothbrush, toothpaste, deodorant, razor only

---

## Rule Group 3: Weather Rules

Uses `temp_min_f` as the cold anchor (the coldest it will be determines what to pack).

**Temperature thresholds:**

| `temp_min_f` | Items added |
|---|---|
| < 32°F | Heavy coat, thermals, gloves, wool socks, hat, scarf |
| 32–45°F | Warm jacket, long pants (all), scarf |
| 46–60°F | Light jacket, mix of long and short pants |
| > 60°F | Light clothing, shorts ok — no jacket added |

**Precipitation:**

| Condition | Items added |
|---|---|
| `rain_days >= 1` | Umbrella, packable rain jacket |
| `rain_days >= 3` | Waterproof shoes or rain boots |
| `snow_days >= 1` | Waterproof boots, heavy layers (stacks with cold temp rules) |

**Reasoning strings use real forecast data:**
- "Paris will be 52–63°F during your trip — pack a light jacket"
- "2 rainy days expected — umbrella recommended"
- "Pack a rain jacket — 4 of your 5 days have rain forecast"

**If `is_forecast: false`** (seasonal estimate): reasoning strings use softer language:
- "April in Paris is typically 50–60°F — a light jacket is recommended"

---

## Rule Group 4: Activity Rules

Activities layer on top of the base + trip type list. Multiple activities stack.

| Activity | Items added |
|---|---|
| `sightseeing` | Comfortable walking shoes, daypack, portable charger, offline maps reminder |
| `fine_dining` | Dressy outfit (1×), dress shoes or heels |
| `hiking` | Hiking boots, moisture-wicking socks, daypack, water bottle, trail snacks, sunscreen |
| `business` | Business attire (2×), laptop + charger, business cards |
| `swimming` | Swimsuit (if not already added by beach trip type), goggles |
| `skiing` | Ski layers, ski goggles, helmet reminder, waterproof gloves — stacks with `cold_weather` |

---

## Rule Group 5: Companions Rules

Companions affect quantities, not item types. Everyone packs the same categories — families just pack more of everything.

| Companions | Quantity multiplier |
|---|---|
| `solo` | 1× (baseline) |
| `couple` | No change to items — each person packs their own list |
| `family` | Sunscreen quantity ↑, first aid kit added, snacks added |
| `friends` | No quantity change |

For MVP, the list is generated as a single shared list, not per-person. Families get a note like "pack sunscreen for each person."

---

## Rule Group 7: Essential Items

An item is marked `is_essential: true` only if **forgetting it ruins or prevents the trip.** The threshold is intentionally high.

| Item | Condition |
|---|---|
| Passport | `trip_type === "international"` |
| Prescriptions | Always — every trip |
| Power adapter | `trip_type === "international"` |
| Travel insurance docs | `trip_type === "international"` |
| Sunscreen | `trip_type === "beach"` |

Everything else is important but not essential. A forgotten toothbrush is annoying; a forgotten passport is a cancelled trip.

Essential items are always placed in the `Essential Items` category with `sort_order` 1–9, displayed at the top of the list with distinct visual treatment.

---

## Category Organization

Items are distributed into categories and displayed in this fixed order:

| Order | Category | sort_order range | Example items |
|---|---|---|---|
| 1 | Essential Items | 1–9 | Passport, Prescriptions, Power adapter |
| 2 | Documents & Money | 10–19 | Travel insurance, Foreign currency, Copies of docs |
| 3 | Clothing | 20–39 | Jacket, Shirts, Underwear, Socks, Shoes |
| 4 | Toiletries | 40–54 | Toothbrush, Deodorant, Shampoo |
| 5 | Health & Safety | 55–64 | Pain relievers, Bandages, First aid kit |
| 6 | Electronics | 65–74 | Phone charger, Portable charger, Laptop |
| 7 | Activity Specific | 75–99 | Hiking boots, Dressy outfit, Daypack |

The frontend renders items in `sort_order` order — no client-side sorting required.

---

## Deduplication

After all rule groups run, items are deduped by name (case-insensitive). If two rules add the same item (e.g., `cold_weather` trip type and weather rules both add "Heavy coat"), keep one instance. The higher `sort_order` (category precedence) wins.

---

## Extensibility

The rule engine is a standalone module (`rules/engine.go`). To add a new rule:
1. Add the trigger condition
2. Add the items it produces
3. Assign category and sort_order

No other files change. The Claude API "Smart Mode" upgrade path works by replacing or augmenting the rule engine output — the same item structure is returned either way.
