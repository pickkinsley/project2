-- name: InsertTrip :exec
INSERT INTO trips (id, destination, dest_lat, dest_lon, departure_date, return_date, trip_type, companions, activities)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTripByID :one
SELECT id, destination, dest_lat, dest_lon, departure_date, return_date, trip_type, companions, activities, created_at
FROM trips
WHERE id = ?;

-- name: UpdateTrip :exec
UPDATE trips
SET destination = ?, dest_lat = ?, dest_lon = ?,
    departure_date = ?, return_date = ?,
    trip_type = ?, companions = ?, activities = ?
WHERE id = ?;

-- name: DeleteTrip :exec
DELETE FROM trips WHERE id = ?;
