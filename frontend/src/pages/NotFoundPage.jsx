import { Link } from 'react-router-dom'

export default function NotFoundPage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-4 bg-gradient-to-br from-rose-200 via-pink-100 to-rose-50">
      <h1 className="text-4xl font-bold text-gray-800">404</h1>
      <p className="text-gray-500">Page not found.</p>
      <Link to="/" className="text-pink-600 hover:underline">← Back to home</Link>
    </div>
  )
}
