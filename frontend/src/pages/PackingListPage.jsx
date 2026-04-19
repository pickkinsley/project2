import { useParams } from 'react-router-dom'

export default function PackingListPage() {
  const { tripId } = useParams()
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <p className="text-gray-500">PackingListPage — trip {tripId} — coming soon</p>
    </div>
  )
}
