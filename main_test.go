package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetArtists(t *testing.T) {
	artistsMu.Lock()
	defer artistsMu.Unlock()
	artists = defaultArtists()
}

func TestGetArtists(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getArtists)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var artists map[string]Artist
	err = json.Unmarshal(rr.Body.Bytes(), &artists)
	if err != nil {
		t.Errorf("response body is not valid JSON: %v", err)
	}

	if len(artists) != 2 {
		t.Errorf("expected 2 artists, got %d", len(artists))
	}

	if artists["1"].Name != "30 Seconds To Mars" {
		t.Errorf("artist 1 name mismatch: got %v want '30 Seconds To Mars'", artists["1"].Name)
	}

	if artists["2"].Name != "Garbage" {
		t.Errorf("artist 2 name mismatch: got %v want 'Garbage'", artists["2"].Name)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header mismatch: got %v want 'application/json'", rr.Header().Get("Content-Type"))
	}
}

func TestGetArtist(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("GET", "/artists/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var artist Artist
	err = json.Unmarshal(rr.Body.Bytes(), &artist)
	if err != nil {
		t.Errorf("response body is not valid JSON: %v", err)
	}

	if artist.ID != "1" {
		t.Errorf("artist ID mismatch: got %v want '1'", artist.ID)
	}

	if artist.Name != "30 Seconds To Mars" {
		t.Errorf("artist name mismatch: got %v want '30 Seconds To Mars'", artist.Name)
	}

	if artist.Genre != "alternative" {
		t.Errorf("artist genre mismatch: got %v want 'alternative'", artist.Genre)
	}

	if len(artist.Songs) != 4 {
		t.Errorf("expected 4 songs, got %d", len(artist.Songs))
	}
}

func TestGetArtistNotFound(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("GET", "/artists/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestPostArtist(t *testing.T) {
	resetArtists(t)

	newArtist := Artist{
		ID:    "3",
		Name:  "The Beatles",
		Born:  "1960",
		Genre: "rock",
		Songs: []string{"Hey Jude", "Let It Be", "Yesterday"},
	}

	body, err := json.Marshal(newArtist)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/artists", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Verify the artist was actually added
	artistsMu.RLock()
	addedArtist, exists := artists["3"]
	artistsMu.RUnlock()

	if !exists {
		t.Errorf("artist was not added to the map")
	}

	if addedArtist.Name != "The Beatles" {
		t.Errorf("added artist name mismatch: got %v want 'The Beatles'", addedArtist.Name)
	}
}

func TestPostArtistInvalidJSON(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("POST", "/artists", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostArtistEmptyBody(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("POST", "/artists", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestDeleteArtist(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("DELETE", "/artists/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the artist was actually deleted
	artistsMu.RLock()
	_, exists := artists["1"]
	artistsMu.RUnlock()

	if exists {
		t.Errorf("artist was not deleted from the map")
	}

	// Verify other artists still exist
	artistsMu.RLock()
	_, exists = artists["2"]
	artistsMu.RUnlock()

	if !exists {
		t.Errorf("other artists should not be deleted")
	}
}

func TestDeleteArtistNotFound(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("DELETE", "/artists/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDeleteArtists(t *testing.T) {
	resetArtists(t)

	req, err := http.NewRequest("DELETE", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify all artists were deleted
	artistsMu.RLock()
	artistCount := len(artists)
	artistsMu.RUnlock()

	if artistCount != 0 {
		t.Errorf("expected 0 artists after delete all, got %d", artistCount)
	}
}

func TestGetArtistsAfterDeleteArtists(t *testing.T) {
	resetArtists(t)

	// Delete all artists
	deleteReq, err := http.NewRequest("DELETE", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	deleteRr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(deleteRr, deleteReq)

	// Now get all artists
	getReq, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	getRr := httptest.NewRecorder()
	router.ServeHTTP(getRr, getReq)

	if status := getRr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resultArtists map[string]Artist
	err = json.Unmarshal(getRr.Body.Bytes(), &resultArtists)
	if err != nil {
		t.Errorf("response body is not valid JSON: %v", err)
	}

	if len(resultArtists) != 0 {
		t.Errorf("expected empty artists map, got %d artists", len(resultArtists))
	}
}

func TestUpdateArtistByPost(t *testing.T) {
	resetArtists(t)

	// Update existing artist with new songs
	updatedArtist := Artist{
		ID:    "1",
		Name:  "30 Seconds To Mars",
		Born:  "1998",
		Genre: "alternative rock",
		Songs: []string{"The Kill", "A Beautiful Lie", "Attack", "Live Like A Dream", "City"},
	}

	body, err := json.Marshal(updatedArtist)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/artists", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Verify the artist was updated
	artistsMu.RLock()
	modifiedArtist := artists["1"]
	artistsMu.RUnlock()

	if modifiedArtist.Born != "1998" {
		t.Errorf("artist born year not updated: got %v want '1998'", modifiedArtist.Born)
	}

	if modifiedArtist.Genre != "alternative rock" {
		t.Errorf("artist genre not updated: got %v want 'alternative rock'", modifiedArtist.Genre)
	}
}
