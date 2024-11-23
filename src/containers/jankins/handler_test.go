package main

/*
 * handler_test.go
 * Tests for handler.go
 * By J. Stuart McMurray
 * Created 20241102
 * Last Modified 20241117
 */

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"embed"

	"golang.org/x/tools/txtar"
)

const (
	testUsername      = "test_username"
	testPassword      = "test_password"
	testAdminUsername = "test_admin_username"
	testAdminPassword = "test_admin_password"
)

//go:embed testdata
var testData embed.FS

func newTestHandler(t *testing.T) (Handler, *bytes.Buffer) {
	lb := new(bytes.Buffer)
	return Handler{
		wg:            new(sync.WaitGroup),
		jobsDir:       t.TempDir(),
		codeDir:       ContainerCodePath,
		image:         ContainerImage,
		network:       "bridge", /* We won't have jankins yet. */
		command:       DefaultCommand,
		username:      testUsername,
		password:      testPassword,
		adminUsername: testAdminUsername,
		adminPassword: testAdminPassword,
		logger:        log.New(lb, "", 0),
	}, lb
}

func TestHandlerUnauthed(t *testing.T) {
	/* Set up a mock server. */
	h, _ := newTestHandler(t)
	/* Make a request with r . */
	var r *http.Request
	run := func(t *testing.T) {
		/* Make an unauthed request. */
		w := httptest.NewRecorder()
		h.HandleJob(w, r)
		/* Make sure we got an Unauthorized back. */
		if want := http.StatusUnauthorized; want != w.Code {
			t.Errorf(
				"Incorrect status code: got: %d\nwant: %d",
				w.Code,
				want,
			)
		}
	}

	/* Without auth at all. */
	r = httptest.NewRequest(http.MethodPost, "/", nil)
	t.Run("no_auth", run)

	/* Wrong auth */
	r = httptest.NewRequest(http.MethodPost, "/", nil)
	r.SetBasicAuth("dummy", "dummy")
	t.Run("wrong_auth", run)
}

func TestHandlerHandleJob(t *testing.T) {
	txtarPath := "testdata/handler/HandleJob/normal/code.txtar"
	/* Set up a mock server. */
	h, lb := newTestHandler(t)
	/* Don't bother if we don't have docker. */
	if _, err := exec.LookPath("docker"); nil != err {
		t.Skipf("Unable to find docker: %s", err)
	}

	/* Set up a mock request. */
	ab, err := testData.ReadFile(txtarPath)
	if nil != err {
		t.Fatalf("Error opening %s in test data: %s", txtarPath, err)
	}
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		url.Values{archiveFormKey: []string{string(ab)}}.Encode(),
	))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(testUsername, testPassword)

	/* Run the job. */
	h.HandleJob(w, r)

	/* Make sure the log looks ok. */
	ls := strings.TrimSpace(lb.String())
	lls := strings.Split(ls, "\n")
	wantN := 3
	if got := len(lls); got != wantN {
		t.Fatalf(
			"Incorrect number of log lines:\n"+
				" got: %d\n"+
				"want: %d\n"+
				" log:\n%s",
			got,
			wantN,
			ls,
		)
	}
	/* First line should tell us where we unpacked it. */
	wantPrefix := "[192.0.2.1:1234] Unpacked archive to "
	var tdir string
	if !strings.HasPrefix(lls[0], wantPrefix) {
		t.Errorf(
			"Incorrect first log line (prefix)"+
				"\n got: %s\n"+
				"want: %s",
			lls[0],
			wantPrefix+"...",
		)
	} else {
		tdir = strings.TrimPrefix(lls[0], wantPrefix)
		/* Is it in the right place? */
		if !strings.HasPrefix(tdir, h.jobsDir) {
			t.Errorf(
				"Unpack directory in wrong place:\n"+
					" got: %s\n"+
					"want: %s",
				tdir,
				h.jobsDir,
			)
		}
	}
	/* Second line should be the build starting. */
	wantLine := "[192.0.2.1:1234] Starting build..."
	if got := lls[1]; got != wantLine {
		t.Errorf(
			"Incorrect second log line:\n got: %s\nwant: %s",
			got,
			wantLine,
		)
	}
	/* Third line should be how long it took. */
	wantPrefix = "[192.0.2.1:1234] Finished in "
	if !strings.HasPrefix(lls[2], wantPrefix) {
		t.Errorf(
			"Third log line has incorrect prefix:\n"+
				" got: %s\n"+
				"want: %s",
			lls[2],
			wantPrefix,
		)
	}

	/* Check the output now. */
	output := strings.TrimSpace(w.Body.String())
	ols := strings.Split(output, "\n")
	wantN = 8
	if got := len(ols); got != wantN {
		t.Fatalf(
			"Incorrect number of output lines:\n"+
				"   got: %d\n"+
				"  want: %d\n"+
				"output:\n%s",
			got,
			wantN,
			output,
		)
	}
	/* Output should end in the last log line, minus the address. */
	if _, ll, ok := strings.Cut(lls[2], " "); !ok {
		t.Errorf("Third log line had no space")
	} else if got := ols[len(ols)-1]; got != ll {
		t.Errorf(
			"Last output line incorrect\n got: %s\nwant: %s",
			got,
			ll,
		)
	}
	/* Should have got a pass, too. */
	wantLines := []string{"prove", "t/test.t .. ok"}
	if got := ols[2:4]; !slices.Equal(got, wantLines) {
		t.Errorf(
			"Output did not show test script result:\n"+
				" got: %s\n"+
				"want: %s",
			got,
			wantLine,
		)
	}
	wantLine = "Result: PASS"
	if got := ols[6]; got != wantLine {
		t.Errorf(
			"Output did not show tests passed:\n"+
				" got: %s\n"+
				"want: %s",
			got,
			wantLine,
		)
	}
}

