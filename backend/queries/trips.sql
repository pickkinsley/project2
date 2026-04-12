-- name: InsertTrip :exec
INSERT INTO trips (id, destination, dest_lat, dest_lon, departure_date, return_date, trip_type, companions, activities)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTripByID :one
SELECT id, destination, dest_lat, dest_lon, departure_date, return_date, trip_type, companions, activities, created_at
FROM trips
WHERE id = ?;
