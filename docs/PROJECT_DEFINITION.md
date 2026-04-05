# Project Definition — PackSmart

## Project Name
PackSmart

## Target Audience
Anyone who travels — from weekend visitors to international travelers. No technical skill required; the app should work for anyone who can fill out a form.

## Core Problem
People don't know what to pack for trips. They either forget important items or overpack out of uncertainty. Generic packing lists don't account for actual weather, trip type, or who they're traveling with.

## Value Proposition
PackSmart generates a personalized packing list in seconds by combining your trip details (destination, dates, activities, companions) with the real weather forecast. Instead of "pack a jacket," it tells you *why* — "Pack layers — Paris will be 45–58°F with 2 rainy days during your trip."

## Must-Have Features (MVP)

1. **Trip input form** — destination, dates, travel companions, activities, trip category
2. **Automatic weather forecast** via Open-Meteo (free, no API key) for trips within 16 days; seasonal fallback message for trips further out
3. **Rule-based packing list generation** personalized to weather + trip type + duration + activities
4. **Checklist with checkboxes** to mark items as packed, persisted in MySQL (shareable across devices via trip URL)
5. **Trip categories** that meaningfully change the list (international vs. beach vs. cold weather vs. staying with friends) — ski trips use `cold_weather` trip type with `skiing` activity

## Deferred Features (Post-MVP)

1. User accounts and saving multiple trips
2. Sharing lists with travel companions
3. Shopping links for missing items
4. Claude API "Smart Mode" for more nuanced, conversational suggestions
5. Mobile app (iOS/Android)
6. Custom list templates the user can create and reuse
