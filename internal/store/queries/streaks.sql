-- name: AddStreak :one
INSERT INTO streaks (habit_id, streak_start, streak_end)
VALUES (?, ?, ?)
RETURNING id;

-- name: UpdateStreakEnd :exec
UPDATE streaks
SET streak_end = ?
WHERE id = ?;

-- name: GetLatestStreak :one
SELECT id, habit_id, streak_start, streak_end
FROM streaks
ORDER BY streak_end DESC
LIMIT 1;

-- name: GetLatestStreakForHabit :one
SELECT id, habit_id, streak_start, streak_end
FROM streaks
WHERE habit_id = ?
ORDER BY streak_end DESC
LIMIT 1;

-- name: GetMaxStreak :one
SELECT CAST(COALESCE(MAX(julianday(streak_end) - julianday(streak_start) + 1), 0) AS INTEGER) as max_streak_days
FROM streaks;

-- name: GetMaxStreakQuittingHabit :one
SELECT CAST(COALESCE(MAX(julianday(streak_end) - julianday(streak_start)), 0) AS INTEGER) as max_streak_days
FROM streaks
WHERE habit_id = ?;

-- name: GetMaxStreakForHabit :one
SELECT CAST(COALESCE(MAX(julianday(streak_end) - julianday(streak_start) + 1), 0) AS INTEGER) as max_streak_days
FROM streaks
WHERE habit_id = ?;

-- name: DeleteStreakByID :exec
DELETE FROM streaks
WHERE id = ?;

-- name: DeleteAllStreaksForHabit :exec
DELETE FROM streaks
WHERE habit_id = ?;

-- name: GetStreaksInRange :many
SELECT streak_start, streak_end
FROM streaks
WHERE streak_end >= ? AND streak_start <= ? AND
habit_id = ?
ORDER BY streak_start;
