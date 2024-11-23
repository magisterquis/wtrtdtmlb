package main

/*
 * handler.go
 * HTTP handler functions
 * By J. Stuart McMurray
 * Created 20241115
 * Last Modified 20241115
 */

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"slices"
)

// Handler contains HTTP handler functions
type Handler struct {
	passwordFile string
	username     string
	password     string
}

// List returns a list of the passwords we know.
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	/* Get the passwords. */
	ps, err := h.readPasswordFile()
	if nil != err {
		http.Error(
			w,
			fmt.Sprintf(
				"%s\n%s\n",
				http.StatusText(http.StatusInternalServerError),
				err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}
	for n := range slices.Values(slices.Sorted(maps.Keys(ps))) {
		fmt.Fprintf(w, "%s\n", n)
	}
}

// Get returns the requested secret.
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	/* Make sure auth is correct. */
	if u, p, ok := r.BasicAuth(); !ok ||
		u != h.username || p != h.password {
		w.Header().Set(
			"WWW-Authenticate",
			`Basic realm="restricted"`,
		)
		http.Error(
			w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized,
		)
		return
	}
	/* Work out which password we want. */
	n := r.PathValue("name")
	if "" == n {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}
	/* Try to get the password. */
	ps, err := h.readPasswordFile()
	if nil != err {
		http.Error(
			w,
			fmt.Sprintf(
				"%s\n%s\n",
				http.StatusText(http.StatusInternalServerError),
				err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}
	p, ok := ps[n]
	if !ok {
		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
		return
	}
	fmt.Fprintf(w, "%s", p)
}

// readPasswordFile reads the JSON file containing passwords.
func (h Handler) readPasswordFile() (map[string]string, error) {
	/* Open the password file. */
	f, err := os.Open(h.passwordFile)
	if nil != err {
		return nil, fmt.Errorf("opening password file: %w", err)
	}
	defer f.Close()

	/* Un-JSON it. */
	var ret map[string]string
	if err := json.NewDecoder(f).Decode(&ret); nil != err {
		return nil, fmt.Errorf("un-JSONing password file: %w", err)
	}

	return ret, nil
}
