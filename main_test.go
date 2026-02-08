package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetArtists() {
	artistsMu.Lock()
	defer artistsMu.Unlock()
	artists = defaultArtists()
}

func TestDeleteArtists(t *testing.T) {
	resetArtists()
	r := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/artists", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	if len(artists) != 0 {
		t.Fatalf("expected artists map to be empty after DELETE /artists, got %d entries", len(artists))
	}
}

func TestGetArtistNotFound(t *testing.T) {
	resetArtists()
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/artists/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestDeleteArtistNotFound(t *testing.T) {
	resetArtists()
	r := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/artists/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestGetArtistsAfterDeleteArtists(t *testing.T) {
	resetArtists()
	r := setupRouter()

	deleteReq := httptest.NewRequest(http.MethodDelete, "/artists", nil)
	deleteRR := httptest.NewRecorder()
	r.ServeHTTP(deleteRR, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, deleteRR.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/artists", nil)
	getRR := httptest.NewRecorder()
	r.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRR.Code)
	}

	var got map[string]Artist
	if err := json.Unmarshal(getRR.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty artists payload after DELETE /artists, got %d entries", len(got))
	}
}
