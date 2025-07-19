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

-- name: CountTotalImproveHabits :one
SELECT COUNT(*) as total_improve_habits
FROM habits 
WHERE habit_type = 'improve';

-- name: CountImproveHabitsLoggedToday :one
SELECT COUNT(DISTINCT h.id) as logged_today_count
FROM habits h
JOIN streaks s ON h.id = s.habit_id
WHERE h.habit_type = 'improve' 
 AND DATE(s.streak_end) = DATE('now');

-- name: GetDaysSinceHabitCreation :one
SELECT CAST(1 + julianday('now') - julianday(DATE(created_at)) AS INTEGER) as days_passed
FROM habits
WHERE id = ?;
