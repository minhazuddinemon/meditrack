package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	dbGen "meditrack/internal/db"
	"meditrack/templates/components"
	"meditrack/templates/pages"
)

type TokenHandler struct {
	queries *dbGen.Queries
	db      *sql.DB // Kept for transaction handling if needed
}

func NewTokenHandler(queries *dbGen.Queries, db *sql.DB) *TokenHandler {
	return &TokenHandler{
		queries: queries,
		db:      db,
	}
}

// Render the main token generator page
func (h *TokenHandler) ShowKiosk(w http.ResponseWriter, r *http.Request) {
	pages.KioskPage().Render(r.Context(), w)
}

// Handle POST /tokens
func (h *TokenHandler) IssueToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse student_id from form
	stIDStr := r.FormValue("student_id")
	stID, err := strconv.ParseInt(stIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Student ID format", http.StatusBadRequest)
		return
	}

	// 1. Verify student exists
	student, err := h.queries.GetStudentByID(ctx, int32(stID))
	if err != nil {
		if err == sql.ErrNoRows {
			// Return 200 so HTMX swaps the warning message into #token-result
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<div class="p-4 bg-amber-50 border border-amber-200 text-amber-800 rounded-lg text-center text-sm font-medium">Student ID not found in system!</div>`))
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 2. Prepare today's date as time.Time (required by sqlc)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 3. Fetch next token ID
	rawNextID, err := h.queries.GetNextTokenIDForToday(ctx, today)

	var nextID int32 = 1
	if err == nil && rawNextID != nil {
		// Type-assert the interface{} returned by sqlc
		switch v := rawNextID.(type) {
		case int64:
			nextID = int32(v)
		case int32:
			nextID = v
		case []byte: // MariaDB driver sometimes returns numbers as byte slices
			if parsed, pErr := strconv.ParseInt(string(v), 10, 32); pErr == nil {
				nextID = int32(parsed)
			}
		}
	}

	// 4. Insert token into database
	err = h.queries.GenerateToken(ctx, dbGen.GenerateTokenParams{
		TokenID:   nextID,
		VisitDate: today,
		VisitTime: now,
		StID:      int32(stID),
	})
	if err != nil {
		http.Error(w, "Failed to issue token", http.StatusInternalServerError)
		return
	}

	// 5. Render HTMX component back
	components.TokenDisplay(int64(nextID), student.Name).Render(ctx, w)
}
