package main

/*
 * env_test.go
 * Test for env.go
 * By J. Stuart McMurray
 * Created 20241115
 * Last Modified 20241115
 */

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestParseEnvVarsFile(t *testing.T) {
	have := `
FOO=bar
FOO=tridge
abc

# FOO=baaz
KITTENS=moose
FOO=quux
`
	want := []string{
		"FOO=quux",
		"KITTENS=moose",
	}
	fn := filepath.Join(t.TempDir(), "have")
	if err := os.WriteFile(fn, []byte(have), 0600); nil != err {
		t.Fatalf("Error writing have data to %s: %s", fn, err)
	}
	got, err := ParseEnvVarsFile(fn)
	if nil != err {
		t.Fatalf("Error parsing %s: %s", fn, err)
	}
	if !slices.Equal(got, want) {
		t.Errorf(
			"Did not get expected variables:\n"+
				"got:\n%s\n"+
				"want:\n%s\n",
			strings.Join(got, "\n"),
			strings.Join(want, "\n"),
		)
	}
}
