package rules

import (
	"github.com/pickkinsley/project2/backend/models"
)

// activityItems adds gear specific to each selected activity.
// Multiple activities stack — a hiking + beach trip gets both sets.
func activityItems(ctx TripContext) []models.PackingItem {
	var items []models.PackingItem

	for _, activity := range ctx.Activities {
		switch activity {
		case "hiking":
			items = append(items, hikingItems(ctx)...)
		case "beach":
			items = append(items, beachActivityItems()...)
		case "skiing":
			items = append(items, skiingItems()...)
		case "water_sports":
			items = append(items, waterSportsItems()...)
		case "nightlife":
			items = append(items, nightlifeItems()...)
		case "fine_dining":
			items = append(items, fineDiningItems()...)
		case "shopping":
			items = append(items, shoppingItems()...)
		case "sightseeing":
			items = append(items, sightseeingItems()...)
		case "photography":
			items = append(items, photographyItems()...)
		case "museums":
			items = append(items, museumsItems()...)
		case "adventure_sports":
			items = append(items, adventureSportsItems()...)
		case "camping":
			items = append(items, campingItems()...)
		case "road_trip":
			items = append(items, roadTripItems()...)
		case "spa":
			items = append(items, spaItems()...)
		case "wine_tasting":
			items = append(items, wineTastingItems()...)
		case "food_tours":
			items = append(items, foodToursItems()...)
		case "concerts":
			items = append(items, concertsItems()...)
		case "theme_parks":
			items = append(items, themeParksItems()...)
		case "festivals":
			items = append(items, festivalsItems()...)
		case "business":
			// business meetings via activity — covered by trip_type_rules for business trips,
			// but handle standalone activity selection too
			items = append(items, businessMeetingItems()...)
		}
	}

	return items
}

func hikingItems(ctx TripContext) []models.PackingItem {
	items := []models.PackingItem{
		item("Hiking boots (broken in)", "Activity Specific",
			"Do not hike in new boots — break them in at home first to avoid blisters", true, 0),
		item("Hiking socks (moisture-wicking, 2+ pairs)", "Activity Specific",
			"Wool or synthetic — cotton socks cause blisters on long trails", true, 1),
		item("Day backpack (20–30L)", "Activity Specific",
			"Carry water, snacks, layers, and a first aid kit on the trail", false, 2),
		item("Reusable water bottle (1L+)", "Activity Specific",
			"Hydration is critical on trails — bring more water than you think you need", true, 3),
		item("Trail snacks / energy bars", "Activity Specific",
			"High-energy snacks for mid-hike fuel: nuts, bars, dried fruit", false, 4),
		item("Trekking poles (optional)", "Activity Specific",
			"Dramatically reduce knee strain on steep descents", false, 5),
		item("Blister pads / moleskin", "Health & Safety",
			"Treat hot spots before they become full blisters", false, 6),
	}

	// Add electrolytes for longer trips with more hiking days
	if ctx.DurationDays >= 4 {
		items = append(items,
			item("Electrolyte packets", "Health & Safety",
				"Replenish sodium and minerals lost through sweat on long hikes", false, 7),
		)
	}

	return items
}

func beachActivityItems() []models.PackingItem {
	// Beach items via activity (complements beach trip type items)
	return []models.PackingItem{
		item("Swimsuit (2)", "Activity Specific",
			"Two suits so one can dry while you wear the other", true, 6),
		item("Beach towel", "Activity Specific",
			"Quick-dry microfiber towel takes up a fraction of the space", false, 7),
		item("Reef-safe sunscreen SPF 50", "Health & Safety",
			"High SPF for beach days — many destinations require reef-safe formulas", true, 3),
		item("Waterproof phone case", "Electronics",
			"Protect your phone from sand and salt water", false, 7),
		item("Flip flops / sandals", "Activity Specific",
			"For walking between beach, pool, and outdoor showers", false, 8),
		item("After-sun lotion / aloe vera", "Health & Safety",
			"Treat sunburn fast — it gets worse overnight without treatment", false, 8),
	}
}

