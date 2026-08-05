-- name: GetDoctorByID :one
SELECT doc_id, name, specialization, contact, password
FROM Doctor
WHERE doc_id = ?;

-- name: GetDoctorQueueToday :many
SELECT
    t.token_id,
    t.visit_date,
    TIME_FORMAT(t.visit_time, '%H:%i') AS visit_time,
    t.status,
    s.st_id,
    s.name AS student_name,
    s.dept,
    s.contact
FROM Token t
JOIN Student s ON t.st_id = s.st_id
WHERE t.visit_date = CURRENT_DATE
  AND t.doc_id = ?
  AND t.status = 'WAITING'
ORDER BY t.token_id ASC;

-- name: CompleteTokenStatus :exec
UPDATE Token
SET status = 'COMPLETED'
WHERE st_id = ? AND visit_date = CURRENT_DATE AND doc_id = ? AND status = 'WAITING';
