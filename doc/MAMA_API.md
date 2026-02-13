# Mama API Documentation

## Overview
The Mama API provides endpoints to manage mama data with thread-safe operations using read-write mutexes.

## Data Model

### Mama Struct
```go
type Mama struct {
	Name string `json:"name"` // Mother's name
	Age  int    `json:"age"`  // Mother's age
	Work string `json:"work"` // Mother's occupation
}
```

## Functions

### defaultMama()
Returns the default mama data as a map.

**Returns:** `map[string]Mama`

**Example:**
```go
mama := defaultMama()
// Output: map[string]Mama{"1": {Name: "Mama", Age: 50, Work: "Housewife"}}
```

### snapshotMama()
Returns a thread-safe copy of the mama data. Similar to `snapshotArtists()` but for Mama entities.

**Returns:** `map[string]Mama`

**Thread Safety:** Uses read lock (RWMutex) to prevent concurrent access issues.

**Example:**
```go
snapshot := snapshotMama()
// Returns a safe copy of the current mama state
```

## Proposed Endpoints

### GET /mama
Retrieves all mama data.

**Response:**
```json
{
  "1": {
    "name": "Mama",
    "age": 50,
    "work": "Housewife"
  }
}
```

**Status Code:** 200 OK

---

### GET /mama/{id}
Retrieves a specific mama by ID.

**Parameters:**
- `id` (path): Mama ID

**Response:**
```json
{
  "name": "Mama",
  "age": 50,
  "work": "Housewife"
}
```

**Status Code:** 200 OK
**Status Code:** 404 Not Found (if mama not found)

---

### POST /mama
Creates or updates a mama entry.

**Request Body:**
```json
{
  "name": "Mama",
  "age": 50,
  "work": "Housewife"
}
```

**Status Code:** 201 Created

---

### DELETE /mama
Deletes all mama entries.

**Status Code:** 204 No Content

---

### DELETE /mama/{id}
Deletes a specific mama by ID.

**Parameters:**
- `id` (path): Mama ID

**Status Code:** 204 No Content
**Status Code:** 404 Not Found (if mama not found)

## Testing

Run tests with cache disabled:
```bash
go test -count=1 ./... -v
```

**Test Functions:**
- `TestDefaultMama()` - Validates default mama data
- `TestSnapshotMama()` - Tests snapshot functionality
- `TestMamaStruct()` - Tests Mama struct creation