func skiingItems() []models.PackingItem {
	return []models.PackingItem{
		item("Ski goggles", "Activity Specific",
			"UV-protected goggles are essential — sun reflects intensely off snow", true, 0),
		item("Ski gloves / mittens", "Activity Specific",
			"Waterproof and insulated — mittens are warmer than gloves for very cold days", true, 1),
		item("Helmet (rent or bring)", "Activity Specific",
			"Ski resort rentals are fine — bring your own if you ski regularly", true, 2),
		item("Ski socks (2–3 pairs)", "Activity Specific",
			"Tall wool socks designed for ski boots — regular socks cause blisters and cold feet", true, 3),
		item("Neck gaiter / balaclava", "Activity Specific",
			"Protects your face and neck from wind chill and snow spray", false, 4),
		item("Hand warmers (multi-pack)", "Health & Safety",
			"Stuff into gloves on the coldest runs — lifesavers at altitude", false, 5),
		item("Lip balm SPF 30 (heavy-duty)", "Toiletries",
			"Wind, cold, and UV at altitude destroy lips — apply constantly", true, 12),
		item("High SPF sunscreen for snow", "Health & Safety",
			"UV is stronger at altitude and reflects off snow — burns happen fast", true, 3),
		item("Ski jacket (if not renting)", "Activity Specific",
			"Waterproof, insulated — check if your resort has rental packages", false, 5),
		item("Ski pants / bibs (if not renting)", "Activity Specific",
			"Waterproof and warm — bibs keep snow out better than regular pants", false, 6),
	}
}

func waterSportsItems() []models.PackingItem {
	return []models.PackingItem{
		item("Rash guard / sun shirt", "Activity Specific",
			"UV protection for long hours on the water — more effective than sunscreen alone", false, 0),
		item("Water shoes", "Activity Specific",
			"Protect feet on rocky shores, reefs, and boat decks", false, 1),
		item("Waterproof dry bag", "Activity Specific",
			"Keep your valuables dry when kayaking, paddleboarding, or on a boat", true, 2),
		item("Waterproof sunscreen SPF 50", "Health & Safety",
			"Water-resistant formula that stays on longer in and around water", true, 3),
		item("Quick-dry towel", "Activity Specific",
			"Microfiber towel dries fast and takes up minimal space", false, 3),
	}
}

func nightlifeItems() []models.PackingItem {
	return []models.PackingItem{
		item("Going-out outfit", "Clothing",
			"Something you feel great in — many clubs enforce dress codes", false, 12),
		item("Dressy shoes", "Clothing",
			"Sneakers are fine for some venues; have a smarter option just in case", false, 13),
		item("Small crossbody bag or clutch", "Essential Items",
			"Secure and hands-free for nights out — just enough for phone, ID, and card", false, 10),
		item("Portable phone charger", "Electronics",
			"Long nights out drain your battery — don't get stranded without GPS or Uber", true, 5),
		item("Earplugs", "Health & Safety",
			"Protect your hearing in loud venues — also great for sleeping", false, 9),
	}
}

func fineDiningItems() []models.PackingItem {
	return []models.PackingItem{
		item("Smart casual or formal outfit", "Clothing",
			"Fine dining restaurants often enforce dress codes — check in advance", true, 12),
		item("Dress shoes or heels", "Clothing",
			"Complete the look — sneakers won't meet most fine dining standards", false, 13),
		item("Reservations printout / confirmation", "Documents & Money",
			"Have your booking reference ready — some restaurants are strict about walk-ins", false, 7),
	}
}

func shoppingItems() []models.PackingItem {
	return []models.PackingItem{
		item("Foldable reusable bag", "Essential Items",
			"Many countries charge for plastic bags — a packable tote weighs almost nothing", false, 11),
		item("Extra luggage space / soft duffel", "Essential Items",
			"Leave room in your bag or pack a foldable duffel for shopping hauls", false, 12),
		item("Comfortable walking shoes (extra pair)", "Clothing",
			"Shopping means hours on your feet — have a comfortable backup pair", false, 14),
		item("Tape measure (small)", "Essential Items",
			"For furniture, art, or clothing with unfamiliar sizing — saves return shipping", false, 13),
	}
}

func sightseeingItems() []models.PackingItem {
	return []models.PackingItem{
		item("Comfortable walking shoes", "Clothing",
			"Sightseeing means 8–15k steps per day — your feet will decide the day", true, 5),
		item("Day bag / backpack", "Essential Items",
			"Carry water, snacks, a jacket, and your camera without straining your back", false, 11),
		item("Reusable water bottle", "Activity Specific",
			"Stay hydrated while exploring — refill at fountains and cafes to save money", false, 10),
		item("City map or offline maps downloaded", "Electronics",
			"Download your destination in Google Maps offline — no data needed for navigation", false, 8),
	}
}

func photographyItems() []models.PackingItem {
	return []models.PackingItem{
		item("Camera + extra batteries", "Electronics",
			"Cold weather and long days drain batteries — bring spares", true, 4),
		item("Extra memory cards", "Electronics",
			"Run out of storage on day 1 otherwise — 64GB minimum", true, 5),
		item("Camera cleaning kit", "Electronics",
			"Dust and smudges ruin shots — a blower and cloth take up no space", false, 6),
		item("Lightweight tripod or gorilla pod", "Activity Specific",
			"For low-light shots, self-portraits, and long exposures", false, 0),
		item("Lens filters (ND / polarizer)", "Activity Specific",
			"Polarizer cuts glare on water and glass; ND for long exposures in daylight", false, 1),
		item("Camera rain cover", "Activity Specific",
			"Protect your gear in rain or dusty environments", false, 2),
	}
}

