package services

import (
	"context"
	"database/sql"
	"fmt"

	"meditrack/internal/db"
)

type PrescriptionInput struct {
	StudentID int32
	DoctorID  int32
	Symptoms  []string
	Medicines []PrescribedMed
	Tests     []string
}

type PrescribedMed struct {
	Name   string
	Type   string
	Dosage string
}

func CreateFullPrescription(ctx context.Context, rawDB *sql.DB, input PrescriptionInput) (int64, error) {
	tx, err := rawDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := db.New(rawDB).WithTx(tx)

	// 1. Create Prescription record
	res, err := qtx.CreatePrescription(ctx, db.CreatePrescriptionParams{
		StID:  input.StudentID,
		DocID: input.DoctorID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create prescription: %w", err)
	}

	pID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get prescription id: %w", err)
	}

	// 2. Insert Pre_symptoms
	for _, sym := range input.Symptoms {
		if sym == "" {
			continue
		}
		err := qtx.AddSymptom(ctx, db.AddSymptomParams{
			PID:     int32(pID),
			Symptom: sym,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to add symptom '%s': %w", sym, err)
		}
	}

	// 3. Insert Prescribed Medicines (Contain)
	for _, med := range input.Medicines {
		if med.Name == "" || med.Type == "" {
			continue
		}
		err := qtx.AddPrescribedMedicine(ctx, db.AddPrescribedMedicineParams{
			PID:          int32(pID),
			MedicineName: med.Name,
			MedicineType: med.Type,
			Dosage:       sql.NullString{String: med.Dosage, Valid: med.Dosage != ""},
		})
		if err != nil {
			return 0, fmt.Errorf("failed to add medicine '%s': %w", med.Name, err)
		}
	}

	// 4. Insert Required Lab Tests (Requires)
	for _, test := range input.Tests {
		if test == "" {
			continue
		}
		err := qtx.AddRequiredTest(ctx, db.AddRequiredTestParams{
			PID:      int32(pID),
			TestName: test,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to add test '%s': %w", test, err)
		}
	}

	// 5. Mark Token as COMPLETED (removes from active queue)
	err = qtx.CompleteTokenStatus(ctx, db.CompleteTokenStatusParams{
		StID:  input.StudentID,
		DocID: input.DoctorID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to complete token status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return pID, nil
}
