import { Link } from 'react-router-dom'

const FEATURES = [
  {
    icon: '🌤',
    title: 'Weather-Aware Recommendations',
    description: 'Fetches a real forecast for your destination and adjusts the packing list — umbrellas for rainy days, extra layers for cold snaps.',
  },
  {
    icon: '✈️',
    title: 'Trip-Type Specific Items',
    description: 'International trips add passport and power adapters. Beach trips add sunscreen. Business trips add professional attire. The list matches your trip.',
  },
  {
    icon: '⭐',
    title: 'Essential Items Highlighted',
    description: 'Critical items like prescriptions and travel documents are marked Essential so they stand out and never get skipped.',
  },
  {
    icon: '✅',
    title: 'Interactive Checklist',
    description: 'Check off items as you pack. Progress is saved to the database and the shareable URL lets travel companions see the list too.',
  },
]

const STEPS = [
  { number: '1', text: 'Enter your destination, dates, trip type, and activities.' },
  { number: '2', text: 'PackSmart fetches the weather forecast and runs the packing rules.' },
  { number: '3', text: 'Receive a personalized list, organized by category and priority.' },
  { number: '4', text: 'Check off items as you pack — your progress is saved automatically.' },
]

export default function AboutPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-violet-700 to-indigo-600 px-4 py-12">
      <div className="max-w-2xl mx-auto space-y-6">

        {/* Hero */}
        <div className="text-center text-white py-6">
          <div className="text-5xl mb-4">🎒</div>
          <h1 className="text-4xl font-bold mb-3">PackSmart</h1>
          <p className="text-xl font-medium opacity-90">Smart Packing Lists for Every Trip</p>
          <p className="mt-3 text-white/70 text-sm max-w-md mx-auto leading-relaxed">
            Stop forgetting things. PackSmart builds a personalized packing list in seconds
            based on where you're going, what you're doing, and what the weather will be.
          </p>
        </div>

        {/* What it does */}
        <div className="bg-white rounded-2xl shadow-sm p-6">
          <h2 className="text-lg font-bold text-gray-900 mb-3">What PackSmart Does</h2>
          <ul className="space-y-3 text-sm text-gray-600">
            <li className="flex items-start gap-2">
              <span className="text-violet-500 mt-0.5">▸</span>
              Creates a personalized packing list based on your destination, trip type, companions, and planned activities.
            </li>
            <li className="flex items-start gap-2">
              <span className="text-violet-500 mt-0.5">▸</span>
              Fetches a real weather forecast for your travel dates and adds weather-appropriate items automatically.
            </li>
            <li className="flex items-start gap-2">
              <span className="text-violet-500 mt-0.5">▸</span>
              Highlights essential items — documents, medications, and must-haves — so critical things never get left behind.
            </li>
            <li className="flex items-start gap-2">
              <span className="text-violet-500 mt-0.5">▸</span>
              Saves your list to a unique URL you can revisit or share with travel companions.
            </li>
          </ul>
        </div>

        {/* Features */}
        <div className="bg-white rounded-2xl shadow-sm p-6">
          <h2 className="text-lg font-bold text-gray-900 mb-4">Key Features</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {FEATURES.map((f) => (
              <div key={f.title} className="bg-violet-50 rounded-xl p-4">
                <div className="text-2xl mb-2">{f.icon}</div>
                <h3 className="text-sm font-semibold text-gray-800 mb-1">{f.title}</h3>
                <p className="text-xs text-gray-500 leading-relaxed">{f.description}</p>
              </div>
            ))}
          </div>
        </div>

        {/* How it works */}
        <div className="bg-white rounded-2xl shadow-sm p-6">
          <h2 className="text-lg font-bold text-gray-900 mb-4">How It Works</h2>
          <div className="space-y-3">
            {STEPS.map((s) => (
              <div key={s.number} className="flex items-start gap-3">
                <span className="shrink-0 w-7 h-7 rounded-full bg-violet-600 text-white text-xs font-bold flex items-center justify-center">
                  {s.number}
                </span>
                <p className="text-sm text-gray-600 leading-relaxed pt-0.5">{s.text}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Tech stack */}
        <div className="bg-white/10 rounded-2xl px-6 py-5 text-white/80 text-xs text-center leading-relaxed">
          Built with <span className="font-semibold text-white">Go + Gin</span> backend ·{' '}
          <span className="font-semibold text-white">MySQL</span> database ·{' '}
          <span className="font-semibold text-white">React + Vite</span> frontend ·{' '}
          <span className="font-semibold text-white">TanStack Query</span> for server state ·{' '}
          Deployed on <span className="font-semibold text-white">AWS Lightsail</span>
        </div>

        {/* CTA */}
        <div className="text-center pb-4">
          <Link
            to="/"
            className="inline-block bg-white text-violet-700 font-bold px-8 py-3 rounded-xl shadow-lg hover:bg-violet-50 transition-colors text-sm"
          >
            Create Your Packing List →
          </Link>
        </div>

      </div>
    </div>
  )
}
