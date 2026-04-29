# PackSmart Code Cleanup — Lesson 2
**Date:** April 29, 2026
**Scope:** Project structure, shared utilities, naming consistency, and error handling

---

## Project Structure Improvements

### Files Moved

Two documentation files were living in `frontend/` where they did not belong. Documentation describes the project as a whole, not the frontend specifically.

| From | To |
|------|----|
| `frontend/CODE_REVIEW.md` | `docs/CODE_REVIEW.md` |
| `frontend/POLISH_IMPROVEMENTS.md` | `docs/POLISH_IMPROVEMENTS.md` |

All project documentation now lives in `docs/`, keeping the frontend directory focused on source code.

### Files Deleted

Five files were confirmed unused and removed entirely.

| File | Reason |
|------|--------|
| `frontend/src/App.css` | 185 lines of Vite boilerplate — never imported anywhere in the app |
| `frontend/src/lib/` | Empty folder with no contents |
| `frontend/src/assets/react.svg` | Default Vite asset — not referenced in any component |
| `frontend/src/assets/vite.svg` | Default Vite asset — not referenced in any component |
| `frontend/src/assets/hero.png` | Unused image asset |

No functionality was lost. These files were left over from the initial Vite scaffold and were never integrated into the app.

### Git Configuration

`backend/server_linux` was added to `.gitignore`. The compiled Go binary for Linux (`server_linux`) was already excluded for macOS (`server`), but the Linux variant used in production deployment was missing from the ignore list.

```
# Go build output
backend/server
backend/server_linux   ← added
backend/*.exe
```

Both compiled binaries are now excluded from version control.

---

## Code Quality Improvements

### Extracted Shared Utility

`formatDate` was defined identically in two separate page files. Duplicated logic means two places to update if the format ever changes, and two places where bugs can independently appear.

**Created:** `frontend/src/utils/formatDate.js`

```js
export function formatDate(dateStr) {
  const [y, m, d] = dateStr.split('-')
  return new Date(y, m - 1, d).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
}
```

**Updated imports:**

| File | Change |
|------|--------|
| `PackingListPage.jsx` | Removed local definition, added import |
| `MyTripsPage.jsx` | Removed local definition, added import |
| `EditTripPage.jsx` | No change needed — did not actually use `formatDate` |

There is now a single source of truth for date formatting across the entire frontend.

### Naming Improvements

Four variables were renamed to improve clarity and consistency.

**`Header.jsx` — `linkClass` → `navLinkClass`**

The original name `linkClass` was generic. `navLinkClass` makes it immediately clear this function computes CSS classes for navigation links specifically, distinguishing it from any other link styling that might exist in the file.

**`HomePage.jsx` — `mutation` → `createMutation`**

Generic `mutation` names are ambiguous when a component grows to have multiple mutations. Naming it `createMutation` reflects what it does (creates a trip) and matches the naming pattern used in other pages.

**`EditTripPage.jsx` — `mutation` → `updateMutation`**

Same rationale. `updateMutation` is self-documenting — it updates an existing trip — and is consistent with `createMutation` on the home page.

**`PackingListPage.jsx` — `previous` → `snapshot`**

In TanStack Query's optimistic update pattern, the variable holding the pre-mutation cache state is conventionally named `snapshot`. The original name `previous` was functional but non-standard. Renamed in both `toggleMutation` and `addMutation`.

```js
// Before
const previous = queryClient.getQueryData(['trip', tripId])
return { previous }
// ...
queryClient.setQueryData(['trip', tripId], context.previous)

// After
const snapshot = queryClient.getQueryData(['trip', tripId])
return { snapshot }
// ...
queryClient.setQueryData(['trip', tripId], context.snapshot)
```

### Consistency Fixes

**Replaced `alert()` in delete error handling with inline error state**

The delete button was the only place in the app that used a native browser `alert()` dialog for error handling. Every other mutation used an inline error message rendered in the component. This inconsistency was fixed.

**Before:**
```js
onError: (err) => alert(err?.message ?? 'Failed to delete trip. Please try again.')
```

**After:**
```js
onError: (err) => {
  setDeleteError(err?.message ?? 'Failed to delete trip. Please try again.')
  setTimeout(() => setDeleteError(null), 3000)
},
```

A `deleteError` state variable was added alongside the existing `toggleError`, and the error message is displayed inline near the delete button — the same pattern used throughout the rest of the app. The error auto-clears after 3 seconds, matching `toggleError` behavior.

---

## Bugs Introduced and Fixed

### `mutationFn` key corrupted by `replace_all` rename

**Files affected:** `HomePage.jsx`, `EditTripPage.jsx`

