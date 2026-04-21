const BASE = '/api'

async function parseResponse(res) {
  const text = await res.text()
  if (!text) {
    if (res.ok) return {}
    throw new Error(`Server error: ${res.status}`)
  }
  try {
    const data = JSON.parse(text)
    if (!res.ok) throw data
    return data
  } catch (e) {
    if (!res.ok) throw new Error(`Server error: ${res.status}`)
    throw e
  }
}

export async function createTrip(formData) {
  const res = await fetch(`${BASE}/trips`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData),
  })
  return parseResponse(res)
}

export async function getTrip(tripId) {
  const res = await fetch(`${BASE}/trips/${tripId}`)
  return parseResponse(res)
}

export async function updateTrip(tripId, formData) {
  const res = await fetch(`${BASE}/trips/${tripId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData),
  })
  return parseResponse(res)
}

export async function deleteTrip(tripId) {
  const res = await fetch(`${BASE}/trips/${tripId}`, {
    method: 'DELETE',
  })
  return parseResponse(res)
}

export async function createPackingItem(tripId, itemData) {
  const res = await fetch(`${BASE}/trips/${tripId}/items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(itemData),
  })
  return parseResponse(res)
}

export async function updateItemChecked(tripId, itemId, isChecked) {
  const res = await fetch(`${BASE}/trips/${tripId}/items/${itemId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ is_checked: isChecked }),
  })
  return parseResponse(res)
}
