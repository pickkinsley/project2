# Navigation Flow — PackSmart

## Flow Diagram

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
│  🌤 Forecast: 52–63°F · 2 rainy days                            │
│  [Mon ⛅ 63°] [Tue 🌧 55°] [Wed ☀️ 61°] ...                   │
│                                                                  │
│  ⚠️  ESSENTIAL ITEMS                                             │
│     ☐ Passport   ☐ Prescriptions   ☐ Power Adapter              │
│                                                                  │
│  📋 Documents & Money  (2/5)                                     │
│  👕 Clothing           (8/14)                                    │
│  🧴 Toiletries         (6/9)  ...                               │
│                                                                  │
│         [ ☑ checkbox ] → PATCH /api/trips/{uuid}/items/{id}     │
└──────────────────────────────────────────────────────────────────┘
         │              │                │
    [← New Trip]    Shared link      All items checked
    or back button  to someone else       │
         │              │                ▼
         │              │        "🎉 You're all packed!"
         │              │        message stays on page
         │              │
         ▼              ▼
┌──────────┐   ┌─────────────────────────────┐
│  HOME /  │   │  SAME PACKING LIST PAGE     │
│  (blank  │   │  Loads same data from MySQL │
│   form)  │   │  via GET /api/trips/{uuid}  │
└──────────┘   └─────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                  UUID               UUID
                  found             not found
                    │                   │
                    ▼                   ▼
              Show trip           ┌──────────────┐
              and list            │  404 PAGE  * │
                                  │              │
                                  │ "Trip not    │
                                  │  found"      │
                                  │              │
                                  │ [← Plan a    │
                                  │  new trip]   │
                                  └──────────────┘
                                         │
                                         ▼
                                      HOME /
```

---

## Flow Summary Table

| Trigger | Where they are | Where they go |
|---|---|---|
| First visit | — | `/` |
| Form submit (valid) | `/` | Loading → `/packing-list/{uuid}` |
| Form submit (invalid) | `/` | Stay on `/`, show inline errors |
| API error on submit | `/` | Stay on `/`, show error banner |
| "← New Trip" button | `/packing-list/{uuid}` | `/` |
| Browser back | `/packing-list/{uuid}` | `/` |
| Shared UUID link | — | `/packing-list/{uuid}` |
| Bad UUID in URL | — | 404 → `/` |
| Unknown route | — | 404 → `/` |
| Check a checkbox | `/packing-list/{uuid}` | Stay on page, PATCH to API |
