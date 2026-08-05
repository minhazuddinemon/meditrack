package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"meditrack/internal/db"
	"meditrack/templates/pages"
)

type ReceptionHandler struct {
	Queries *db.Queries
}

func (h *ReceptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokens, err := h.Queries.GetTodayTokens(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doctors, err := h.Queries.GetAllDoctors(ctx)
	if err != nil {
		doctors = []db.GetAllDoctorsRow{} // Now matches []db.Doctor!
	}

	pages.Reception(tokens, doctors).Render(ctx, w)
}

func (h *ReceptionHandler) IssueToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stID, err := strconv.Atoi(r.FormValue("st_id"))
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	docID, err := strconv.Atoi(r.FormValue("doc_id"))
	if err != nil {
		http.Error(w, "Please select a doctor", http.StatusBadRequest)
		return
	}

	// 1. Get next token number for today
	nextToken, err := h.Queries.GetNextTokenIDForToday(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var intToken int32 = 1

	switch v := nextToken.(type) {
	case int:
		intToken = int32(v)
	case int64:
		intToken = int32(v)
	case int32:
		intToken = v
	default:
		fmt.Println(fmt.Errorf("something wrong token type encountered"), v)
	}
	// 2. Generate token entry with assigned doctor and WAITING status
	err = h.Queries.GenerateToken(ctx, db.GenerateTokenParams{
		TokenID: intToken,
		StID:    int32(stID),
		DocID:   int32(docID),
	})
	if err != nil {
		http.Error(w, "Error generating token: ensure Student ID and Doctor exist!", http.StatusBadRequest)
		return
	}

	// 3. Return updated queue partial for HTMX swap
	tokens, _ := h.Queries.GetTodayTokens(ctx)
	pages.TodayQueueTable(tokens).Render(ctx, w)
}

func (h *ReceptionHandler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	stID, _ := strconv.Atoi(r.FormValue("st_id"))

	err := h.Queries.CreateStudent(r.Context(), db.CreateStudentParams{
		StID:       int32(stID),
		Name:       r.FormValue("name"),
		Gender:     r.FormValue("gender"),
		Contact:    sql.NullString{String: r.FormValue("contact"), Valid: r.FormValue("contact") != ""},
		Dept:       sql.NullString{String: r.FormValue("dept"), Valid: r.FormValue("dept") != ""},
		BloodGroup: sql.NullString{String: r.FormValue("blood_group"), Valid: r.FormValue("blood_group") != ""},
	})

	if err != nil {
		w.Write([]byte(`<p style="color: red;">Failed to register: Student ID may already exist.</p>`))
		return
	}

	w.Write([]byte(`<p style="color: green;">✅ Student registered successfully!</p>`))
}
