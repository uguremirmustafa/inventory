-- name: ListManufacturers :many
SELECT * FROM manufacturer where deleted_at is null;
