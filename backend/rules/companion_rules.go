package rules

import (
	"fmt"

	"github.com/pickkinsley/project2/backend/models"
)

// companionItems adds items based on who is traveling.
func companionItems(ctx TripContext) []models.PackingItem {
	switch ctx.Companions {
	case "family":
		return familyItems(ctx)
	case "couple":
		return coupleItems(ctx)
	case "group":
		return groupItems()
	default:
		// solo — no additional companion-specific items beyond the base list
		return nil
	}
}

func familyItems(ctx TripContext) []models.PackingItem {
	items := []models.PackingItem{
		// ── Kids essentials ────────────────────────────────────────────────
		item("Extra changes of clothes for kids", "Clothing",
			"Kids go through clothes faster than adults — pack at least 2 extra outfits per child", true, 16),
		item("Baby wipes / wet wipes", "Health & Safety",
			"Indispensable for quick cleanups regardless of kids' ages", true, 7),
		item("Children's pain reliever / fever reducer", "Health & Safety",
			"Acetaminophen or ibuprofen formulated for children's weight — travel with it, not without it", true, 8),
		item("Children's sunscreen SPF 50", "Health & Safety",
			"Kids' skin burns faster — use dedicated children's formula with no harsh chemicals", true, 9),
		item("Insect repellent (child-safe formula)", "Health & Safety",
			"DEET-free formula for young kids — check age guidelines on the label", false, 10),
		item("Hand sanitizer + hand wipes", "Health & Safety",
			"Kids touch everything — quick sanitation before meals is essential", false, 11),

		// ── Entertainment ──────────────────────────────────────────────────
		item("Travel games (card games, small board games)", "Activity Specific",
			"Keep everyone engaged on planes, trains, and restaurant waits", false, 13),
		item("Kids' entertainment (tablet loaded with shows/games)", "Electronics",
			"Download content before you leave — in-flight Wi-Fi is expensive and unreliable", false, 9),
		item("Headphones for kids", "Electronics",
			"Volume-limiting headphones protect kids' hearing on flights", false, 10),
		item("Coloring books + crayons / colored pencils", "Activity Specific",
			"Screen-free entertainment for restaurants and long waits", false, 14),
		item("Small comfort item / stuffed animal", "Essential Items",
			"A familiar toy makes unfamiliar places feel safe for young kids", false, 14),

		// ── Logistics ─────────────────────────────────────────────────────
		item("Stroller or carrier (age-appropriate)", "Essential Items",
			"For young children who can't walk long distances yet — check airline policies", false, 15),
		item("Portable first aid kit", "Health & Safety",
			"Band-aids, antiseptic, tweezers, thermometer — more complete than the base kit", true, 12),
		item("Snacks (kid-approved, non-perishable)", "Activity Specific",
			"Hungry kids derail everyone's plans — always carry something they'll actually eat", true, 15),
		item("Reusable water bottles for kids", "Activity Specific",
			"Staying hydrated is harder to enforce with kids — keep bottles accessible", true, 16),
	}

	// For longer family trips, add a portable sound machine and laundry bags
	if ctx.DurationDays >= 5 {
		items = append(items,
			item("Portable white noise machine", "Electronics",
				"Helps kids sleep in unfamiliar hotel rooms — also blocks corridor noise", false, 11),
			item("Laundry bags (separate dirty clothes)", "Essential Items",
				"Keep dirty and clean clothes separated in a shared suitcase", false, 16),
		)
	}

	// Night items for trips with young children
	items = append(items,
		item("Nightlight (portable)", "Electronics",
			"Unfamiliar dark rooms are scary for small kids — a compact nightlight costs nothing", false, 12),
		item(fmt.Sprintf("Diapers / pull-ups (%d days supply)", ctx.DurationDays), "Essential Items",
			"Bring enough for the trip plus a buffer — availability varies by destination", true, 17),
	)

	return items
}

func coupleItems(ctx TripContext) []models.PackingItem {
	items := []models.PackingItem{
		item("Nice outfit for a special dinner", "Clothing",
			"At least one nicer option for a romantic dinner out — you'll want it", false, 13),
		item("Camera or phone camera accessories", "Electronics",
			"Capture the memories together — a small tripod means you don't need to ask strangers", false, 9),
	}

	// For longer couple trips, add connection and comfort items
	if ctx.DurationDays >= 5 {
		items = append(items,
			item("Travel games for two (card game or travel version)", "Activity Specific",
				"Something to do on slow evenings or long transit — Taco Cat Goat Cheese Pizza, Bananagrams, etc.", false, 15),
		)
	}

	return items
}

func groupItems() []models.PackingItem {
	return []models.PackingItem{
		item("Portable Bluetooth speaker", "Electronics",
			"Great for group hangouts, beach days, and pre-dinner gatherings", false, 10),
		item("Group activity supplies (cards, etc.)", "Activity Specific",
			"A card game works anywhere and brings everyone together — UNO, Cabo, etc.", false, 16),
		item("Shared snack supply", "Activity Specific",
			"Buy snacks as a group on arrival — cheaper and keeps everyone fueled", false, 17),
		item("Power strip with USB ports", "Electronics",
			"One hotel room, multiple people charging devices — a shared power strip solves it", false, 11),
	}
}
