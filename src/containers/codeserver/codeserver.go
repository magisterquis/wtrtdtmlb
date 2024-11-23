// Program codeserver - A read-only server for git repos
package main

/*
 * codeserver.go
 * A read-only server for git repos
 * By J. Stuart McMurray
 * Created 20241116
 * Last Modified 20241116
 */

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	/* Command-line flags. */
	var (
		lAddr = flag.String(
			"listen",
			"0.0.0.0:80",
			"Listen `address`",
		)
		serveRoot = flag.String(
			"dir",
			"/git",
			"Root `directory` to serve",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

A read-only server for git repos

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Set up the read-only server. */
	fs := http.FileServer(http.Dir(*serveRoot))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s", r.RemoteAddr, r.URL)
		fs.ServeHTTP(w, r)
	})

	/* Listen. */
	l, err := net.Listen("tcp", *lAddr)
	if nil != err {
		log.Fatalf("Error listening on %s: %s", *lAddr, err)
	}
	log.Printf("Listening on %s", *lAddr)

	/* Serve */
	log.Fatalf("Error: %s", http.Serve(l, nil))
}