func museumsItems() []models.PackingItem {
	return []models.PackingItem{
		item("Comfortable walking shoes", "Clothing",
			"Museum visits mean hours on hard floors — comfort matters more than style", true, 5),
		item("Museum membership card(s)", "Documents & Money",
			"Some memberships grant reciprocal access at affiliated museums — check before you go", false, 8),
		item("Small notebook for notes", "Essential Items",
			"Jot down artwork names, dates, and things to research later", false, 13),
	}
}

func adventureSportsItems() []models.PackingItem {
	return []models.PackingItem{
		item("Athletic clothing (moisture-wicking)", "Activity Specific",
			"Breathable, quick-dry fabric for physical activities", true, 0),
		item("Sturdy athletic shoes or approach shoes", "Activity Specific",
			"More grip and support than regular sneakers for outdoor adventures", true, 1),
		item("Helmet (activity-specific)", "Activity Specific",
			"Many adventure sports require helmets — check if rentals are available at your destination", true, 2),
		item("First aid kit (extended)", "Health & Safety",
			"Adventure sports carry injury risk — carry more than basic band-aids", true, 3),
		item("Emergency whistle", "Health & Safety",
			"Weighs nothing and can save your life if you need to signal for help", false, 4),
		item("Lightweight rain jacket", "Clothing",
			"Weather changes fast in outdoor adventure settings", false, 19),
	}
}

func campingItems() []models.PackingItem {
	return []models.PackingItem{
		item("Tent (if not provided)", "Activity Specific",
			"Confirm whether your campsite provides shelters or you need your own", true, 0),
		item("Sleeping bag (temperature-rated)", "Activity Specific",
			"Match the bag rating to nighttime temps — going too light means a miserable night", true, 1),
		item("Sleeping pad or inflatable mat", "Activity Specific",
			"Insulates you from the cold ground and cushions for better sleep", false, 2),
		item("Headlamp + extra batteries", "Activity Specific",
			"Hands-free light is essential at night — phone flashlights don't cut it", true, 3),
		item("Camp cooking kit (if cooking)", "Activity Specific",
			"Lightweight titanium pot and utensils — check if your site has facilities first", false, 4),
		item("Water filter or purification tablets", "Health & Safety",
			"For backcountry camping where tap water isn't available", false, 5),
		item("Bear canister or hang bag (if backcountry)", "Activity Specific",
			"Required in many national parks — keeps wildlife and your food safe", false, 6),
		item("Insect repellent (DEET or Picaridin)", "Health & Safety",
			"Mosquitoes and ticks are serious — apply before dusk and check for ticks daily", true, 6),
		item("Biodegradable soap", "Toiletries",
			"Leave no trace — regular soap harms ecosystems when rinsed into nature", false, 13),
	}
}

func roadTripItems() []models.PackingItem {
	return []models.PackingItem{
		item("Car phone mount", "Electronics",
			"Keep your phone accessible for navigation without holding it", true, 7),
		item("Car charger / USB adapter", "Electronics",
			"Keep devices charged on the road", true, 8),
		item("Roadside emergency kit", "Health & Safety",
			"Jumper cables, reflective triangle, flashlight — keep in the trunk", true, 5),
		item("Snacks + drinks cooler", "Activity Specific",
			"Road trip staple — saves money and time vs. stopping for every meal", false, 7),
		item("Reusable water bottles", "Activity Specific",
			"Stay hydrated on long drives — fill up at rest stops", false, 8),
		item("Paper map / road atlas (backup)", "Documents & Money",
			"GPS fails — a physical map is your backup in dead zones", false, 9),
		item("Sunglasses", "Clothing",
			"Driving into sun is dangerous — polarized lenses reduce glare", true, 15),
		item("Trash bags (small)", "Essential Items",
			"Keep the car clean — one bag per seat area", false, 14),
	}
}

func spaItems() []models.PackingItem {
	return []models.PackingItem{
		item("Spa sandals / flip flops", "Activity Specific",
			"Required footwear in most spa facilities", false, 9),
		item("Comfortable lounge wear", "Clothing",
			"Something relaxed to wear between treatments", false, 14),
		item("Hair tie / clips", "Toiletries",
			"Keep hair out of the way during facials and body treatments", false, 13),
	}
}

