package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokeapi"
	"github.com/dontsitdowncauseimovedyourchair/pokedex-go/internal/pokecache"
)

func newTestConfig() *config {
	return &config{
		caughtPokemon: make(map[string]pokeapi.Pokemon),
		cache:         pokecache.NewCache(time.Minute),
	}
}

// captureStdout runs fn with os.Stdout redirected and returns what it printed,
// so we can assert on user-facing output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	w.Close()
	os.Stdout = orig
	return <-done
}

// The command registry should be internally consistent: the key you type must
// match the command's own name, and every command needs a description and a
// runnable callback. This catches copy/paste slips in getCommands().
func TestGetCommandsConsistency(t *testing.T) {
	for key, cmd := range getCommands() {
		if cmd.name != key {
			t.Errorf("command registered under key %q has mismatched name %q", key, cmd.name)
		}
		if strings.TrimSpace(cmd.description) == "" {
			t.Errorf("command %q has an empty description", key)
		}
		if cmd.callback == nil {
			t.Errorf("command %q has a nil callback", key)
		}
	}
}

// Commands that require an argument must reject the wrong number of arguments
// with an error, rather than panicking or firing a request with junk input.
func TestCommandsValidateArgs(t *testing.T) {
	cfg := newTestConfig()
	cases := []struct {
		name string
		fn   func(*config, []string) error
		args []string
	}{
		{"explore with no args", commandExplore, []string{}},
		{"catch with no args", commandCatch, []string{}},
		{"catch with too many args", commandCatch, []string{"a", "b"}},
		{"inspect with no args", commandInspect, []string{}},
		{"inspect with too many args", commandInspect, []string{"a", "b"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.fn(cfg, c.args); err == nil {
				t.Errorf("expected a usage error, got nil")
			}
		})
	}
}

// Inspecting a pokemon you haven't caught should be handled gracefully: a
// friendly message and no error (never a crash or a network call).
func TestInspectUncaught(t *testing.T) {
	cfg := newTestConfig()

	out := captureStdout(t, func() {
		if err := commandInspect(cfg, []string{"pikachu"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(strings.ToLower(out), "caught") {
		t.Errorf("expected a 'not caught yet' message, got %q", out)
	}
}

// Lookups should be case-insensitive: a pokemon caught as "pikachu" must be
// inspectable via "PIKACHU".
func TestInspectIsCaseInsensitive(t *testing.T) {
	cfg := newTestConfig()
	cfg.caughtPokemon["pikachu"] = pokeapi.Pokemon{Name: "pikachu", Height: 4, Weight: 60}

	out := captureStdout(t, func() {
		if err := commandInspect(cfg, []string{"PIKACHU"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "pikachu") {
		t.Errorf("expected inspect output to include the pokemon name, got %q", out)
	}
}

// Catching a pokemon you already have should short-circuit (tell the user,
// no error, no network call).
func TestCatchAlreadyCaught(t *testing.T) {
	cfg := newTestConfig()
	cfg.caughtPokemon["pikachu"] = pokeapi.Pokemon{Name: "pikachu"}

	out := captureStdout(t, func() {
		if err := commandCatch(cfg, []string{"pikachu"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(strings.ToLower(out), "already caught") {
		t.Errorf("expected an 'already caught' message, got %q", out)
	}
}

// On the first page there is no previous page; mapb should say so and return
// without an error or a request.
func TestMapbOnFirstPage(t *testing.T) {
	cfg := newTestConfig()
	cfg.Previous = nil

	out := captureStdout(t, func() {
		if err := commandMapb(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(strings.ToLower(out), "first page") {
		t.Errorf("expected a 'first page' message, got %q", out)
	}
}
