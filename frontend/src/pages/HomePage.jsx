import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { createTrip } from '../api/trips.js'

const TRIP_TYPES = [
  { value: 'international', label: 'International' },
  { value: 'domestic',      label: 'Domestic' },
  { value: 'beach',         label: 'Beach' },
  { value: 'cold_weather',  label: 'Cold Weather' },
  { value: 'business',      label: 'Business' },
]

const COMPANIONS = [
  { value: 'solo',   label: 'Solo' },
  { value: 'couple', label: 'Couple' },
  { value: 'family', label: 'Family' },
  { value: 'group',  label: 'Group' },
]

const ACTIVITIES = [
  { value: 'sightseeing', label: 'Sightseeing' },
  { value: 'beach',       label: 'Beach' },
  { value: 'hiking',      label: 'Hiking' },
  { value: 'fine_dining', label: 'Fine Dining' },
  { value: 'shopping',    label: 'Shopping' },
  { value: 'nightlife',   label: 'Nightlife' },
]

const EMPTY_FORM = {
  destination: '',
  departure_date: '',
  return_date: '',
  trip_type: '',
  companions: '',
  activities: [],
}

export default function HomePage() {
  const navigate = useNavigate()
  const [form, setForm] = useState(EMPTY_FORM)
  const [errors, setErrors] = useState({})

  const mutation = useMutation({
    mutationFn: createTrip,
    onSuccess: (data) => navigate(`/packing-list/${data.id}`),
  })

  function clearFieldError(field) {
    if (errors[field]) setErrors((e) => ({ ...e, [field]: undefined }))
  }

  function handleChange(e) {
    const { name, value } = e.target
    setForm((f) => ({ ...f, [name]: value }))
    clearFieldError(name)
    mutation.reset()
  }

  function handleActivityToggle(value) {
    setForm((f) => {
      const next = f.activities.includes(value)
        ? f.activities.filter((a) => a !== value)
        : [...f.activities, value]
      return { ...f, activities: next }
    })
    clearFieldError('activities')
    mutation.reset()
  }

  function validate() {
    const e = {}
    if (!form.destination.trim()) e.destination = 'Destination is required.'
    if (!form.departure_date)     e.departure_date = 'Departure date is required.'
    if (!form.return_date)        e.return_date = 'Return date is required.'
    if (form.departure_date && form.return_date && form.return_date <= form.departure_date)
      e.return_date = 'Return date must be after departure date.'
    if (!form.trip_type)          e.trip_type = 'Trip type is required.'
    if (!form.companions)         e.companions = 'Companions is required.'
    if (form.activities.length === 0) e.activities = 'Select at least one activity.'
    return e
  }

  function handleSubmit(e) {
    e.preventDefault()
    const errs = validate()
    if (Object.keys(errs).length > 0) {
      setErrors(errs)
      return
    }
    mutation.mutate(form)
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-violet-700 to-indigo-600 flex items-center justify-center px-4 py-12">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-xl p-8">

        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">🎒 PackSmart</h1>
          <p className="mt-2 text-gray-500 text-sm">Enter your trip details and get a personalized packing list.</p>
        </div>

        {/* API error banner */}
        {mutation.isError && (
          <div className="mb-6 rounded-lg bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-700">
            {mutation.error?.message ?? 'Something went wrong. Please try again.'}
          </div>
        )}

        <form onSubmit={handleSubmit} noValidate className="space-y-5">

          {/* Destination */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Destination <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="destination"
              value={form.destination}
              onChange={handleChange}
              placeholder="e.g., Paris, France"
              className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 ${errors.destination ? 'border-red-400' : 'border-gray-300'}`}
            />
            {errors.destination && <p className="mt-1 text-xs text-red-600">{errors.destination}</p>}
          </div>

          {/* Dates */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Departure Date <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                name="departure_date"
                value={form.departure_date}
                onChange={handleChange}
                className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-violet-500 ${errors.departure_date ? 'border-red-400' : 'border-gray-300'}`}
              />
              {errors.departure_date && <p className="mt-1 text-xs text-red-600">{errors.departure_date}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Return Date <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                name="return_date"
                value={form.return_date}
                onChange={handleChange}
                className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-violet-500 ${errors.return_date ? 'border-red-400' : 'border-gray-300'}`}
              />
              {errors.return_date && <p className="mt-1 text-xs text-red-600">{errors.return_date}</p>}
            </div>
          </div>

          {/* Trip Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Trip Type <span className="text-red-500">*</span>
            </label>
            <select
              name="trip_type"
              value={form.trip_type}
              onChange={handleChange}
              className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-violet-500 bg-white ${errors.trip_type ? 'border-red-400' : 'border-gray-300'}`}
            >
              <option value="">Select trip type…</option>
              {TRIP_TYPES.map((t) => (
                <option key={t.value} value={t.value}>{t.label}</option>
              ))}
            </select>
            {errors.trip_type && <p className="mt-1 text-xs text-red-600">{errors.trip_type}</p>}
          </div>

          {/* Companions */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Traveling With <span className="text-red-500">*</span>
            </label>
            <select
              name="companions"
              value={form.companions}
              onChange={handleChange}
              className={`w-full rounded-lg border px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-violet-500 bg-white ${errors.companions ? 'border-red-400' : 'border-gray-300'}`}
            >
              <option value="">Select companions…</option>
              {COMPANIONS.map((c) => (
                <option key={c.value} value={c.value}>{c.label}</option>
              ))}
            </select>
            {errors.companions && <p className="mt-1 text-xs text-red-600">{errors.companions}</p>}
          </div>

          {/* Activities */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Activities <span className="text-red-500">*</span>
            </label>
            <div className="grid grid-cols-2 gap-2">
              {ACTIVITIES.map((a) => {
                const checked = form.activities.includes(a.value)
                return (
                  <label
                    key={a.value}
                    className={`flex items-center gap-2 rounded-lg border px-3 py-2 text-sm cursor-pointer select-none transition-colors ${checked ? 'border-violet-500 bg-violet-50 text-violet-800' : 'border-gray-200 text-gray-700 hover:border-gray-300'}`}
                  >
                    <input
                      type="checkbox"
                      className="accent-violet-600"
                      checked={checked}
                      onChange={() => handleActivityToggle(a.value)}
                    />
                    {a.label}
                  </label>
                )
              })}
            </div>
            {errors.activities && <p className="mt-1 text-xs text-red-600">{errors.activities}</p>}
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={mutation.isPending}
            className="w-full rounded-lg bg-violet-600 hover:bg-violet-700 disabled:opacity-60 disabled:cursor-not-allowed text-white font-semibold py-3 text-sm transition-colors mt-2"
          >
            {mutation.isPending ? 'Creating…' : 'Generate Packing List'}
          </button>

        </form>
      </div>
    </div>
  )
}
