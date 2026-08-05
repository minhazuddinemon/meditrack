-- name: CreatePrescription :execresult
INSERT INTO Prescription (st_id, doc_id)
VALUES (?, ?);

-- name: AddSymptom :exec
INSERT INTO Pre_symptoms (p_id, symptom)
VALUES (?, ?);

-- name: AddPrescribedMedicine :exec
INSERT INTO Contain (p_id, medicine_name, medicine_type, dosage)
VALUES (?, ?, ?, ?);

-- name: AddRequiredTest :exec
INSERT INTO Requires (p_id, test_name)
VALUES (?, ?);

-- name: GetPrescriptionDetails :one
SELECT p.p_id, p.p_date, s.st_id, s.name AS student_name, s.dept, s.blood_group, d.name AS doctor_name, d.specialization
FROM Prescription p
JOIN Student s ON p.st_id = s.st_id
JOIN Doctor d ON p.doc_id = d.doc_id
WHERE p.p_id = ?;

-- name: GetPrescriptionSymptoms :many
SELECT symptom FROM Pre_symptoms WHERE p_id = ?;

-- name: GetPrescriptionMedicines :many
SELECT m.medicine_name, m.medicine_type, c.dosage
FROM Contain c
JOIN Medicine m ON c.medicine_name = m.medicine_name AND c.medicine_type = m.medicine_type
WHERE c.p_id = ?;

-- name: GetStudentPrescriptionHistory :many
SELECT p.p_id, p.p_date, d.name AS doctor_name, d.specialization
FROM Prescription p
JOIN Doctor d ON p.doc_id = d.doc_id
WHERE p.st_id = ?
ORDER BY p.p_date DESC;

-- name: GetStudentLabResults :many
SELECT r.p_id, r.test_name, r.test_date, r.test_result, p.p_date
FROM Requires r
JOIN Prescription p ON r.p_id = p.p_id
WHERE p.st_id = ?
ORDER BY p.p_date DESC;

-- name: GetAllMedicines :many
SELECT medicine_name, medicine_type, manufacturer FROM Medicine ORDER BY medicine_name ASC;

-- name: GetAllTests :many
SELECT test_name, test_fee FROM Test ORDER BY test_name ASC;
