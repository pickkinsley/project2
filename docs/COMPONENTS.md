# Component Plan — PackSmart

## Component Tree

```
<Layout>
  ├── / (Home)
  │   └── <TripForm>
  │
  ├── /packing-list/{uuid}
  │   ├── <WeatherCard>
  │   ├── <EssentialItems>
  │   │   └── <PackingItem> (repeated)
  │   └── <PackingCategory> (repeated)
  │       └── <PackingItem> (repeated)
  │
  └── * (404)
      └── <ErrorMessage>

Shared anywhere:
  <LoadingState>
  <ErrorMessage>
```

---

## Shared Components (used on 2+ pages)

### `<Layout>`
Wraps every page. Contains the app name/header and a consistent outer container. One component means one place to update the header across the whole app.

### `<LoadingState>`
Used in two places: during form submission (2–4s wait for API) and during packing list page load (GET from API). Same visual treatment, different message text passed as a prop.

### `<ErrorMessage>`
Handles form validation errors, API errors, and the 404 "trip not found" state. One component with a message prop covers all cases.

---

## Page-Specific Components

### `<TripForm>` — Home page only
The full trip input form. Complex enough (5 fields, validation, submit + loading handling, API call) to deserve its own file even though it only lives on `/`.

**Contains:** destination input, date pickers, trip type card grid, companions card grid, activities tag multi-select, submit button.

### `<WeatherCard>` — Packing list page only
Displays the forecast summary: temperature range, rain/snow day count, and a daily strip with icons and temps. Self-contained with its own display logic — isolated from the packing list.

### `<EssentialItems>` — Packing list page only
The ⚠️ warning section at the top of the packing list. Visually distinct from regular categories (different color/border). Always shown, never collapsed. Renders a list of `<PackingItem>` components.

**Separate from `<PackingCategory>` because:** different visual treatment, different rule logic (these items are always included regardless of trip type).

### `<PackingCategory>` — Packing list page only
One collapsible section of the checklist (Documents & Money, Clothing, Toiletries, etc.). The packing list page renders several of these in a loop. Shows a count badge (e.g., "8 of 14 packed"). Renders a list of `<PackingItem>` components.

### `<PackingItem>` — Used inside `<EssentialItems>` and `<PackingCategory>`
One checkbox row. Contains the checkbox, item name, and reason line ("Paris will be 52–63°F with 2 rainy days"). Checking/unchecking fires a PATCH to `/api/trips/{uuid}/items/{itemId}`. Isolating this component keeps the API call logic contained.

---

## What We're Not Planning Yet

- `<Header>` — for a 2-page app with no navigation links, the header is just the app name inside `<Layout>`. Not worth a separate component until there's nav to put in it.
- `<ProgressBar>` — likely just a few lines inside the packing list page; extract only if it becomes complex.
- Form sub-components (`<DatePicker>`, `<TripTypeCard>`, etc.) — identify during implementation, not now.

---

## Total: 9 Components

| Component | Type | Pages |
|---|---|---|
| `<Layout>` | Shared | All |
| `<LoadingState>` | Shared | `/`, `/packing-list/{uuid}` |
| `<ErrorMessage>` | Shared | All |
| `<TripForm>` | Page-specific | `/` |
| `<WeatherCard>` | Page-specific | `/packing-list/{uuid}` |
| `<EssentialItems>` | Page-specific | `/packing-list/{uuid}` |
| `<PackingCategory>` | Page-specific | `/packing-list/{uuid}` |
| `<PackingItem>` | Sub-component | `/packing-list/{uuid}` |
| 404 content | Inline | `*` |
