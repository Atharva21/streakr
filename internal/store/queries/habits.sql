-- name: AddHabit :exec
INSERT INTO habits (name, description, habit_type, created_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP);

-- name: GetHabit :one
SELECT * FROM habits WHERE id = ?;

-- name: GetHabitByName :one
SELECT * FROM habits WHERE name = ?;

-- name: GetHabitByAlias :one
SELECT h.*
FROM habits h
JOIN habit_aliases ha ON h.id = ha.habit_id
WHERE ha.alias = ?;

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
