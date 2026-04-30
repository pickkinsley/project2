package rules

import (
	"fmt"

	"github.com/pickkinsley/project2/backend/models"
)

// baseItems returns items every traveler needs regardless of trip type,
// with clothing quantities scaled to trip duration.
func baseItems(ctx TripContext) []models.PackingItem {
	d := ctx.DurationDays

	// Clothing quantities — standard travel math
	underwear := clamp(d+1, 3, 10)
	socks := clamp(d+1, 3, 10)
	shirts := clamp((d+1)/2+1, 2, 7)
	pants := clamp(d/3+1, 1, 4)

	items := []models.PackingItem{
		// ── Essentials ──────────────────────────────────────────────────────
		item("Prescriptions", "Essential Items",
			"Never travel without your medications — pack more than you think you need", true, 0),
		item("Health insurance card", "Essential Items",
			"Required for medical care abroad or out of network", true, 1),
		item("Emergency contact list", "Essential Items",
			"Keep a written copy in case your phone dies", true, 2),

		// ── Clothing ────────────────────────────────────────────────────────
		item(fmt.Sprintf("Underwear (%d pairs)", underwear), "Clothing",
			fmt.Sprintf("1 per day + 1 spare for a %d-day trip", d), false, 0),
		item(fmt.Sprintf("Socks (%d pairs)", socks), "Clothing",
			fmt.Sprintf("1 per day + 1 spare for a %d-day trip", d), false, 1),
		item(fmt.Sprintf("T-shirts / tops (%d)", shirts), "Clothing",
			"Plan to re-wear — you can always do laundry", false, 2),
		item(fmt.Sprintf("Pants / bottoms (%d)", pants), "Clothing",
			"Jeans or versatile trousers that work for multiple occasions", false, 3),
		item("Sleepwear", "Clothing",
			"Comfortable clothes for sleeping", false, 4),
		item("Comfortable walking shoes", "Clothing",
			"Broken-in shoes you can walk miles in — blisters ruin trips", true, 5),
		item("Casual outfit", "Clothing",
			"Something relaxed for downtime and light activities", false, 6),

		// ── Toiletries ──────────────────────────────────────────────────────
		item("Toothbrush + toothpaste", "Toiletries",
			"Non-negotiable daily essential", false, 0),
		item("Deodorant", "Toiletries",
			"Travel-size if flying — under 3.4 oz for carry-on", false, 1),
		item("Shampoo + conditioner", "Toiletries",
			"Travel-size or check if your accommodation provides them", false, 2),
		item("Body wash / soap", "Toiletries",
			"Travel-size bar soap takes up less space than liquid", false, 3),
		item("Face wash + moisturizer", "Toiletries",
			"Skin gets dry when traveling — keep your routine", false, 4),
		item("Lip balm", "Toiletries",
			"Planes and new climates are notoriously drying", false, 5),
		item("Razor + shaving cream", "Toiletries",
			"Travel-size razor or disposable", false, 6),
		item("Hair brush / comb", "Toiletries",
			"Easy to forget, impossible to live without", false, 7),
		item("Feminine hygiene products", "Toiletries",
			"Bring enough for the trip — availability varies by destination", false, 8),

		// ── Health & Safety ─────────────────────────────────────────────────
		item("Pain reliever (ibuprofen / acetaminophen)", "Health & Safety",
			"Headaches and minor pain are common when traveling", false, 0),
		item("Antacids / stomach medicine", "Health & Safety",
			"New food can upset your stomach — be prepared", false, 1),
		item("Band-aids + antiseptic wipes", "Health & Safety",
			"Small cuts and blisters happen on every trip", false, 2),
		item("Hand sanitizer", "Health & Safety",
			"Airports, planes, and public transit are germ-heavy", false, 3),
		item("Sunscreen SPF 30+", "Health & Safety",
			"Sun exposure adds up quickly when you're outside all day", true, 4),

		// ── Electronics ─────────────────────────────────────────────────────
		item("Phone + charger", "Electronics",
			"Your most important travel tool — charge every night", true, 0),
		item("Portable battery / power bank", "Electronics",
			"Essential for long days of navigation and photos", false, 1),
		item("Headphones / earbuds", "Electronics",
			"For flights, transit, and blocking out hotel noise", false, 2),
		item("Camera (optional)", "Electronics",
			"Phone cameras are great — bring a dedicated camera if you care about photo quality", false, 3),
	}

	// Add laundry supplies for longer trips
	if d >= 7 {
		items = append(items,
			item("Travel laundry detergent", "Toiletries",
				"For week+ trips — doing laundry mid-trip saves suitcase space", false, 9),
		)
	}

	// Add a packing cube or compression bag for longer trips
	if d >= 5 {
		items = append(items,
			item("Packing cubes", "Essential Items",
				"Keeps clothes organized and makes repacking fast", false, 3),
		)
	}

	return items
}
