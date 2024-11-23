package main

/*
 * handler.go
 * Serve up the main page and run a job.
 * By J. Stuart McMurray
 * Created 20241102
 * Last Modified 20241119
 */

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/magisterquis/mqd"
	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/txtar"
)

// archiveFormKey
const archiveFormKey = "archive"

// indexPage is what we serve up to most requests
//
//go:embed index.html
var indexPage []byte

// Handler handles HTTP requests.
type Handler struct {
	wg            *sync.WaitGroup
	jobsDir       string /* Outside container. */
	codeDir       string /* In container. */
	image         string
	network       string /* Container network. */
	command       string
	username      string
	password      string
	adminUsername string
	adminPassword string
	allowCustom   bool
	envVars       []string
	logger        *log.Logger
}

// Handle handles a non-job request.
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Write(indexPage)
}

// HandleJob starts a new job.
func (h Handler) HandleJob(w http.ResponseWriter, r *http.Request) {
	/* Auth also enables a timing attack.  Whee! */
	var (
		isAdmin         bool
		ruser, rpass, _ = r.BasicAuth()
	)
	if ruser == h.adminUsername && rpass == h.adminPassword {
		isAdmin = true
	} else if ruser != h.username || rpass != h.password {
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

	rc := http.NewResponseController(w)
	/* Logger to log to output as well as log. */
	l := func(f string, v ...any) {
		msg := fmt.Sprintf(f, v...)
		h.logger.Printf("[%s] %s", r.RemoteAddr, msg)
		fmt.Fprintf(w, "%s\n", msg)
		rc.Flush()
	}
	dn, comment, err := h.unpackCode(l, w, r)
	if nil != err {
		l("Error preparing code: %s", err)
		return
	}
	l("Unpacked archive to %s", dn)
	defer os.RemoveAll(dn)

	/* Work out what to run. */
	flags := []string{
		"--init",
		"--network", h.network,
		"--rm",
		"--volume", dn + ":" + h.codeDir,
		"--workdir", h.codeDir,
	}
	for _, v := range h.envVars {
		flags = append(flags, "--env", v)
	}
	if isAdmin {
		flags = append(flags, "--privileged")
		flags = append(flags, "--env", "JANKINS_IS_ADMIN=true")
	}
	argv := []string{
		"docker",
		"run",
	}
	argv = append(argv, flags...)
	argv = append(
		argv,
		h.image,
		"sh", "-c", /* buildCommand will be here. */
	)
	var buildCommand string
	if 0 != len(comment) && h.allowCustom {
		buildCommand = comment
	} else {
		buildCommand = h.command
	}
	argv = append(argv, buildCommand)

	/* Spawn a container to run the thing. */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, ectx := errgroup.WithContext(ctx)
	cmd := exec.CommandContext(ectx, argv[0], argv[1:]...)
	cmd.Stdout = w
	cmd.Stderr = w
	l("Starting build...")
	start := time.Now()
	/* Run the command. */
	eg.Go(func() error { defer cancel(); return cmd.Run() })
	/* Flush output every so often. */
	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				rc.Flush()
			case <-ectx.Done():
				return nil
			}
		}
	})
	err = eg.Wait()
	d := time.Since(start).Round(time.Millisecond)
	if nil != err {
		l("Error building after %s: %s", d, err)
		return
	}
	l("Finished in %s", d)
}

// unpackCode unpacks the code sent in the archive in the body of the request.
// The directory with the code is returned.  l is a log function.
func (h Handler) unpackCode(
	l func(string, ...any),
	w http.ResponseWriter,
	r *http.Request,
) (dir, comment string, err error) {
	/* Get the archive sent to us. */
	v := r.FormValue(archiveFormKey)
	if "" == v {
		return "", "", fmt.Errorf("empty request")
	}
	v = strings.ReplaceAll(v, "\r\n", "\n")

	mqd.Logf("Archive:\n%s", v)

	/* Unpack as a txtar archive and clean paths. */
	ta := txtar.Parse([]byte(v))
	for i, f := range ta.Files {
		ta.Files[i].Name = filepath.Clean(f.Name)
	}
	afs, err := txtar.FS(ta)
	if nil != err {
		return "", "", fmt.Errorf("parsing archive: %s", err)
	}
	/* Unpack to a temporary directory. */
	dn, err := os.MkdirTemp(
		h.jobsDir,
		"files-"+strings.ReplaceAll(r.RemoteAddr, ":", "_")+"-",
	)
	if nil != err {
		return "", "", fmt.Errorf(
			"creating temporary directory for files: %s",
			err,
		)
	}
	if err := os.CopyFS(dn, afs); nil != err {
		err := fmt.Errorf("unpacking archive: %s", err)
		if err := os.RemoveAll(dn); nil != err {
			l("Error removing temporary directory %s: %s", dn, err)
		}
		return "", "", err
	}

	return dn, string(ta.Comment), nil
}
