// Program passwordstore - A simple read-only password store
package main

/*
 * passwordstore.go
 * A simple read-only password store
 * By J. Stuart McMurray
 * Created 20241115
 * Last Modified 20241123
 */

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
)

var (
	//go:embed default_username
	defaultUsername string
	//go:embed default_password
	defaultPassword string
	//go:embed 008.bakedin.flag
	bakedInFlag string
)

func main() {
	/* Command-line flags. */
	var (
		lAddr = flag.String(
			"listen",
			"0.0.0.0:80",
			"Listen `address`",
		)
		baUser = flag.String(
			"user",
			strings.TrimSpace(defaultUsername),
			"Basic auth `username`",
		)
		baPass = flag.String(
			"pass",
			strings.TrimSpace(defaultPassword),
			"Basic auth `password`",
		)
		printBakedInFlag = flag.Uint64(
			"print-flag",
			0,
			"Print the flag if the magic `number` is correct",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options] passwords.json

A simple read-only password store.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Print the flag if someone guesses a random number. */
	if 0 != *printBakedInFlag && *printBakedInFlag == rand.Uint64() {
		log.Printf("Flag: %s", bakedInFlag)
		return
	}

	/* Spawn scrapable perl. */
	go FlagChild()

	/* Make sure we have a filename. */
	if 0 == flag.NArg() {
		log.Fatalf("Need a password file")
	}

	/* Set up handlers. */
	handler := Handler{
		passwordFile: flag.Arg(0),
		username:     *baUser,
		password:     *baPass,
	}
	/* Serve HTTP. */
	http.HandleFunc("/", handler.List)
	http.HandleFunc("/{name}", handler.Get)
	log.Fatalf("Error: %s", http.ListenAndServe(*lAddr, nil))
}