func TestHandlerHandleJob_Comment(t *testing.T) {
	txtarPath := "testdata/handler/HandleJob/comment/code.txtar"
	/* Set up a mock server. */
	h, _ := newTestHandler(t)
	h.allowCustom = true
	/* Don't bother if we don't have docker. */
	if _, err := exec.LookPath("docker"); nil != err {
		t.Skipf("Unable to find docker: %s", err)
	}

	/* Set up a mock request. */
	ab, err := testData.ReadFile(txtarPath)
	if nil != err {
		t.Fatalf("Error opening %s in test data: %s", txtarPath, err)
	}
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		url.Values{archiveFormKey: []string{string(ab)}}.Encode(),
	))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(testUsername, testPassword)

	/* Run the job. */
	h.HandleJob(w, r)

	/* Should have five lines, of which we care about two.  Unfortunately
	sometimes the effect of -x comes after the command itself. */
	bls := strings.Split(strings.TrimSpace(w.Body.String()), "\n")
	wantN := 4
	if got := len(bls); wantN != got {
		t.Fatalf(
			"Incorrect number of response body lines:\n"+
				" got: %d\n"+
				"want: %d",
			got,
			wantN,
		)
	}
	wantL := "/code"
	if got := bls[2]; got != wantL {
		t.Errorf(
			"Command output incorrect:\n got: %s\nwant: %s",
			got,
			wantL,
		)
	}
}

func TestHandlerUnpackCode(t *testing.T) {
	txtarPath := "testdata/handler/unpackCode/normal/code.txtar"
	/* Set up a mock server. */
	h, _ := newTestHandler(t)
	/* Set up a mock request. */
	b := new(bytes.Buffer)
	l := func(f string, v ...any) {
		fmt.Fprintf(b, f, v...)
		b.WriteString("\n")
	}
	ab, err := testData.ReadFile(txtarPath)
	if nil != err {
		t.Fatalf("Error opening %s in test data: %s", txtarPath, err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		url.Values{archiveFormKey: []string{string(ab)}}.Encode(),
	))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(testUsername, testPassword)

	/* Unpack the data */
	td, _, err := h.unpackCode(l, w, r)
	if nil != err {
		t.Fatalf("Error: %s", err)
	}
	if "" == td {
		t.Fatalf("Got empty directory")
	}

	/* Parse the txtar so we can check unpacking. */
	ta := txtar.Parse(ab)

	/* Make sure our temp directory only has the one directory. */
	des, err := os.ReadDir(h.jobsDir)
	if nil != err {
		t.Fatalf(
			"Error reading temporary directory %s: %s",
			h.jobsDir,
			err,
		)
	}
	if 0 == len(des) {
		t.Fatalf("Temporary directory empty")
	} else if 1 != len(des) {
		ns := make([]string, len(des))
		for i, de := range des {
			ns[i] = de.Name()
		}
		t.Fatalf(
			"Temporary directory has too many children\n"+
				" got: %d\n"+
				"want: 1\n"+
				"Children:\n\t%s",
			len(des),
			strings.Join(ns, "\n\t"),
		)
	} else if got := filepath.Join(h.jobsDir, des[0].Name()); td != got {
		t.Fatalf(
			"Incorrect code directory\n got: %s\nwant: %s",
			des[0].Name(),
			td,
		)
	}

	/* Make sure it has what we expect. */
	seen := make(map[string]struct{})
	if err := filepath.WalkDir(td, func(
		path string,
		d fs.DirEntry,
		err error,
	) error {
		/* Silly... */
		if nil != err {
			return err
		}
		/* Only want regular files. */
		if d.IsDir() {
			return nil
		} else if ft := d.Type(); !ft.IsRegular() {
			t.Errorf("Irregular file %s, type %s", path, ft)
			return nil
		}
		/* Make sure this path should be here.  O(n**2).  Eh. */
		ind := slices.IndexFunc(ta.Files, func(f txtar.File) bool {
			return path == filepath.Join(td, f.Name)
		})
		if -1 == ind {
			t.Errorf("Extra file: %s", path)
			return nil
		}
		/* Note we've seen it. */
		seen[path] = struct{}{}
		/* Make sure the contents are as expected. */
		b, err := os.ReadFile(path)
		if nil != err {
			t.Errorf("Unable to read %s: %s", path, err)
			return nil
		}
		if !bytes.Equal(b, ta.Files[ind].Data) {
			t.Errorf("Found incorrect contents: %s", path)
		}
		return nil
	}); nil != err {
		t.Fatalf("Error checking extracted files: %s", err)
	}

	/* Make sure we got all the files. */
	for _, f := range ta.Files {
		n := filepath.Join(td, f.Name)
		if _, ok := seen[n]; !ok {
			t.Errorf("Missing file: %s", n)
			continue
		}
		delete(seen, n)
	}
	for f := range seen {
		t.Errorf("Unexpected extra file: %s", f)
	}
}
