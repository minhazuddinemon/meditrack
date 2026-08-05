package handler

import (
	"net/http"
	"strconv"

	"meditrack/internal/db"
	"meditrack/templates/pages"
)

type VaultHandler struct {
	Queries *db.Queries
}

func (h *VaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.VaultSearchPage().Render(r.Context(), w)
}

func (h *VaultHandler) GetStudentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stIDStr := r.FormValue("st_id")

	stID, err := strconv.Atoi(stIDStr)
	if err != nil {
		w.Write([]byte(`<article style="color: red;"><p>Please enter a valid numeric Student ID.</p></article>`))
		return
	}

	student, err := h.Queries.GetStudentByID(ctx, int32(stID))
	if err != nil {
		w.Write([]byte(`<article style="color: red;"><p>No student found with ID ` + stIDStr + `.</p></article>`))
		return
	}

	history, _ := h.Queries.GetStudentPrescriptionList(ctx, int32(stID))
	labResults, _ := h.Queries.GetStudentLabResults(ctx, int32(stID))

	pages.VaultStudentHistory(student, history, labResults).Render(ctx, w)
}

func (h *VaultHandler) ViewPrescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pIDStr := r.PathValue("p_id")

	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		http.Error(w, "Invalid Prescription ID", http.StatusBadRequest)
		return
	}

	header, err := h.Queries.GetPrescriptionFullHeader(ctx, int32(pID))
	if err != nil {
		http.Error(w, "Prescription not found", http.StatusNotFound)
		return
	}

	symptoms, _ := h.Queries.GetPrescriptionSymptoms(ctx, int32(pID))
	medicines, _ := h.Queries.GetPrescriptionMedicines(ctx, int32(pID))
	tests, _ := h.Queries.GetPrescriptionLabTests(ctx, int32(pID))

	pages.PrescriptionPrintDocument(header, symptoms, medicines, tests).Render(ctx, w)
}
