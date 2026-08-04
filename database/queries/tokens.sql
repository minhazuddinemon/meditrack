-- name: GetStudentByID :one
SELECT st_id, name, gender, contact, dept, dob, blood_group
FROM Student
WHERE st_id = ?;

-- name: GetNextTokenIDForToday :one
SELECT COALESCE(MAX(token_id) + 1, 1)
FROM Token
WHERE visit_date = ?;

-- name: GenerateToken :exec
INSERT INTO Token (token_id, visit_date, visit_time, st_id)
VALUES (?, ?, ?, ?);

-- name: GetTodayTokens :many
SELECT t.token_id, t.visit_date, t.visit_time, s.st_id, s.name AS student_name, s.dept, s.contact
FROM Token t
JOIN Student s ON t.st_id = s.st_id
WHERE t.visit_date = CURRENT_DATE
ORDER BY t.token_id ASC;
