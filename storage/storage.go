package storage

import (
	"errors"
	"sync"
)

// URLStore defines the interface for URL storage
type URLStore interface {
	// Save stores a URL and returns a unique short key
	Save(url string) (string, error)
	
	// Load retrieves the original URL for a given short key
	Load(shortKey string) (string, error)
}

// InMemoryStore implements URLStore using an in-memory map
type InMemoryStore struct {
	urls  map[string]string
	mutex sync.RWMutex
	count int
}

// NewInMemoryStore creates a new in-memory URL store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		urls: make(map[string]string),
	}
}

// Save stores a URL and returns a unique short key
func (s *InMemoryStore) Save(url string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Generate a simple key based on the count
	s.count++
	key := generateKey(s.count)
	
	s.urls[key] = url
	return key, nil
}

// Load retrieves the original URL for a given short key
func (s *InMemoryStore) Load(shortKey string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	url, exists := s.urls[shortKey]
	if !exists {
		return "", errors.New("URL not found")
	}
	
	return url, nil
}

// generateKey creates a short key from an integer using base62 encoding
func generateKey(n int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n == 0 {
		return string(charset[0])
	}
	
	result := ""
	base := len(charset)
	
	for n > 0 {
		result = string(charset[n%base]) + result
		n /= base
	}
	
	return result
}