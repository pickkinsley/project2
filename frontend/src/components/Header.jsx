import { useState } from 'react'
import { Link, NavLink } from 'react-router-dom'

const NAV_LINKS = [
  { to: '/about',        label: 'About' },
  { to: '/',             label: 'Create Trip' },
  { to: '/my-trips',     label: 'My Trips' },
  { to: '/how-it-works', label: 'How It Works' },
]

const STRIPE_BG = {
  background:
    'repeating-linear-gradient(90deg, #fecdd3 0px, #fecdd3 24px, #bbf7d0 24px, #bbf7d0 48px)',
}

const TEXT_SHADOW = { textShadow: '0 1px 3px rgba(255,255,255,0.8)' }

export default function Header() {
  const [menuOpen, setMenuOpen] = useState(false)

  const linkClass = ({ isActive }) =>
    `text-sm font-bold transition-colors ${
      isActive ? 'text-black' : 'text-black/70 hover:text-black'
    }`

  return (
    <header className="sticky top-0 z-50 shadow-sm" style={STRIPE_BG}>

      {/* Main nav row */}
      <div className="max-w-5xl mx-auto px-4 h-14 flex items-center justify-between">

        {/* Brand */}
        <Link
          to="/"
          className="flex items-center gap-2 font-bold text-white hover:text-white/80 transition-colors"
          style={TEXT_SHADOW}
        >
          <span className="text-xl">🎒</span>
          <span>PackSmart</span>
        </Link>

        {/* Desktop nav */}
        <nav className="hidden sm:flex items-center gap-6">
          {NAV_LINKS.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={linkClass}
              style={TEXT_SHADOW}
            >
              {label}
            </NavLink>
          ))}
        </nav>

        {/* Mobile hamburger */}
        <button
          className="sm:hidden p-2 rounded-md text-white hover:text-white/80 transition-colors"
          style={TEXT_SHADOW}
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
        <nav className="sm:hidden border-t border-pink-200 px-4 py-3 flex flex-col gap-1" style={STRIPE_BG}>
          {NAV_LINKS.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                `block px-3 py-2 rounded-lg text-sm font-semibold transition-colors ${
                  isActive
                    ? 'bg-white/30 text-black'
                    : 'text-black/70 hover:bg-white/20 hover:text-black'
                }`
              }
              style={TEXT_SHADOW}
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
