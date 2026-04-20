import { Link } from 'react-router-dom'

export default function MyTripsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-violet-700 to-indigo-600 flex items-center justify-center px-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md p-10 text-center">
        <div className="text-5xl mb-4">🗂️</div>
        <h1 className="text-2xl font-bold text-gray-900 mb-2">My Trips</h1>
        <p className="text-gray-500 text-sm mb-1">Your saved trips will appear here soon.</p>
        <p className="text-gray-400 text-sm mb-8">Create your first packing list to get started!</p>
        <Link
          to="/"
          className="inline-block bg-violet-600 hover:bg-violet-700 text-white font-semibold px-6 py-2.5 rounded-lg text-sm transition-colors"
        >
          Create New Trip
        </Link>
      </div>
    </div>
  )
}
