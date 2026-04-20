import { Routes, Route } from 'react-router-dom'
import Header from './components/Header.jsx'
import HomePage from './pages/HomePage.jsx'
import PackingListPage from './pages/PackingListPage.jsx'
import AboutPage from './pages/AboutPage.jsx'
import NotFoundPage from './pages/NotFoundPage.jsx'

function App() {
  return (
    <>
      <Header />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/packing-list/:tripId" element={<PackingListPage />} />
        <Route path="/about" element={<AboutPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </>
  )
}

export default App
