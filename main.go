package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Artist struct {
	ID    string   `json:"id"`    // id коллектива
	Name  string   `json:"name"`  // название группы
	Born  string   `json:"born"`  // год основания группы
	Genre string   `json:"genre"` // жанр
	Songs []string `json:"songs"` // популярные песни, это слайс строк, так как песен может быть несколько
}

// ...
func defaultArtists() map[string]Artist {
	return map[string]Artist{
		"1": {
			ID:    "1",
			Name:  "30 Seconds To Mars",
			Born:  "",
			Genre: "alternative",
			Songs: []string{
				"The Kill",
				"A Beautiful Lie",
				"Attack",
				"Live Like A Dream",
			},
		},
		"2": {
			ID:    "2",
			Name:  "Garbage",
			Born:  "1994",
			Genre: "alternative",
			Songs: []string{
				"Queer",
				"Shut Your Mouth",
				"Cup of Coffee",
				"Til the Day I Die",
			},
		},
	}
}

var artists = defaultArtists()
var artistsMu sync.RWMutex

func snapshotArtists() map[string]Artist {
	artistsMu.RLock()
	defer artistsMu.RUnlock()

	copyArtists := make(map[string]Artist, len(artists))
	for id, artist := range artists {
		copyArtists[id] = artist
	}

	return copyArtists
}

func getArtists(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	resp, err := json.Marshal(snapshotArtists())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}
func getArtist(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	id := chi.URLParam(r, "id")
	artistsMu.RLock()
	artist, ok := artists[id]
	artistsMu.RUnlock()

	if !ok {
		http.Error(w, "artist not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}
func postArtist(w http.ResponseWriter, r *http.Request) {
	var artist Artist
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &artist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	artistsMu.Lock()
	artists[artist.ID] = artist
	artistsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
func deleteArtists(w http.ResponseWriter, r *http.Request) {
	artistsMu.Lock()
	artists = map[string]Artist{}
	artistsMu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func deleteArtist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	artistsMu.Lock()
	defer artistsMu.Unlock()

	if _, ok := artists[id]; !ok {
		http.Error(w, "artist not found", http.StatusNotFound)
		return
	}

	delete(artists, id)
	w.WriteHeader(http.StatusNoContent)
}
func setupRouter() *chi.Mux {
	// создаём новый роутер
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинты для artists
	r.Get("/artists", getArtists)
	r.Post("/artists", postArtist)
	r.Delete("/artists", deleteArtists)
	r.Get("/artists/{id}", getArtist)
	r.Delete("/artists/{id}", deleteArtist)

	return r
}

func main() {
	r := setupRouter()

	// запускаем сервер
	if err := http.ListenAndServe(":9999", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
