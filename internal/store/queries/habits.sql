-- name: AddHabit :one
INSERT INTO habits (name, description, habit_type, created_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
RETURNING id;

-- name: GetHabit :one
SELECT * FROM habits WHERE id = ?;

-- name: GetHabitByName :one
SELECT * FROM habits WHERE name = ?;

-- name: ListHabits :many
SELECT * FROM habits;

-- name: DeleteHabit :exec
DELETE FROM habits WHERE id = ?;

-- name: DeleteHabitByName :exec
DELETE FROM habits WHERE name = ?;
