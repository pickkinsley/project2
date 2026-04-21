import { Link } from 'react-router-dom'

const STEPS = [
  {
    number: '1',
    title: 'Enter Your Trip Details',
    description: 'Tell us where you\'re going, when, and what you\'ll be doing. Destination, dates, trip type, companions, and activities all shape your list.',
  },
  {
    number: '2',
    title: 'We Analyze Your Destination',
    description: 'PackSmart checks the weather forecast for your travel dates and factors in your trip type and activities to determine exactly what you\'ll need.',
  },
  {
    number: '3',
    title: 'Get Your Personalized List',
    description: 'Receive a smart packing list tailored to your trip — organized by category, with essential items highlighted so nothing critical gets missed.',
  },
  {
    number: '4',
    title: 'Check Off As You Pack',
    description: 'Use the interactive checklist to tick off items as you pack. Your progress saves automatically and the list is accessible from any device via its unique URL.',
  },
]

export default function HowItWorksPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50 px-4 py-12">
      <div className="max-w-xl mx-auto space-y-5">

        {/* Header */}
        <div className="text-center py-4">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">How PackSmart Works</h1>
          <p className="text-gray-500 text-sm">From trip details to a ready-to-pack checklist in seconds.</p>
        </div>

        {/* Steps */}
        {STEPS.map((step) => (
          <div key={step.number} className="bg-white rounded-2xl shadow-sm p-6 flex items-start gap-4">
            <span className="shrink-0 w-10 h-10 rounded-full bg-emerald-400 text-white font-bold text-lg flex items-center justify-center">
              {step.number}
            </span>
            <div>
              <h2 className="font-semibold text-gray-900 mb-1">{step.title}</h2>
              <p className="text-sm text-gray-500 leading-relaxed">{step.description}</p>
            </div>
          </div>
        ))}

        {/* CTA */}
        <div className="text-center pt-2 pb-4">
          <Link
            to="/"
            className="inline-block bg-emerald-400 hover:bg-emerald-500 text-white font-bold px-8 py-3 rounded-xl shadow-lg transition-colors text-sm"
          >
            Get Started →
          </Link>
        </div>

      </div>
    </div>
  )
}
