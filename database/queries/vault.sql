-- name: GetStudentPrescriptionList :many
SELECT p.p_id, p.p_date, d.name AS doctor_name, d.specialization
FROM Prescription p
JOIN Doctor d ON p.doc_id = d.doc_id
WHERE p.st_id = ?
ORDER BY p.p_date DESC, p.p_id DESC;

-- name: GetPrescriptionFullHeader :one
SELECT p.p_id, p.p_date,
       s.st_id, s.name AS student_name, s.dept, s.gender, s.blood_group, s.dob,
       d.name AS doctor_name, d.specialization, d.contact AS doctor_contact
FROM Prescription p
JOIN Student s ON p.st_id = s.st_id
JOIN Doctor d ON p.doc_id = d.doc_id
WHERE p.p_id = ?;

-- name: GetPrescriptionLabTests :many
SELECT r.test_name, r.test_date, r.test_result, t.test_fee
FROM Requires r
JOIN Test t ON r.test_name = t.test_name
WHERE r.p_id = ?;
