package server

import (
	"crypto/rand"
	"encoding/json"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

func (s *Server) shortenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if !isValidURL(req.URL) {
		http.Error(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	existingCode, exists, err := s.db.GetShortCode(ctx, req.URL)
	if err != nil {
		slog.Warn("s.db.URLExists:", "err", err)
		http.Error(w, "Error checking URL", http.StatusInternalServerError)
		return
	}

	var response ShortenResponse
	if exists {
		response = ShortenResponse{
			ShortURL:    s.baseURL + "/" + existingCode,
			ShortCode:   existingCode,
			OriginalURL: req.URL,
		}
	} else {
		shortCode, err := generateShortCode(8)
		if err != nil {
			http.Error(w, "Error generating short code", http.StatusInternalServerError)
			return
		}

		err = s.db.StoreURL(ctx, shortCode, req.URL)
		if err != nil {
			slog.Warn("s.db.StoreURL:", "err", err)
			http.Error(w, "Error storing URL", http.StatusInternalServerError)
			return
		}

		response = ShortenResponse{
			ShortURL:    s.baseURL + "/" + shortCode,
			ShortCode:   shortCode,
			OriginalURL: req.URL,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func isValidURL(urlStr string) bool {
	url, err := url.Parse(urlStr)
	if err != nil {
		slog.Warn("url.Parse", "err", err)
		return false
	}

	if url.Scheme != "https" {
		slog.Warn("invalid url scheme", "scheme", url.Scheme)
		return false
	}

	if url.Host == "" {
		slog.Warn("url.Host", "invalid", url.Host)
		return false
	}

	return true
}

func generateShortCode(length int) (string, error) {
	b := make([]byte, length)
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"

	for i := range len(b) {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}

	return string(b), nil
}
