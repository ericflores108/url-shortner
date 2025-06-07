package server

import "net/http"

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.PathValue("shortCode")

	if shortCode == "" {
		http.Error(w, "Short code required", http.StatusBadRequest)
		return
	}

	url, exists, err := s.db.GetURL(ctx, shortCode)

	if err != nil {
		http.Error(w, "Error getting URL", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Short code does not exists", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
