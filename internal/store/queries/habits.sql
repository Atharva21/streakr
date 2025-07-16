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

-- name: UpdateHabit :exec
UPDATE habits SET name = ?, description = ?, habit_type = ?, last_logged = ?
WHERE id = ?;

-- name: UpdateHabitLastLogged :exec
UPDATE habits SET last_logged = CURRENT_TIMESTAMP
WHERE id = ?;
