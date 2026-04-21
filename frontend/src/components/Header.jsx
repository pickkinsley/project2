import { useState } from 'react'
import { Link, NavLink } from 'react-router-dom'

const NAV_LINKS = [
  { to: '/about',        label: 'About' },
  { to: '/',             label: 'Create Trip' },
  { to: '/my-trips',     label: 'My Trips' },
  { to: '/how-it-works', label: 'How It Works' },
]

export default function Header() {
  const [menuOpen, setMenuOpen] = useState(false)

  const linkClass = ({ isActive }) =>
    `text-sm font-medium transition-colors ${
      isActive ? 'text-pink-600' : 'text-gray-600 hover:text-pink-500'
    }`

  return (
    <header className="sticky top-0 z-50 bg-white shadow-sm">
      {/* Pink + green stripe bar */}
      <div
        className="h-2"
        style={{
          background:
            'repeating-linear-gradient(90deg, #fda4af 0px, #fda4af 20px, #86efac 20px, #86efac 40px)',
        }}
      />

      <div className="max-w-5xl mx-auto px-4 h-13 flex items-center justify-between py-2">

        {/* Brand */}
        <Link to="/" className="flex items-center gap-2 font-bold text-gray-900 hover:text-pink-600 transition-colors">
          <span className="text-xl">🎒</span>
          <span>PackSmart</span>
        </Link>

        {/* Desktop nav */}
        <nav className="hidden sm:flex items-center gap-6">
          {NAV_LINKS.map(({ to, label }) => (
            <NavLink key={to} to={to} end={to === '/'} className={linkClass}>
              {label}
            </NavLink>
          ))}
        </nav>

        {/* Mobile hamburger */}
        <button
          className="sm:hidden p-2 rounded-md text-gray-600 hover:text-pink-500 hover:bg-pink-50 transition-colors"
          onClick={() => setMenuOpen((o) => !o)}
          aria-label="Toggle menu"
        >
          {menuOpen ? (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          ) : (
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          )}
        </button>
      </div>

      {/* Mobile menu */}
      {menuOpen && (
        <nav className="sm:hidden border-t border-pink-100 bg-white px-4 py-3 flex flex-col gap-1">
          {NAV_LINKS.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                `block px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-pink-50 text-pink-600'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-pink-500'
                }`
              }
              onClick={() => setMenuOpen(false)}
            >
              {label}
            </NavLink>
          ))}
        </nav>
      )}
    </header>
  )
}
