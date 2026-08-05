-- name: GetStudentByID :one
SELECT st_id, name, gender, contact, dept, dob, blood_group
FROM Student
WHERE st_id = ?;

-- name: CreateStudent :exec
INSERT INTO Student (st_id, name, gender, contact, dept, dob, blood_group)
VALUES (?, ?, ?, ?, ?, ?, ?);


-- name: GetAllDoctors :many
SELECT doc_id, name, specialization, contact
FROM Doctor
ORDER BY name ASC;


-- name: GetNextTokenIDForToday :one
SELECT COALESCE(MAX(token_id) + 1, 1)
FROM Token
WHERE visit_date = CURRENT_DATE;

-- name: GenerateToken :exec
INSERT INTO Token (token_id, visit_date, visit_time, st_id, doc_id, status)
VALUES (?, CURRENT_DATE, CURRENT_TIME, ?, ?, 'WAITING');

-- name: GetTodayTokens :many
SELECT
    t.token_id,
    t.visit_date,
    TIME_FORMAT(t.visit_time, '%H:%i') AS visit_time,
    t.status,
    s.st_id,
    s.name AS student_name,
    s.dept,
    s.contact,
    d.doc_id,
    d.name AS doctor_name,
    d.specialization
FROM Token t
JOIN Student s ON t.st_id = s.st_id
JOIN Doctor d ON t.doc_id = d.doc_id
WHERE t.visit_date = CURRENT_DATE
ORDER BY t.token_id ASC;

-- name: UpdateTokenStatus :exec
UPDATE Token
SET status = ?
WHERE token_id = ? AND visit_date = CURRENT_DATE;
