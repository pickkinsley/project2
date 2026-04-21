import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getTrip, updateItemChecked, deleteTrip, createPackingItem } from '../api/trips.js'

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

// ─── Add Item Form ────────────────────────────────────────────────────────────

const CATEGORIES = [
  'Essential Items',
  'Clothing',
  'Toiletries',
  'Electronics',
  'Documents',
  'Activity Specific',
  'Beach Gear',
  'Cold Weather Gear',
  'Other',
]

function AddItemForm({ tripId, onSuccess, onCancel }) {
  const queryClient = useQueryClient()
  const [name, setName] = useState('')
  const [category, setCategory] = useState('')
  const [errors, setErrors] = useState({})

  const addMutation = useMutation({
    mutationFn: (itemData) => createPackingItem(tripId, itemData),
    onMutate: async (itemData) => {
      await queryClient.cancelQueries({ queryKey: ['trip', tripId] })
      const previous = queryClient.getQueryData(['trip', tripId])
      // Optimistically add with a temporary negative ID to avoid collisions
      const tempItem = {
        id: -Date.now(),
        name: itemData.name,
        category: itemData.category,
        is_essential: false,
        reason: '',
        is_checked: false,
        sort_order: 9999,
      }
      queryClient.setQueryData(['trip', tripId], (old) => ({
        ...old,
        items: [...old.items, tempItem],
      }))
      return { previous }
    },
    onError: (_err, _vars, context) => {
      queryClient.setQueryData(['trip', tripId], context.previous)
    },
    onSuccess: (serverItem) => {
      // Replace the temp item with the real one from the server
      queryClient.setQueryData(['trip', tripId], (old) => ({
        ...old,
        items: old.items.map((item) =>
          item.id < 0 ? { ...serverItem } : item
        ),
      }))
      onSuccess()
    },
  })

  function validate() {
    const e = {}
    const trimmed = name.trim()
    if (!trimmed) e.name = 'Item name is required.'
    else if (trimmed.length > 100) e.name = 'Item name must be 100 characters or fewer.'
    if (!category) e.category = 'Category is required.'
    return e
  }

  function handleSubmit(e) {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length > 0) {
      setErrors(errs)
      return
    }
    addMutation.mutate({ name: name.trim(), category })
  }

  return (
    <div className="bg-white rounded-2xl shadow-sm p-6">
      <h3 className="font-semibold text-gray-800 mb-4">Add Custom Item</h3>

      {addMutation.isError && (
        <div className="mb-4 rounded-lg bg-red-50 border border-red-200 px-4 py-2.5 text-sm text-red-700">
          {addMutation.error?.message ?? 'Failed to add item. Please try again.'}
        </div>
      )}

      <form onSubmit={handleSubmit} noValidate className="space-y-3">
        <div>
          <input
            type="text"
            value={name}
            onChange={(e) => {
              setName(e.target.value)
              if (errors.name) setErrors((prev) => ({ ...prev, name: undefined }))
            }}
            placeholder="e.g., Beach towel, Sunscreen..."
            maxLength={100}
            className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-pink-400 ${errors.name ? 'border-red-400' : 'border-gray-300'}`}
          />
          {errors.name && <p className="mt-1 text-xs text-red-600">{errors.name}</p>}
        </div>

        <div>
          <select
            value={category}
            onChange={(e) => {
              setCategory(e.target.value)
              if (errors.category) setErrors((prev) => ({ ...prev, category: undefined }))
            }}
            className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 bg-white focus:outline-none focus:ring-2 focus:ring-pink-400 ${errors.category ? 'border-red-400' : 'border-gray-300'}`}
          >
            <option value="">Select category…</option>
            {CATEGORIES.map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>
          {errors.category && <p className="mt-1 text-xs text-red-600">{errors.category}</p>}
        </div>

        <div className="flex gap-2 pt-1">
          <button
            type="submit"
            disabled={addMutation.isPending}
            className="flex-1 rounded-lg bg-pink-400 hover:bg-pink-500 disabled:opacity-60 disabled:cursor-not-allowed text-white font-semibold py-2 text-sm transition-colors"
          >
            {addMutation.isPending ? 'Adding…' : 'Add Item'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="flex-1 rounded-lg border border-gray-300 text-gray-600 hover:bg-gray-50 font-semibold py-2 text-sm transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
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
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [showAddForm, setShowAddForm] = useState(false)

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

  const deleteMutation = useMutation({
    mutationFn: () => deleteTrip(tripId),
    onSuccess: () => navigate('/'),
    onError: (err) => alert(err?.message ?? 'Failed to delete trip. Please try again.'),
  })

  function handleToggle(itemId, isChecked) {
    toggleMutation.mutate({ itemId, isChecked })
  }

  function handleDelete() {
    if (window.confirm('Are you sure you want to delete this trip? This cannot be undone.')) {
      deleteMutation.mutate()
    }
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

          <div className="flex gap-2 mt-4">
            <Link
              to={`/edit-trip/${tripId}`}
              className="flex items-center gap-1.5 px-4 py-2 rounded-lg bg-pink-50 border border-pink-300 text-pink-700 text-sm font-semibold hover:bg-pink-100 transition-colors"
            >
              ✏️ Edit Trip
            </Link>
            <button
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
              className="flex items-center gap-1.5 px-4 py-2 rounded-lg border border-rose-300 text-rose-600 text-sm font-semibold hover:bg-rose-50 disabled:opacity-60 disabled:cursor-not-allowed transition-colors"
            >
              {deleteMutation.isPending ? 'Deleting…' : '🗑️ Delete'}
            </button>
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

        {/* Add Item button / form */}
        {showAddForm ? (
          <AddItemForm
            tripId={tripId}
            onSuccess={() => setShowAddForm(false)}
            onCancel={() => setShowAddForm(false)}
          />
        ) : (
          <button
            onClick={() => setShowAddForm(true)}
            className="w-full rounded-2xl bg-pink-400 hover:bg-pink-500 text-white font-semibold py-3 text-sm transition-colors"
          >
            ➕ Add Custom Item
          </button>
        )}

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
