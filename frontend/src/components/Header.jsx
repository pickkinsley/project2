import { useState } from 'react'
import { Link, NavLink } from 'react-router-dom'

const NAV_LINKS = [
  { to: '/about',     label: 'About' },
  { to: '/',          label: 'Create Trip' },
  { to: '/my-trips',  label: 'My Trips' },
  { to: '/how-it-works', label: 'How It Works' },
]

export default function Header() {
  const [menuOpen, setMenuOpen] = useState(false)

  const linkClass = ({ isActive }) =>
    `text-sm font-medium transition-colors ${
      isActive ? 'text-violet-700' : 'text-gray-600 hover:text-violet-600'
    }`

  return (
    <header className="sticky top-0 z-50 bg-white shadow-sm">
      <div className="max-w-5xl mx-auto px-4 h-14 flex items-center justify-between">

        {/* Brand */}
        <Link to="/" className="flex items-center gap-2 font-bold text-gray-900 hover:text-violet-700 transition-colors">
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
          className="sm:hidden p-2 rounded-md text-gray-600 hover:text-violet-600 hover:bg-violet-50 transition-colors"
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
        <nav className="sm:hidden border-t border-gray-100 bg-white px-4 py-3 flex flex-col gap-1">
          {NAV_LINKS.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                `block px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-violet-50 text-violet-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-violet-600'
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
