-- name: InsertPackingItem :exec
INSERT INTO packing_items (trip_id, name, category, is_essential, reason, is_checked, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetItemsByTripID :many
SELECT id, trip_id, name, category, is_essential, reason, is_checked, sort_order, created_at
FROM packing_items
WHERE trip_id = ?
ORDER BY sort_order ASC;

-- name: UpdateItemChecked :exec
UPDATE packing_items
SET is_checked = ?
WHERE id = ? AND trip_id = ?;

-- name: CreatePackingItem :execlastid
INSERT INTO packing_items (trip_id, name, category, is_essential, reason, is_checked, sort_order)
VALUES (?, ?, ?, false, NULL, false, 9999);

-- name: GetPackingItemByID :one
SELECT id, trip_id, name, category, is_essential, reason, is_checked, sort_order, created_at
FROM packing_items
WHERE id = ? AND trip_id = ?;
