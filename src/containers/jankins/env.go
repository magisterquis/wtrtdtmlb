package main

/*
 * env.go
 * Environment variables in a file
 * By J. Stuart McMurray
 * Created 20241115
 * Last Modified 20241115
 */

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// ParseEnvVarsFile parses the KEY=value pairs from fn and returns them as a
// sorted list.
func ParseEnvVarsFile(fn string) ([]string, error) {
	/* Open the file for line-by-line reading. */
	f, err := os.Open(fn)
	if nil != err {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return ParseEnvVars(f)
}

// ParseEnvVars parses the KEY=value pairs read from r and returns them as
// a sorted list.
func ParseEnvVars(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)

	/* Read lines and parse into variables. */
	vars := make(map[string]string)
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		/* Ignore blanks and comments. */
		if "" == l || strings.HasPrefix(l, "#") {
			continue
		}
		/* Grab the key and value, and skip lines which don't have
		both. */
		k, v, ok := strings.Cut(l, "=")
		if !ok {
			continue
		}
		/* Save this pair, deduping. */
		vars[k] = v
	}
	if err := scanner.Err(); nil != err {
		return nil, fmt.Errorf("reading line: %w", err)
	}

	/* Return a list of pairs. */
	ret := make([]string, 0, len(vars))
	for k, v := range vars {
		ret = append(ret, k+"="+v)
	}
	slices.Sort(ret)
	return ret, nil
}