When renaming `mutation` → `createMutation` and `mutation` → `updateMutation` using a global find-and-replace, the rename also hit the `mutationFn` option key inside the `useMutation` config object:

```js
// Intended result
const createMutation = useMutation({
  mutationFn: createTrip,   // ← key name, should not be renamed
  ...
})

// Actual result after replace_all
const createMutation = useMutation({
  createMutationFn: createTrip,   // ← corrupted — TanStack Query doesn't recognize this
  ...
})
```

TanStack Query silently ignores unrecognized keys and falls back to no `mutationFn`, producing the runtime error: **"No mutationFn found"** when the form was submitted.

**Fix:** Restored `mutationFn` (and `updateMutationFn` → `mutationFn`) in both files.

**Lesson:** `replace_all` renames are blunt instruments. The string being renamed (`mutation`) appeared both as a variable name *and* as a substring of an unrelated config key (`mutationFn`). A more targeted replace — or a manual review pass after the rename — would have caught this immediately. When a rename touches a high-traffic keyword, verify that no option keys, object properties, or string literals were unintentionally affected.

---

## Impact Summary

| Metric | Result |
|--------|--------|
| Lines of code removed | ~200 (185 from App.css, ~15 from duplicate `formatDate` definitions) |
| Files deleted | 5 |
| Files moved to correct location | 2 |
| Duplicate function definitions eliminated | 1 (`formatDate`, 2 copies → 1 shared utility) |
| Variables renamed for clarity | 4 |
| Native browser dialogs removed | 1 (`alert()` → inline error state) |

The shared `formatDate` utility is the most maintainable change: date formatting requirements can evolve (locale, format options, timezone handling), and a single file to update is strictly better than hunting down all call sites.

The naming improvements have no runtime effect but reduce the mental load of reading the code. A developer encountering `createMutation` and `updateMutation` side by side immediately understands what each does and which page they belong to.

The `alert()` removal means error handling is now uniform across all mutations in the app. Uniform patterns are easier to reason about, easier to test, and easier to extend.

---

## What Was Not Changed

Several refactors were considered and deliberately declined. Abstraction has a cost — more files, more imports, more indirection — and that cost only pays off when the abstraction is actually reused or when the original code is genuinely hard to understand.

**`WeatherCard` was not extracted into its own file.**
The weather display block is only rendered in one place (`PackingListPage`). Extracting it would create a new file that exists solely to be imported once. The complexity added by the extraction outweighs any benefit.

**`CategoryCard` / `PackingItem` sub-components were not extracted.**
These are small, tightly coupled to the packing list context, and used only within `PackingListPage`. Keeping them co-located makes it easier to see how they fit together.

**`useDeleteTrip` hook was not created.**
The delete mutation is eight lines. A custom hook wrapping eight lines adds a file, an import, and an abstraction layer that provides no clarity benefit at this size.

**`AddItemForm` was not extracted into its own file.**
This was deferred — it is a reasonable candidate for extraction if the form grows, but in its current state the co-location with `PackingListPage` is clearer.

The principle: extract when there is actual duplication or when the component is genuinely complex enough that isolation helps. Do not extract for its own sake.

---

## Lessons Learned

**Remove unused code immediately.**
Unused files (boilerplate, dead assets, empty folders) are not harmless. They create noise that makes it harder to understand what the project actually contains. Deleting them early — before they get committed alongside real work — keeps the repository clean and intentional from the start.

**Extract duplicated utilities, not single-use code.**
`formatDate` appeared twice with identical implementations. That is the right trigger for extraction: real duplication, a function with a clear boundary, and a utility general enough to live outside any single component. Single-use code belongs where it is used.

**Consistent naming reduces cognitive load.**
`mutation`, `createMutation`, and `updateMutation` all refer to TanStack Query mutations. The generic name forces the reader to look at the function body to understand what it does. The specific names are self-documenting. Small naming improvements compound across an entire codebase — every function that reads clearly is one less thing a future reader has to decode.

**Consistent patterns are a form of documentation.**
When every mutation in the app handles errors the same way (inline state, auto-clear timeout), a developer reading any one mutation immediately knows how error handling works everywhere. When one mutation uses `alert()`, that consistency is broken and the outlier demands extra attention. Uniformity makes a codebase easier to maintain because there are fewer special cases to remember.

**`replace_all` renames require a verification pass.**
Global string replacement is fast but indiscriminate. Renaming `mutation` hit not just the variable but also the `mutationFn` config key, breaking both pages silently at runtime. After any broad rename, scan the diff for unintended collateral changes — especially when the renamed string is a common substring.
