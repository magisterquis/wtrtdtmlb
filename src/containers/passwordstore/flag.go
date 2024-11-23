package main

/*
 * flag.go
 * Flag only in process memory
 * By J. Stuart McMurray
 * Created 20241123
 * Last Modified 20241123
 */

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//go:embed local/009.mem.flag.hex
var memFlag string

// FlagChild spawns a perl process which reads a flag and dies when stdin
// dies.
func FlagChild() {
	/* Decode the flag. */
	dec, err := hex.DecodeString(strings.TrimSpace(memFlag))
	if nil != err {
		log.Fatalf("Decoding memory flag: %s", err)
	}
	/* Spawn perl and send it. */
	perl := exec.Command("perl", "-e", "while(<>){}")
	si, err := perl.StdinPipe()
	if nil != err {
		log.Fatalf("Getting perl's stdin: %s", err)
	}
	perl.Stdout = os.Stdout
	perl.Stderr = os.Stderr
	go func() {
		fmt.Fprintf(si, "%s\n", dec)
	}()
	if err := perl.Run(); nil != err {
		log.Fatalf("Perl died: %s", err)
	}
}
