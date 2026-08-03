package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokecache"
)

// A minimal shape so these tests don't depend on the real API's schema.
type testResource struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// A successful response should be fetched and decoded into T.
func TestGetResourceDecodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"name":"pikachu","value":42}`)
	}))
	defer server.Close()

	got, err := GetResource[testResource](server.URL, pokecache.NewCache(time.Minute))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "pikachu" || got.Value != 42 {
		t.Errorf("decoded %+v, want {Name:pikachu Value:42}", got)
	}
}

// A second request for the same URL should be served from cache, not refetched.
func TestGetResourceCachesResponse(t *testing.T) {
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		fmt.Fprint(w, `{"name":"pikachu","value":42}`)
	}))
	defer server.Close()

	cache := pokecache.NewCache(time.Minute)

	first, err := GetResource[testResource](server.URL, cache)
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}
	second, err := GetResource[testResource](server.URL, cache)
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}

	if hits != 1 {
		t.Errorf("expected the server to be hit once (second served from cache), got %d hits", hits)
	}
	if first != second {
		t.Errorf("cached result %+v differs from the original %+v", second, first)
	}
}

// A 404 should surface as an error and a good client must NOT cache the
// failure, so a retry actually hits the server again.
func TestGetResource404IsErrorAndNotCached(t *testing.T) {
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cache := pokecache.NewCache(time.Minute)

	if _, err := GetResource[testResource](server.URL, cache); err == nil {
		t.Errorf("expected an error for a 404 response, got nil")
	}
	if _, err := GetResource[testResource](server.URL, cache); err == nil {
		t.Errorf("expected an error on retry, got nil")
	}
	if hits != 2 {
		t.Errorf("expected 2 server hits (errors must not be cached), got %d", hits)
	}
}

// Any non-200 (or 200 osmething) status should be reported as an error.
func TestGetResourceServerErrorIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	if _, err := GetResource[testResource](server.URL, pokecache.NewCache(time.Minute)); err == nil {
		t.Errorf("expected an error for a 500 response, got nil")
	}
}

// Malformed JSON must be reported, not silently swallowed into a zero value.
func TestGetResourceInvalidJSONIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `this is not json`)
	}))
	defer server.Close()

	if _, err := GetResource[testResource](server.URL, pokecache.NewCache(time.Minute)); err == nil {
		t.Errorf("expected an error for malformed JSON, got nil")
	}
}
