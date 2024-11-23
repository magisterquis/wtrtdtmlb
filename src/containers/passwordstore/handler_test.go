package main

/*
 * handler.go
 * HTTP handler functions
 * By J. Stuart McMurray
 * Created 20241115
 * Last Modified 20241115
 */

import (
	"maps"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const testPasswordJSON = `{
	"foo": "bar",
	"tridge": "baaz",
	"quux": "kittens"
}`

var testPasswords = map[string]string{
	"foo":    "bar",
	"tridge": "baaz",
	"quux":   "kittens",
}

const (
	testUsername = "moose"
	testPassword = "squirrel"
)

func newTestHandler(t *testing.T) Handler {
	fn := filepath.Join(t.TempDir(), "p.json")
	if err := os.WriteFile(fn, []byte(testPasswordJSON), 0600); nil != err {
		t.Fatalf("Error writing password JSON to %s: %s", fn, err)
	}
	return Handler{
		passwordFile: fn,
		username:     testUsername,
		password:     testPassword,
	}
}

func TestHandlerReadPasswordFile(t *testing.T) {
	h := newTestHandler(t)
	got, err := h.readPasswordFile()
	if nil != err {
		t.Fatalf("Error reading password file: %s", err)
	}
	if !maps.Equal(got, testPasswords) {
		t.Errorf(
			"Password file read incorrect:\n"+
				" got: %+v\n"+
				"want: %+v\n",
			got,
			testPasswords,
		)
	}
}

func TestHandlerList(t *testing.T) {
	var want string
	for n := range slices.Values(slices.Sorted(maps.Keys(testPasswords))) {
		want += n + "\n"
	}
	req := httptest.NewRequest("", "/", nil)
	rr := httptest.NewRecorder()
	newTestHandler(t).List(rr, req)
	if got := rr.Body.String(); got != want {
		t.Errorf("Incorrect list\n got: %s\nwant: %s", got, want)
	}
}

func TestHandlerGet(t *testing.T) {
	h := newTestHandler(t)
	for n, want := range testPasswords {
		t.Run(n, func(t *testing.T) {
			req := httptest.NewRequest("", "/"+n, nil)
			req.SetBasicAuth(h.username, h.password)
			req.SetPathValue("name", n)
			rr := httptest.NewRecorder()
			h.Get(rr, req)
			if got := rr.Body.String(); got != want {
				t.Errorf(
					"Incorrect password:\n "+
						"got: %s\n"+
						"want: %s",
					got,
					want,
				)
			}
		})
	}
}
