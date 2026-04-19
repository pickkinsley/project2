const BASE = '/api'

export async function createTrip(formData) {
  const res = await fetch(`${BASE}/trips`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData),
  })
  const data = await res.json()
  if (!res.ok) throw data
  return data
}

export async function getTrip(tripId) {
  const res = await fetch(`${BASE}/trips/${tripId}`)
  const data = await res.json()
  if (!res.ok) throw data
  return data
}

export async function updateItemChecked(tripId, itemId, isChecked) {
  const res = await fetch(`${BASE}/trips/${tripId}/items/${itemId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ is_checked: isChecked }),
  })
  const data = await res.json()
  if (!res.ok) throw data
  return data
}
