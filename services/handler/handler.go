package handler

import (
	"encoding/json"
	"net/http"

	"morris-backend.com/main/services/helper"
	"morris-backend.com/main/services/models"
)

func PartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostPartHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetPartHandler(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func PostPartHandler(w http.ResponseWriter, r *http.Request) {

	var part models.Part

	if err := json.NewDecoder(r.Body).Decode(&part); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := helper.PostPart(part.PartNumber, part.RemainPartNumber, part.PartDescription, part.FgWisonPartNumber, part.SuperSSNumber, part.Weight, part.Coo, part.HsCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	part.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(part)
}

func GetPartHandler(w http.ResponseWriter, r *http.Request) {

	part, err := helper.GetPart()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(part)

}

func GetPartHandlerByPartNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract part_number from URL query parameter
	partNumber := r.URL.Query().Get("part_number")
	if partNumber == "" {
		http.Error(w, "part_number parameter is required", http.StatusBadRequest)
		return
	}

	// Retrieve parts from repository
	parts, err := helper.GetPartByPartNumber(partNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(parts) == 0 {
		http.Error(w, "No parts found", http.StatusNotFound)
		return
	}

	// Serialize parts to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(parts)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
