import { Link, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { listTrips } from '../api/trips.js'

function formatDate(dateStr) {
  const [y, m, d] = dateStr.split('-')
  return new Date(y, m - 1, d).toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
}

function TripCard({ trip }) {
  const navigate = useNavigate()

  return (
    <div
      onClick={() => navigate(`/packing-list/${trip.id}`)}
      className="bg-white rounded-2xl shadow-sm border border-transparent hover:border-pink-300 hover:shadow-md cursor-pointer transition-all duration-200 p-6 flex flex-col gap-3"
    >
      {/* Destination */}
      <div>
        <h2 className="text-lg font-bold text-gray-900 leading-snug">{trip.destination}</h2>
        <p className="text-sm text-gray-500 mt-0.5">
          {formatDate(trip.departure_date)} – {formatDate(trip.return_date)}
          <span className="ml-2 text-gray-400">· {trip.duration_days} day{trip.duration_days !== 1 ? 's' : ''}</span>
        </p>
      </div>

      {/* Badges */}
      <div className="flex flex-wrap gap-2">
        <span className="text-xs bg-pink-100 text-pink-700 px-2.5 py-1 rounded-full font-medium capitalize">
          {trip.trip_type.replace(/_/g, ' ')}
        </span>
        <span className="text-xs bg-pink-100 text-pink-700 px-2.5 py-1 rounded-full font-medium capitalize">
          {trip.companions}
        </span>
      </div>

      {/* Footer */}
      <div className="flex items-center justify-between mt-auto pt-2 border-t border-gray-100">
        <span className="text-xs text-gray-400">
          Created {new Date(trip.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
        </span>
        <span className="text-xs font-semibold text-pink-500">View list →</span>
      </div>
    </div>
  )
}

export default function MyTripsPage() {
  const { data: trips, isLoading, isError, error } = useQuery({
    queryKey: ['trips'],
    queryFn: listTrips,
  })

  // ── Loading ──
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-4xl mb-4">🎒</div>
          <p className="text-lg font-medium text-gray-700">Loading your trips…</p>
        </div>
      </div>
    )
  }

  // ── Error ──
  if (isError) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 flex items-center justify-center px-4">
        <div className="bg-white rounded-2xl shadow-2xl p-10 text-center max-w-sm w-full">
          <div className="text-4xl mb-4">😕</div>
          <h1 className="text-xl font-bold text-gray-800 mb-2">Couldn't load trips</h1>
          <p className="text-sm text-gray-500 mb-6">{error?.message ?? 'Something went wrong.'}</p>
          <Link to="/" className="inline-block bg-pink-500 hover:bg-pink-600 text-white text-sm font-semibold px-5 py-2.5 rounded-lg transition-colors">
            ← Create a Trip
          </Link>
        </div>
      </div>
    )
  }

  // ── Empty state ──
  if (!trips || trips.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 flex items-center justify-center px-4">
        <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md p-10 text-center">
          <div className="text-5xl mb-4">🗂️</div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">My Trips</h1>
          <p className="text-gray-500 text-sm mb-1">Your saved trips will appear here.</p>
          <p className="text-gray-400 text-sm mb-8">Create your first packing list to get started!</p>
          <Link
            to="/"
            className="inline-block bg-pink-500 hover:bg-pink-600 text-white font-semibold px-6 py-2.5 rounded-lg text-sm transition-colors"
          >
            Create New Trip
          </Link>
        </div>
      </div>
    )
  }

  // ── Trips grid ──
  return (
    <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 px-4 py-10">
      <div className="max-w-5xl mx-auto">

        {/* Page header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">My Trips</h1>
            <p className="text-gray-500 text-sm mt-1">
              {trips.length} trip{trips.length !== 1 ? 's' : ''} saved
            </p>
          </div>
          <Link
            to="/"
            className="bg-pink-500 hover:bg-pink-600 text-white font-semibold px-5 py-2.5 rounded-lg text-sm transition-colors"
          >
            ➕ Create New Trip
          </Link>
        </div>

        {/* Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {trips.map((trip) => (
            <TripCard key={trip.id} trip={trip} />
          ))}
        </div>

      </div>
    </div>
  )
}
