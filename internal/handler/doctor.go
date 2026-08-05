package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"meditrack/internal/db"
	"meditrack/internal/services"
	"meditrack/templates/pages"
)

type DoctorHandler struct {
	Queries *db.Queries
	RawDB   *sql.DB
}

func (h *DoctorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if doctor is logged in via cookie
	cookie, err := r.Cookie("doc_id")
	if err != nil || cookie.Value == "" {
		pages.DoctorLogin("").Render(ctx, w)
		return
	}

	docID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		pages.DoctorLogin("Invalid session").Render(ctx, w)
		return
	}

	doc, err := h.Queries.GetDoctorByID(ctx, int32(docID))
	if err != nil {
		pages.DoctorLogin("Doctor account not found").Render(ctx, w)
		return
	}

	tokens, _ := h.Queries.GetDoctorQueueToday(ctx, int32(docID))
	medicines, _ := h.Queries.GetAllMedicines(ctx)
	tests, _ := h.Queries.GetAllTests(ctx)

	pages.DoctorConsole(doc, tokens, medicines, tests).Render(ctx, w)
}

func (h *DoctorHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	docIDStr := r.FormValue("doc_id")
	password := r.FormValue("password")

	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		pages.DoctorLogin("Please enter a valid Doctor ID").Render(ctx, w)
		return
	}

	doc, err := h.Queries.GetDoctorByID(ctx, int32(docID))
	if err != nil || doc.Password != password {
		pages.DoctorLogin("Invalid Doctor ID or Password").Render(ctx, w)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "doc_id",
		Value:    strconv.Itoa(int(doc.DocID)),
		Path:     "/",
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/doctor", http.StatusSeeOther)
}

func (h *DoctorHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "doc_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/doctor", http.StatusSeeOther)
}

func (h *DoctorHandler) LoadPatientWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stIDStr := r.PathValue("st_id")
	stID, err := strconv.Atoi(stIDStr)
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("doc_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	docID, _ := strconv.Atoi(cookie.Value)

	student, err := h.Queries.GetStudentByID(ctx, int32(stID))
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	history, _ := h.Queries.GetStudentPrescriptionHistory(ctx, int32(stID))
	medicines, _ := h.Queries.GetAllMedicines(ctx)
	tests, _ := h.Queries.GetAllTests(ctx)

	pages.PatientWorkspace(student, history, medicines, tests, int32(docID)).Render(ctx, w)
}

func (h *DoctorHandler) GetMedicineRow(w http.ResponseWriter, r *http.Request) {
	medicines, _ := h.Queries.GetAllMedicines(r.Context())
	pages.MedicineRow(medicines).Render(r.Context(), w)
}

func (h *DoctorHandler) CreatePrescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.ParseForm()

	stID, _ := strconv.Atoi(r.FormValue("st_id"))
	docID, _ := strconv.Atoi(r.FormValue("doc_id"))

	rawSymptoms := strings.Split(r.FormValue("symptoms"), ",")
	var symptoms []string
	for _, s := range rawSymptoms {
		symptoms = append(symptoms, strings.TrimSpace(s))
	}

	medKeys := r.Form["medicine_key"]
	dosages := r.Form["dosage"]
	var prescribedMeds []services.PrescribedMed

	for i, key := range medKeys {
		parts := strings.Split(key, "|")
		if len(parts) == 2 {
			dosage := ""
			if i < len(dosages) {
				dosage = dosages[i]
			}
			prescribedMeds = append(prescribedMeds, services.PrescribedMed{
				Name:   parts[0],
				Type:   parts[1],
				Dosage: dosage,
			})
		}
	}

	tests := r.Form["tests"]

	pID, err := services.CreateFullPrescription(ctx, h.RawDB, services.PrescriptionInput{
		StudentID: int32(stID),
		DoctorID:  int32(docID),
		Symptoms:  symptoms,
		Medicines: prescribedMeds,
		Tests:     tests,
	})

	if err != nil {
		w.Write([]byte(`<p style="color: red;">Failed to save Rx: ` + err.Error() + `</p>`))
		return
	}

	// 1. Success message
	w.Write([]byte(fmt.Sprintf(`<p style="color: green; font-weight: bold;">✅ Rx #%d saved & Patient removed from queue!</p>`, pID)))

	// 2. Return updated queue via HTMX Out-of-Band (OOB) Swap to instantly update left panel
	updatedTokens, _ := h.Queries.GetDoctorQueueToday(ctx, int32(docID))
	w.Write([]byte(`<div id="doctor-queue" hx-swap-oob="true">`))
	pages.DoctorQueueList(updatedTokens).Render(ctx, w)
	w.Write([]byte(`</div>`))
}
