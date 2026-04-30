package rules

import (
	"github.com/pickkinsley/project2/backend/models"
)

// tripTypeItems adds items specific to the trip type.
func tripTypeItems(ctx TripContext) []models.PackingItem {
	switch ctx.TripType {
	case "international":
		return internationalItems(ctx)
	case "beach":
		return beachItems(ctx)
	case "cold_weather":
		return coldWeatherItems(ctx)
	case "business":
		return businessItems(ctx)
	case "domestic":
		return domesticItems(ctx)
	default:
		return nil
	}
}

func internationalItems(ctx TripContext) []models.PackingItem {
	items := []models.PackingItem{
		item("Passport", "Documents & Money",
			"Required for all international travel — check expiry (must be valid 6 months beyond return)", true, 0),
		item("Visa / entry documents", "Documents & Money",
			"Check requirements for your destination — apply weeks in advance if needed", true, 1),
		item("Travel insurance documents", "Documents & Money",
			"Print a copy and save a photo — essential if something goes wrong abroad", true, 2),
		item("Foreign currency / local cash", "Documents & Money",
			"Have some local cash for taxis, tips, and places that don't take cards", false, 3),
		item("International credit card (no foreign fees)", "Documents & Money",
			"Avoid 3% foreign transaction fees — Charles Schwab, Chase Sapphire, or similar", false, 4),
		item("Universal power adapter", "Electronics",
			"Outlet shapes and voltages vary by country — one universal adapter covers everywhere", true, 4),
		item("Copies of important documents", "Documents & Money",
			"Photo copies of passport, insurance, and visa stored separately from originals", true, 5),
		item("Translation app (downloaded offline)", "Electronics",
			"Download offline language packs before you leave — no data needed", false, 5),
		item("Locks for luggage", "Essential Items",
			"TSA-approved locks for checked bags, combination padlock for hostel lockers", false, 6),
	}

	// For longer international trips, add data/communication items
	if ctx.DurationDays >= 7 {
		items = append(items,
			item("International SIM card or travel eSIM", "Electronics",
				"Avoid roaming charges — buy a local SIM or activate an eSIM before departure", false, 6),
		)
	}

	return items
}

func beachItems(_ TripContext) []models.PackingItem {
	return []models.PackingItem{
		item("Swimsuit (2)", "Clothing",
			"Two suits so one can dry while you wear the other", true, 7),
		item("Cover-up / beach dress", "Clothing",
			"For walking between beach and restaurant", false, 8),
		item("Flip flops / sandals", "Clothing",
			"Easy on/off for beach and pool areas", false, 9),
		item("Beach bag", "Essential Items",
			"Waterproof or sand-resistant bag for towel, sunscreen, and snacks", false, 7),
		item("Beach towel", "Essential Items",
			"Hotels may charge extra for pool/beach towels — bring your own", false, 8),
		item("Waterproof phone case", "Electronics",
			"Protect your phone from sand and water at the beach", false, 6),
		item("After-sun lotion / aloe vera", "Health & Safety",
			"Treat sunburn quickly — it turns a bad day into a recovery day", false, 5),
		item("Snorkel gear (optional)", "Activity Specific",
			"Cheaper to bring your own than rent — great for reef exploring", false, 1),
		item("Reef-safe sunscreen SPF 50", "Health & Safety",
			"High SPF for beach days — many destinations require reef-safe formulas", true, 3),
	}
}

func coldWeatherItems(_ TripContext) []models.PackingItem {
	return []models.PackingItem{
		item("Heavy winter coat", "Clothing",
			"A real coat — not just a hoodie", true, 10),
		item("Thermal base layer (top + bottom)", "Clothing",
			"Wear under everything — traps heat close to your body", true, 11),
		item("Fleece or wool mid-layer", "Clothing",
			"The layer between your base and outer shell", false, 12),
		item("Insulated gloves", "Clothing",
			"Touchscreen-compatible so you can still use your phone", true, 13),
		item("Warm hat / beanie", "Clothing",
			"Covers your ears — a key heat loss point", true, 14),
		item("Neck gaiter or scarf", "Clothing",
			"Versatile — can pull up over your face in very cold wind", false, 15),
		item("Wool or thermal socks (extra pairs)", "Clothing",
			"Cold feet ruin cold-weather trips — wool stays warm even when damp", true, 16),
		item("Insulated, waterproof boots", "Clothing",
			"Non-negotiable for cold weather — your feet will be outside all day", true, 17),
		item("Hand warmers (multi-pack)", "Health & Safety",
			"Disposable warmers for gloves and boots on the coldest days", false, 5),
		item("Lip balm (heavy-duty, SPF)", "Toiletries",
			"Cold wind chaps lips fast — keep one in every pocket", false, 9),
		item("Moisturizer / thick hand cream", "Toiletries",
			"Cold, dry air cracks skin — apply morning and night", false, 10),
	}
}

func businessItems(_ TripContext) []models.PackingItem {
	return []models.PackingItem{
		item("Business attire (2 outfits)", "Clothing",
			"Professional clothing appropriate for meetings", true, 7),
		item("Dress shoes", "Clothing",
			"Polished shoes complete a professional look", false, 8),
		item("Belt", "Clothing",
			"Often forgotten — needed with dress pants", false, 9),
		item("Laptop + charger", "Electronics",
			"Bring your laptop and its charger — hotel rooms rarely have spares", true, 4),
		item("Laptop bag or briefcase", "Essential Items",
			"Protects your equipment and looks professional in meetings", false, 7),
		item("Business cards", "Documents & Money",
			"Still expected at many professional events and conferences", false, 6),
		item("Portable charger / battery pack", "Electronics",
			"Long conference days drain your phone — always have backup power", true, 5),
		item("Wrinkle-release spray", "Toiletries",
			"Packed clothes wrinkle — spray + hang in a steamy bathroom fixes most issues", false, 11),
		item("Notebook + pen", "Essential Items",
			"For meeting notes when devices aren't appropriate or allowed", false, 9),
		item("Formal tie or accessories (optional)", "Clothing",
			"Bring one formal option in case the dress code is stricter than expected", false, 10),
	}
}

func domesticItems(_ TripContext) []models.PackingItem {
	return []models.PackingItem{
		item("Government-issued photo ID", "Documents & Money",
			"Required for domestic flights and many venues — driver's license or state ID", true, 0),
		item("Credit / debit card", "Documents & Money",
			"Primary payment method — bring a backup card in case one is declined", false, 1),
	}
}
