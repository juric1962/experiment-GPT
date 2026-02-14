package main

import "sync"

type Mama struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Work string `json:"work"`
}

// MamaMap is the map type that stores Mama entries keyed by ID.
type MamaMap map[string]Mama

// add map struct for Mama
var mama = defaultMama()
var mamaMu sync.RWMutex

func defaultMama() MamaMap {
	return MamaMap{
		"1": {
			Name: "Mama",
			Age:  50,
			Work: "Housewife",
		},
	}
}

// Mama variable fourm Mama

// Mama function
func Maman() {
}

// defaultMama function called in main.go
// snapshotMama function
// snapshotMama returns a snapshot of the current Mama state as a map.
// It creates a copy of the Mama data to provide a consistent view of the state
// at a particular point in time. This is useful for preventing external modifications
// to the internal Mama state.
// insertMama adds or updates a Mama entry in the package-level map.
func insertMama(id string, m Mama) {
	mamaMu.Lock()
	defer mamaMu.Unlock()
	mama[id] = m
}

func snapshotMama() map[string]Mama {
	mamaMu.RLock()
	defer mamaMu.RUnlock()

	copyMama := make(map[string]Mama, len(mama))
	for id, mm := range mama {
		copyMama[id] = mm
	}
	return copyMama
}

// test snapshotMama
