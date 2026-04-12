-- name: InsertWeatherSnapshot :exec
INSERT INTO weather_snapshots (trip_id, temp_min_f, temp_max_f, rain_days, snow_days, is_forecast, daily_forecast)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetWeatherByTripID :one
SELECT id, trip_id, temp_min_f, temp_max_f, rain_days, snow_days, is_forecast, daily_forecast, fetched_at
FROM weather_snapshots
WHERE trip_id = ?;
