# Page Architecture — PackSmart

## Overview

PackSmart has 2 real pages plus a 404 handler. No user accounts for MVP — trips are identified by UUID and accessible to anyone with the link.

---

## Page 1: Home

**URL:** `/`
**One job:** Collect trip details and send them to the API to generate a packing list.

**How users arrive:**
- Direct visit (first time or returning to plan a new trip)
- Clicking "← New Trip" from the packing list page
- Browser back button from `/packing-list/{uuid}`
- Any unrecognized route (redirect here via 404)

**What's on the page:**
- App name and tagline
- One-sentence explainer
- Trip input form:
  - Destination (text input)
  - Departure and return dates (date pickers, side by side)
  - Trip type (visual card grid — single select)
  - Travel companions (visual card grid — single select)
  - Activities (tag-style multi-select)
- "Generate My Packing List" submit button
- Inline validation errors when required fields are missing

**On submit:**
- Validates all required fields (show inline errors if missing)
- POSTs to `/api/trips`
- Shows loading state: "Building your list…" (expect 2–4 seconds)
- On API success: navigates to `/packing-list/{uuid}`
- On API error: shows error banner, stays on `/`, user can retry

---

## Page 2: Packing List

**URL:** `/packing-list/{uuid}`
**One job:** Show the personalized packing checklist and let the user track what they've packed.

**How users arrive:**
- Form submission on `/` (app navigates here after API responds)
- Shared link (anyone with the URL)
- Bookmark or direct return visit

**What's on the page:**
- Trip summary header: destination · duration · trip type
- "← New Trip" button and "🖨 Print" button
- Progress bar: "X of Y items packed"
- Weather forecast card: temperature range, rain/snow day count, daily strip with icons
- **⚠️ Essential Items section** (visually distinct, always at top): passport, prescriptions, power adapter — items that ruin a trip if forgotten
- Categorized checklist sections below essentials:
  - Documents & Money
  - Clothing
  - Toiletries
  - Health & Safety
  - Electronics
  - Activity Specific
- Each item: checkbox, item name, reason line ("Paris will be 52–63°F with 2 rainy days")
- Celebration message when all items are checked

**States to handle:**
- **Loading:** fetching trip data via `GET /api/trips/{uuid}`
- **Not found:** UUID doesn't exist in the database → friendly message + link to `/`
- **API error:** show error message, don't lose context
- **Fully packed:** celebration message at 100%

**Checkbox behavior:**
- Checking/unchecking an item sends `PATCH /api/trips/{uuid}/items/{itemId}`
- State persists in MySQL — works across devices, works for anyone with the link

---

## Page 3: 404 / Not Found

**URL:** Any unrecognized route
**One job:** Handle wrong URLs gracefully without confusing the user.

**How users arrive:**
- Mistyped URL
- Stale bookmark
- Bad or expired UUID in `/packing-list/{uuid}`

**What's on the page:**
- Simple "Trip not found" or "Page not found" message
- Link back to `/` to start a new trip

---

## Trip ID Format

All trip IDs are **UUIDs** (e.g., `a3f8c2d1-4b5e-4c3d-8f9a-1b2c3d4e5f6a`).

**Why UUID over auto-increment:**
- Without user accounts, sequential IDs (1, 2, 3…) would let anyone enumerate all trips by changing the number
- UUID makes trips private-by-default — only people with the link can access them
- Follows the same pattern as Google Docs, Dropbox, and other shareable-link systems
- Compatible with future user account features
