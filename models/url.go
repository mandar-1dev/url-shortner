package models

import "sync"

// URL represents a single shortened URL entry.
// The `json` tags control how this struct is serialized in API responses.
type URL struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	Clicks      int    `json:"clicks"`
}

// URLStore is an in-memory, thread-safe storage for URLs.
// We use a map for O(1) lookups by short code, and a mutex
// to prevent race conditions when multiple requests hit the
// API at the same time (Gin handles requests concurrently).
type URLStore struct {
	mu   sync.RWMutex
	urls map[string]*URL
}

// NewURLStore creates and returns a new, empty URLStore.
func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]*URL),
	}
}

// Save creates a new URL entry and stores it under the given short code.
func (s *URLStore) Save(code, originalURL string) *URL {
	s.mu.Lock()
	defer s.mu.Unlock()

	url := &URL{
		OriginalURL: originalURL,
		ShortCode:   code,
		Clicks:      0,
	}
	s.urls[code] = url
	return url
}

// Get retrieves a URL by its short code.
// The second return value tells the caller whether it was found.
func (s *URLStore) Get(code string) (*URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[code]
	return url, ok
}

// GetAll returns a slice containing every stored URL.
func (s *URLStore) GetAll() []*URL {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]*URL, 0, len(s.urls))
	for _, u := range s.urls {
		all = append(all, u)
	}
	return all
}

// Delete removes a URL entry by its short code.
// Returns false if the code didn't exist.
func (s *URLStore) Delete(code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.urls[code]; !ok {
		return false
	}
	delete(s.urls, code)
	return true
}

// IncrementClicks increases the click counter for a given short code.
// Called every time someone visits a short URL.
func (s *URLStore) IncrementClicks(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if u, ok := s.urls[code]; ok {
		u.Clicks++
	}
}

// Exists checks whether a short code is already in use.
// Used to avoid collisions when generating new codes.
func (s *URLStore) Exists(code string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.urls[code]
	return ok
}
