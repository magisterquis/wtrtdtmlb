// Program Jankins is like an extremely janky Jenkins, with only the best parts.
package main

/*
 * jankins.go
 * Terrible Jenkins knockoff, just the fun bits.
 * By J. Stuart McMurray
 * Created 20241102
 * Last Modified 20241119
 */

import (
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/magisterquis/curlrevshell/lib/sstls"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// LEStagingURL is Let's Encrypt's staging environment's directory URL.
const LEStagingURL = "https://acme-staging-v02.api.letsencrypt.org/directory"

// Default default creds
var (
	//go:embed default_username
	DefaultUsername string
	//go:embed default_password
	DefaultPassword string
	//go:embed default_admin_username
	DefaultAdminUsername string
	//go:embed default_admin_password
	DefaultAdminPassword string
)

// defaultEnvVars are the environment variables we serve up, by default.
//
//go:embed default_env_vars
var defaultEnvVars string

var (
	// Jobsdir is where we store files for jobs we run.
	JobsDir = filepath.Join(os.TempDir(), "jankins_jobs")

	// ContainerCodePath is the path in the container where we make the
	// code we get in a txtar archive available.
	ContainerCodePath = "/code"

	// DefaultCommand is the default command to run in the container.
	DefaultCommand = "exec bmake"

	// ContainerImage is the name of the container image to use.
	ContainerImage = "jankins"

	// ContainerNetwork is the name of the container network to use.
	ContainerNetwork = "jankins"
)

func main() {
	var (
		leStaging = flag.Bool(
			"letsencrypt-staging",
			false,
			"Use Let's Encrypt's staging environment",
		)
		domain = flag.String(
			"domain",
			"",
			"Let's Encrypt `domain`, or empty for a "+
				"self-signed certificate",
		)
		jobsDir = flag.String(
			"jobs-dir",
			JobsDir,
			"Path to jobs `directory`",
		)
		containerCodePath = flag.String(
			"container-code",
			ContainerCodePath,
			"Path to code `directory` in container",
		)
		containerImage = flag.String(
			"image",
			ContainerImage,
			"Docker container `image` in which to run code",
		)
		containerNetwork = flag.String(
			"network",
			ContainerNetwork,
			"Docker `network` to which to attach containers",
		)
		defaultCommand = flag.String(
			"default-command",
			DefaultCommand,
			"Default build/test `command`",
		)
		httpsAddr = flag.String(
			"https-listen",
			"0.0.0.0:443",
			"HTTPS Listen `address`",
		)
		httpAddr = flag.String(
			"http-addr",
			"127.0.0.1:80",
			"HTTP Listen `address`",
		)
		username = flag.String(
			"username",
			strings.TrimSpace(DefaultUsername),
			"Basic auth `username`",
		)
		password = flag.String(
			"password",
			strings.TrimSpace(DefaultPassword),
			"Basic auth `password`",
		)
		adminUsername = flag.String(
			"admin-username",
			strings.TrimSpace(DefaultAdminUsername),
			"Basic auth admin `username`",
		)
		adminPassword = flag.String(
			"admin-password",
			strings.TrimSpace(DefaultAdminPassword),
			"Basic auth admin `password`",
		)
		allowCustomCommands = flag.Bool(
			"allow-custom-commands",
			false,
			"Allow custom build commands via txtar comments",
		)
		envVarsFile = flag.String(
			"env-vars",
			"",
			"Optional `file` containing KEY=value "+
				"environment variables for containers",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

Serves up a simple web interface which takes a txtar archive and runs it in
a container.

Use of this program to provision a certificate with Let's Encrypt constitutes
acceptance of Let's Encrypt's Terms of Service.

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Get environment variables for our container, if we have them. */
	var (
		envVars []string
		err     error
	)
	if "" != *envVarsFile {
		if envVars, err = ParseEnvVarsFile(*envVarsFile); nil != err {
			log.Fatalf(
				"Error parsing environment variables "+
					"from %s: %s",
				*envVarsFile,
				err,
			)
		}
	} else {
		if envVars, err = ParseEnvVars(strings.NewReader(
			defaultEnvVars,
		)); nil != err {
			log.Fatalf(
				"Error parsing built-in environment "+
					"variables: %s",
				err,
			)
		}
	}

	/* Set up HTTP */
	httpL, err := net.Listen("tcp", *httpAddr)
	if nil != err {
		log.Fatalf(
			"Error listening for HTTP connections on %s: %s",
			*httpAddr,
			err,
		)
	}
	log.Printf("Listening for HTTP connections on %s", httpL.Addr())

	/* Set up TLS. */
	var tlsL net.Listener
	if "" != *domain {
		/* Let's Encrypt cert-grabber. */
		mgr := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*domain),
		}
		/* Work out the whole staging thing. */
		if *leStaging {
			mgr.Client = &acme.Client{DirectoryURL: LEStagingURL}
		} else {
			mgr.Cache = autocert.DirCache(cacheDir())
		}
		/* Start the listener. */
		var err error
		if tlsL, err = tls.Listen(
			"tcp",
			*httpsAddr,
			mgr.TLSConfig(),
		); nil != err {
			log.Fatalf(
				"Error listening for HTTPS connections "+
					"on %s: %s",
				*httpsAddr,
				err,
			)
		}
		log.Printf(
			"Listening for HTTPS connections on %s",
			tlsL.Addr(),
		)
	} else { /* Self-signed cert. */
		sl, err := sstls.Listen(
			"tcp",
			*httpsAddr,
			"",
			0,
			sstls.DefaultCertFile(),
		)
		if nil != err {
			log.Fatalf(
				"Error listening for HTTPS connections with "+
					"a self-signed certificate on %s: %s",
				*httpsAddr,
				err,
			)
		}
		log.Printf(
			"Listening for HTTPS connections on %s with a "+
				"self-signed certificate with fingerprint %s",
			sl.Addr(),
			sl.Fingerprint,
		)
		tlsL = sl
	}

	/* Set up to run jobs. */
	if err := os.MkdirAll(*jobsDir, 0770); nil != err {
		log.Fatalf(
			"Failed to make jobs directory %s: %s",
			*jobsDir,
			err,
		)
	}
	var wg sync.WaitGroup
	handler := Handler{
		wg:            &wg,
		jobsDir:       *jobsDir,
		codeDir:       *containerCodePath,
		image:         *containerImage,
		network:       *containerNetwork,
		command:       *defaultCommand,
		username:      *username,
		password:      *password,
		adminUsername: *adminUsername,
		adminPassword: *adminPassword,
		allowCustom:   *allowCustomCommands,
		envVars:       envVars,
		logger:        log.Default(),
	}
	http.HandleFunc("POST /", handler.HandleJob)
	http.HandleFunc("/", handler.Handle)

	/* Serve requests. */
	ech := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		ech <- http.Serve(tlsL, nil)
	}()
	go func() {
		defer wg.Done()
		ech <- http.Serve(httpL, nil)
	}()
	log.Printf(
		"Fatal error, waiting for jobs to finish: %s",
		<-ech,
	)

	/* Wait for everything else to finish. */
	wg.Wait()
	log.Fatalf("All done.")
}
