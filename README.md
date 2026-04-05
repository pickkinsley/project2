# PackSmart

A smart packing list generator that builds a personalized checklist based on your trip details and real weather forecast.

## What It Does

Enter where you're going, when, who you're with, and what you'll do — PackSmart fetches the actual weather forecast and generates a packing list tailored to your specific trip. Instead of "pack a jacket," it tells you why: *"Paris will be 52–63°F with 2 rainy days during your trip."*

Lists are saved and shareable via URL — anyone with the link can view and check off items.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React + Vite + React Router + TanStack Query |
| Backend | Go |
| Database | MySQL |
| Query layer | sqlc |
| Weather API | Open-Meteo (free, no API key) |
| Deployment | AWS Lightsail + nginx + HTTPS |

## Features

- Trip input: destination, dates, trip type, travel companions, activities
- Real weather forecast for trips within 16 days (Open-Meteo)
- Rule-based packing list: weather × trip type × activities × duration
- Essential items section (passport, prescriptions, power adapter) always shown first
- Shareable URLs — anyone with the link can view the same list
- Checkbox state persists in MySQL across devices

## Project Structure

```
project2/
├── CLAUDE.md          # Implementation guide for Claude Code
├── README.md          # This file
├── docs/
│   ├── PROJECT_PROPOSAL.md    # Complete project proposal
│   ├── PROJECT_DEFINITION.md  # Problem, audience, feature scope
│   ├── DECISIONS.md           # Technical decision log
│   ├── PAGES.md               # Page inventory and states
│   ├── NAVIGATION.md          # Navigation flow diagram
│   ├── COMPONENTS.md          # Component tree
│   ├── DATABASE_SCHEMA.md     # MySQL schema with CREATE TABLE SQL
│   ├── API_ENDPOINTS.md       # 3 endpoints with request/response examples
│   └── RULE_ENGINE.md         # Packing rule framework and execution order
├── backend/           # Go API (Module 5)
└── frontend/          # React app (Module 5)
```

## API Endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/trips` | Create trip, generate packing list |
| `GET` | `/api/trips/{uuid}` | Fetch trip with weather and items |
| `PATCH` | `/api/trips/{uuid}/items/{itemId}` | Check/uncheck an item |

## Pages

| Page | URL | Purpose |
|---|---|---|
| Home | `/` | Trip input form |
| Packing List | `/packing-list/{uuid}` | Checklist, shareable |
| 404 | `*` | Not found handler |

## Status

Module 4 (Planning) complete. Implementation begins in Module 5.
