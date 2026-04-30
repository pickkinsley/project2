// Package rules generates a smart packing list from trip and weather data.
// Call GeneratePackingList — it runs all rule groups in order and returns
// a deduplicated, sorted slice of PackingItem ready to insert into the DB.
package rules

import (
	"strings"

	"github.com/pickkinsley/project2/backend/models"
)

// TripContext is the input to the engine.
type TripContext struct {
	TripType     string
	Companions   string
	Activities   []string
	DurationDays int
	Weather      *models.WeatherResponse // nil when forecast unavailable
}

// sortOrderBase maps category names to their starting sort_order values.
// Items within a category are assigned base + offset so they stay grouped.
var sortOrderBase = map[string]int{
	"Essential Items":  1,
	"Documents & Money": 10,
	"Clothing":         20,
	"Toiletries":       40,
	"Health & Safety":  55,
	"Electronics":      65,
	"Activity Specific": 75,
}

// item is a builder helper — shorthand for constructing a PackingItem.
func item(name, category, reason string, essential bool, offset int) models.PackingItem {
	base, ok := sortOrderBase[category]
	if !ok {
		base = 90
	}
	return models.PackingItem{
		Name:        name,
		Category:    category,
		Reason:      reason,
		IsEssential: essential,
		SortOrder:   base + offset,
	}
}

// GeneratePackingList runs all rule groups in order and returns a
// deduplicated, sort_order-assigned slice of items.
//
// Execution order:
//  1. Base items (toiletries + clothing quantities)
//  2. Trip-type rules
//  3. Weather rules (skipped when ctx.Weather == nil)
//  4. Activity rules
//  5. Companion rules
//  6. Dedupe by name (case-insensitive, first occurrence wins)
func GeneratePackingList(ctx TripContext) []models.PackingItem {
	var all []models.PackingItem

	all = append(all, baseItems(ctx)...)
	all = append(all, tripTypeItems(ctx)...)
	if ctx.Weather != nil {
		all = append(all, weatherItems(ctx)...)
	}
	all = append(all, activityItems(ctx)...)
	all = append(all, companionItems(ctx)...)

	return dedupe(all)
}

// dedupe removes duplicate items by name (case-insensitive).
// The first occurrence of each name is kept; later ones are dropped.
func dedupe(items []models.PackingItem) []models.PackingItem {
	seen := make(map[string]bool, len(items))
	out := make([]models.PackingItem, 0, len(items))
	for _, it := range items {
		key := strings.ToLower(it.Name)
		if !seen[key] {
			seen[key] = true
			out = append(out, it)
		}
	}
	return out
}

// hasActivity returns true if the given activity is in the context.
func hasActivity(ctx TripContext, activity string) bool {
	for _, a := range ctx.Activities {
		if a == activity {
			return true
		}
	}
	return false
}

// clamp returns n clamped to [min, max].
func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
