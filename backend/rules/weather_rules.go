package rules

import (
	"github.com/pickkinsley/project2/backend/models"
)

// weatherItems adds items based on actual forecast data.
// Called only when ctx.Weather != nil.
func weatherItems(ctx TripContext) []models.PackingItem {
	w := ctx.Weather
	var items []models.PackingItem

	// ── Temperature-based clothing ──────────────────────────────────────────

	if w.TempMinF < 32 {
		// Freezing conditions
		items = append(items,
			item("Heavy winter coat", "Clothing",
				"Temperatures will drop below freezing — a real coat is essential", true, 10),
			item("Thermal underwear (top + bottom)", "Clothing",
				"Base layer for sub-freezing temperatures", true, 11),
			item("Insulated gloves", "Clothing",
				"Your hands will thank you below 32°F", true, 12),
			item("Warm hat / beanie", "Clothing",
				"You lose most body heat through your head in the cold", true, 13),
			item("Scarf", "Clothing",
				"Wind chill makes cold temps feel even colder", false, 14),
			item("Wool or thermal socks", "Clothing",
				"Cotton socks lose insulation when wet — wool stays warm", true, 15),
			item("Insulated, waterproof boots", "Clothing",
				"Cold and wet feet can ruin a full day outdoors", true, 16),
		)
	} else if w.TempMinF < 45 {
		// Cold but above freezing
		items = append(items,
			item("Winter coat or heavy jacket", "Clothing",
				"Lows in the 30s–40s require a serious outer layer", true, 10),
			item("Warm sweater or fleece", "Clothing",
				"Mid-layer for cold days and chilly evenings", false, 11),
			item("Light gloves", "Clothing",
				"Handy for cold mornings and evenings", false, 12),
			item("Warm hat", "Clothing",
				"Evenings will be cold even if days are mild", false, 13),
			item("Scarf", "Clothing",
				"Versatile — adds warmth without bulk", false, 14),
		)
	} else if w.TempMinF < 60 {
		// Mild — layers recommended
		items = append(items,
			item("Light jacket or fleece", "Clothing",
				"Evenings can be cool — a layer you can tie around your waist works great", false, 10),
			item("Long-sleeve shirts (2)", "Clothing",
				"For cooler mornings and evenings", false, 11),
		)
	} else {
		// Warm / hot
		items = append(items,
			item("Lightweight, breathable tops (extra)", "Clothing",
				"Hot weather means you'll go through clothes faster", false, 10),
			item("Shorts (2 pairs)", "Clothing",
				"Comfortable for warm-weather sightseeing", false, 11),
		)
	}

	// ── Rain gear ───────────────────────────────────────────────────────────

	if w.RainDays >= 1 {
		items = append(items,
			item("Compact travel umbrella", "Essential Items",
				"Rain is forecast — a pocket umbrella weighs nothing and changes everything", true, 4),
		)
	}
	if w.RainDays >= 2 {
		items = append(items,
			item("Waterproof jacket / rain shell", "Clothing",
				"Multiple rainy days expected — an umbrella won't cut it for outdoor activities", true, 17),
			item("Waterproof shoes or water-resistant sneakers", "Clothing",
				"Wet feet in the morning will ruin your whole day", false, 18),
		)
	}

	// ── High precipitation probability on specific days ──────────────────────

	highPrecipDays := countHighPrecipDays(w.DailyForecast, 60)
	if highPrecipDays >= 2 {
		items = append(items,
			item("Dry bag or waterproof pouch", "Essential Items",
				"Protect your phone, passport, and wallet on heavy rain days", false, 5),
		)
	}

	// ── Wind protection ─────────────────────────────────────────────────────

	if avgWindSpeed(w.DailyForecast) >= 20 {
		items = append(items,
			item("Windproof jacket", "Clothing",
				"Average wind speeds are high — a windproof layer makes a huge difference", false, 19),
		)
	}

	// ── Snow days ───────────────────────────────────────────────────────────

	if w.SnowDays >= 1 {
		items = append(items,
			item("Waterproof snow boots", "Clothing",
				"Snow is forecast — waterproof boots with grip are essential", true, 20),
			item("Boot traction cleats (microspikes)", "Activity Specific",
				"Icy sidewalks are dangerous — clip-on traction gives you confidence on slick surfaces", false, 0),
			item("Hand warmers (pocket-size)", "Health & Safety",
				"Disposable warmers are cheap and invaluable in snow", false, 5),
		)
	}

	return items
}

// countHighPrecipDays returns the number of forecast days with precipitation
// probability above the given threshold.
func countHighPrecipDays(forecast []models.DailyForecast, threshold int) int {
	count := 0
	for _, day := range forecast {
		if day.PrecipProbability >= threshold {
			count++
		}
	}
	return count
}

// avgWindSpeed returns the average daily max wind speed across all forecast days.
func avgWindSpeed(forecast []models.DailyForecast) float64 {
	if len(forecast) == 0 {
		return 0
	}
	total := 0.0
	for _, day := range forecast {
		total += day.WindSpeedMph
	}
	return total / float64(len(forecast))
}
