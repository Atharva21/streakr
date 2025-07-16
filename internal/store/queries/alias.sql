-- name: AddAliasForHabit :exec
INSERT INTO habit_aliases (habit_id, alias)
VALUES (?, ?);

-- name: GetAllAliasesForHabit :many
SELECT alias
FROM habit_aliases
WHERE habit_id = ?;

-- name: DeleteAliasForHabit :exec
DELETE FROM habit_aliases
WHERE habit_id = ? AND alias = ?;

-- name: DeleteAllAliasesForHabit :exec
DELETE FROM habit_aliases
WHERE habit_id = ?;