func wineTastingItems() []models.PackingItem {
	return []models.PackingItem{
		item("Smart casual outfit", "Clothing",
			"Wineries range from casual to formal — smart casual works everywhere", false, 12),
		item("Reusable wine bag or soft carrier", "Activity Specific",
			"For transporting bottles you purchase — protect them in your luggage", false, 0),
		item("Corkscrew / wine key", "Activity Specific",
			"For bottles purchased without screw caps — pack in checked luggage if flying", false, 1),
		item("Water bottle", "Activity Specific",
			"Stay hydrated between tastings — wine dehydrates faster than you'd think", true, 2),
		item("Crackers / bread snacks", "Activity Specific",
			"Eating between pours keeps palate fresh and alcohol absorption slower", false, 3),
	}
}

func foodToursItems() []models.PackingItem {
	return []models.PackingItem{
		item("Comfortable walking shoes", "Clothing",
			"Food tours cover a lot of ground between stops", true, 5),
		item("Light jacket or layer", "Clothing",
			"Tours often include outdoor markets in the morning when it's cooler", false, 11),
		item("Small cash (local currency)", "Documents & Money",
			"Some street food vendors and market stalls are cash only", false, 10),
		item("Antacids", "Health & Safety",
			"Eating many unfamiliar foods in one day can challenge your stomach", false, 2),
	}
}

func concertsItems() []models.PackingItem {
	return []models.PackingItem{
		item("Earplugs (musician-grade)", "Health & Safety",
			"High-fidelity earplugs protect hearing while preserving sound quality — worth it", false, 8),
		item("Portable phone charger", "Electronics",
			"Long events drain your battery — essential for photos, Uber, and emergencies", true, 5),
		item("Comfortable standing shoes", "Clothing",
			"Concerts often mean hours of standing — comfort beats style here", false, 14),
		item("Light layer / jacket", "Clothing",
			"Venues can be cold with AC or cold outside at night", false, 11),
		item("Tickets / confirmation (digital or printed)", "Documents & Money",
			"Screenshot your tickets in case you lose cell signal inside the venue", true, 8),
	}
}

func themeParksItems() []models.PackingItem {
	return []models.PackingItem{
		item("Comfortable walking shoes (broken in)", "Clothing",
			"Theme parks are 15–25k step days — worn-in shoes are essential", true, 5),
		item("Fanny pack / small belt bag", "Essential Items",
			"Keeps hands free on rides and is harder to steal than a backpack", false, 10),
		item("Portable phone charger", "Electronics",
			"Apps, photos, and digital tickets drain your battery fast", true, 5),
		item("Refillable water bottle", "Activity Specific",
			"Many parks have free water stations — save money and stay hydrated", true, 11),
		item("Sunscreen SPF 30+ (reapply)", "Health & Safety",
			"Long outdoor park days mean high UV exposure — reapply every 2 hours", true, 3),
		item("Rain poncho (compact)", "Essential Items",
			"Park ponchos are expensive — bring a $2 compact one from home", false, 12),
		item("Comfortable change of clothes", "Clothing",
			"Water rides guarantee you'll get soaked — have dry clothes in a locker", false, 15),
		item("Snacks (non-perishable)", "Activity Specific",
			"Park food is expensive — approved snacks save money and prevent hunger meltdowns", false, 12),
	}
}

func festivalsItems() []models.PackingItem {
	return []models.PackingItem{
		item("Comfortable festival outfit(s)", "Clothing",
			"Something you don't mind getting dirty, wet, or sweaty", false, 12),
		item("Sturdy boots or closed-toe shoes", "Clothing",
			"Muddy fields and crowds require real footwear — sandals are a mistake", true, 14),
		item("Portable phone charger", "Electronics",
			"Festival cell coverage is poor and your battery will die from photos and maps", true, 5),
		item("Earplugs", "Health & Safety",
			"Protect your hearing — hearing loss from festivals is permanent", false, 8),
		item("Small backpack (theft-resistant)", "Essential Items",
			"Keep essentials close — zippers face inward in crowds", false, 11),
		item("Cash (small bills)", "Documents & Money",
			"Many festival vendors are cash-only — ATM lines get long", false, 10),
		item("Reusable water bottle", "Activity Specific",
			"Most festivals have free water refill stations — staying hydrated is critical", true, 13),
		item("Insect repellent", "Health & Safety",
			"Outdoor evening events attract mosquitoes", false, 6),
		item("Sunscreen SPF 30+", "Health & Safety",
			"Full-day outdoor exposure — reapply every 2 hours", true, 4),
		item("Light rain jacket (packable)", "Clothing",
			"Festival weather is unpredictable — a compact rain layer rolls up to nothing", false, 19),
	}
}

func businessMeetingItems() []models.PackingItem {
	return []models.PackingItem{
		item("Business attire", "Clothing",
			"Professional clothing appropriate for client meetings or presentations", true, 7),
		item("Laptop + charger", "Electronics",
			"Essential for presentations and remote work during downtime", true, 4),
		item("Business cards", "Documents & Money",
			"Still expected at many professional events and conferences", false, 6),
	}
}
