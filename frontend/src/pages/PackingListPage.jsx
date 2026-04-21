import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getTrip, updateItemChecked } from '../api/trips.js'

// ─── Weather ────────────────────────────────────────────────────────────────

const WEATHER_ICONS = {
  sunny:         '☀️',
  partly_cloudy: '⛅',
  cloudy:        '☁️',
  rainy:         '🌧️',
  snowy:         '❄️',
  stormy:        '⛈️',
}

function WeatherCard({ weather }) {
  if (!weather) {
    return (
      <div className="bg-amber-50 border border-amber-200 rounded-2xl px-5 py-4 text-sm text-amber-800">
        ⚠️ Weather forecast unavailable — showing general recommendations for your trip type.
      </div>
    )
  }

  return (
    <div className="bg-white rounded-2xl shadow-sm p-6">
      <h2 className="font-semibold text-gray-800 mb-4">🌤 Weather Forecast</h2>
      <div className="flex flex-wrap gap-4 text-sm text-gray-600 mb-4">
        <span>🌡 {weather.temp_min_f}°F – {weather.temp_max_f}°F</span>
        {weather.rain_days > 0 && (
          <span>🌧 {weather.rain_days} rainy day{weather.rain_days !== 1 ? 's' : ''}</span>
        )}
        {weather.snow_days > 0 && (
          <span>❄️ {weather.snow_days} snowy day{weather.snow_days !== 1 ? 's' : ''}</span>
        )}
      </div>
      {weather.daily_forecast?.length > 0 && (
        <div className="flex gap-2 overflow-x-auto pb-1">
          {weather.daily_forecast.map((day) => {
            const [y, m, d] = day.date.split('-')
            const label = new Date(y, m - 1, d).toLocaleDateString('en-US', { weekday: 'short' })
            return (
              <div key={day.date} className="flex flex-col items-center min-w-[58px] bg-pink-50 rounded-xl p-3">
                <span className="text-xs text-gray-500 mb-1">{label}</span>
                <span className="text-xl mb-1">{WEATHER_ICONS[day.icon] ?? '🌡'}</span>
                <span className="text-xs text-gray-500">{day.min_f}°</span>
                <span className="text-xs font-semibold text-gray-800">{day.max_f}°</span>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

// ─── Packing list ────────────────────────────────────────────────────────────

function PackingItem({ item, onToggle }) {
  return (
    <label className="flex items-start gap-3 py-3 cursor-pointer group">
      <input
        type="checkbox"
        checked={item.is_checked}
        onChange={() => onToggle(item.id, !item.is_checked)}
        className="mt-0.5 w-4 h-4 accent-pink-500 cursor-pointer shrink-0"
      />
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <span className={`text-sm font-medium transition-colors ${item.is_checked ? 'line-through text-gray-400' : 'text-gray-800'}`}>
            {item.name}
          </span>
          {item.is_essential && (
            <span className="text-xs bg-pink-100 text-pink-700 px-1.5 py-0.5 rounded font-semibold leading-none">
              Essential
            </span>
          )}
        </div>
        {item.reason && (
          <p className="text-xs text-gray-400 mt-0.5 leading-relaxed">{item.reason}</p>
        )}
      </div>
    </label>
  )
}

function CategoryCard({ category, items, onToggle }) {
  const checkedCount = items.filter((i) => i.is_checked).length
  const allDone = checkedCount === items.length

  return (
    <div className="bg-white rounded-2xl shadow-sm p-6">
      <div className="flex items-center justify-between mb-2">
        <h3 className="font-semibold text-gray-800">{category}</h3>
        <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${allDone ? 'bg-green-100 text-green-700' : 'bg-pink-50 text-pink-500'}`}>
          {checkedCount}/{items.length}
        </span>
      </div>
      <div className="divide-y divide-gray-100">
        {items.map((item) => (
          <PackingItem key={item.id} item={item} onToggle={onToggle} />
        ))}
      </div>
    </div>
  )
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

function groupByCategory(items) {
  const groups = {}
  for (const item of items) {
    if (!groups[item.category]) groups[item.category] = []
    groups[item.category].push(item)
  }
  return groups
}

function formatDate(dateStr) {
  const [y, m, d] = dateStr.split('-')
  return new Date(y, m - 1, d).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
}

// ─── Page ─────────────────────────────────────────────────────────────────────

export default function PackingListPage() {
  const { tripId } = useParams()
  const queryClient = useQueryClient()

  const { data: trip, isLoading, isError, error } = useQuery({
    queryKey: ['trip', tripId],
    queryFn: () => getTrip(tripId),
    retry: false,
  })

  const toggleMutation = useMutation({
    mutationFn: ({ itemId, isChecked }) => updateItemChecked(tripId, itemId, isChecked),
    onMutate: async ({ itemId, isChecked }) => {
      await queryClient.cancelQueries({ queryKey: ['trip', tripId] })
      const previous = queryClient.getQueryData(['trip', tripId])
      queryClient.setQueryData(['trip', tripId], (old) => ({
        ...old,
        items: old.items.map((item) =>
          item.id === itemId ? { ...item, is_checked: isChecked } : item
        ),
      }))
      return { previous }
    },
    onError: (_err, _vars, context) => {
      queryClient.setQueryData(['trip', tripId], context.previous)
    },
  })

  function handleToggle(itemId, isChecked) {
    toggleMutation.mutate({ itemId, isChecked })
  }

  // ── Loading ──
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-4xl mb-4">🎒</div>
          <p className="text-lg font-medium text-gray-700">Loading your packing list…</p>
        </div>
      </div>
    )
  }

  // ── Error / not found ──
  if (isError) {
    const msg = error?.message ?? 'Trip not found.'
    return (
      <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 flex items-center justify-center px-4">
        <div className="bg-white rounded-2xl shadow-2xl p-10 text-center max-w-sm w-full">
          <div className="text-4xl mb-4">🔍</div>
          <h1 className="text-xl font-bold text-gray-800 mb-2">Trip Not Found</h1>
          <p className="text-sm text-gray-500 mb-6">{msg}</p>
          <Link to="/" className="inline-block bg-pink-500 hover:bg-pink-600 text-white text-sm font-semibold px-5 py-2.5 rounded-lg transition-colors">
            ← Plan a New Trip
          </Link>
        </div>
      </div>
    )
  }

  // ── Success ──
  const checkedCount = trip.items.filter((i) => i.is_checked).length
  const totalCount = trip.items.length
  const progressPct = totalCount === 0 ? 0 : Math.round((checkedCount / totalCount) * 100)
  const categories = groupByCategory(trip.items)

  return (
    <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 px-4 py-10">
      <div className="max-w-2xl mx-auto space-y-4">

        {/* Back link */}
        <Link to="/" className="inline-flex items-center gap-1 text-rose-700 hover:text-rose-900 text-sm transition-colors">
          ← New Trip
        </Link>

        {/* Trip header */}
        <div className="bg-white rounded-2xl shadow-sm p-6">
          <h1 className="text-2xl font-bold text-gray-900">{trip.destination}</h1>
          <p className="text-sm text-gray-500 mt-1">
            {formatDate(trip.departure_date)} – {formatDate(trip.return_date)}
            <span className="ml-2 text-gray-400">· {trip.duration_days} day{trip.duration_days !== 1 ? 's' : ''}</span>
          </p>
          <div className="flex flex-wrap gap-2 mt-3">
            <span className="text-xs bg-pink-100 text-pink-700 px-2.5 py-1 rounded-full font-medium capitalize">
              {trip.trip_type.replace('_', ' ')}
            </span>
            <span className="text-xs bg-pink-100 text-pink-700 px-2.5 py-1 rounded-full font-medium capitalize">
              {trip.companions}
            </span>
          </div>
        </div>

        {/* Progress */}
        <div className="bg-white rounded-2xl shadow-sm px-6 py-4">
          <div className="flex items-center justify-between text-sm mb-2">
            <span className="font-medium text-gray-700">Packing progress</span>
            <span className="text-gray-500">{checkedCount} of {totalCount} items</span>
          </div>
          <div className="h-2 bg-pink-100 rounded-full overflow-hidden">
            <div
              className="h-full bg-pink-400 rounded-full transition-all duration-300"
              style={{ width: `${progressPct}%` }}
            />
          </div>
          {progressPct === 100 && (
            <p className="text-xs text-green-600 font-medium mt-2">✓ All packed — you're ready to go!</p>
          )}
        </div>

        {/* Weather */}
        <WeatherCard weather={trip.weather} />

        {/* Packing list by category */}
        {Object.entries(categories).map(([category, items]) => (
          <CategoryCard
            key={category}
            category={category}
            items={items}
            onToggle={handleToggle}
          />
        ))}

      </div>
    </div>
  )
}
