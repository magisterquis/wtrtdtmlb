# go.mk
# Makefile things for Go projects
# By J. Stuart McMurray
# Created 20241009
# Last Modified 20241102

.ifndef HAVE_GO_MK
HAVE_GO_MK=1

.include "curl.mk"

# Variables.  One can set GOLDFLAGS to pass things to -ldflags, like -X...
GO           = ${MGOPATH}/go/bin/go
MGOPATH      = ${HOME}/go
STATICCHECK  = ${MGOPATH}/bin/staticcheck
GOBUILDFLAGS = -trimpath -ldflags "-w -s ${GOLDFLAGS}"

# Go compiler
${GO}: ${CURL} .NOTMAIN
	mkdir -p $@
	${CURL} -sfL https://dl.google.com/go/$$(\
		${CURL} -sfL 'https://go.dev/VERSION?m=text' | head -n 1\
	).linux-${MACHINE_ARCH}.tar.gz  | tar -C "${MGOPATH}" -xzf -
	echo "export PATH=$$PATH:$$(\
		$@ env GOROOT\
	)/bin:$$(\
		$@ env GOPATH\
	)/bin" >> ${HOME}/.bashrc
	touch $@

# Static code analyzer
${STATICCHECK}: ${GO} .NOTMAIN
	${GO} install honnef.co/go/tools/cmd/staticcheck@latest

# Make sure all the Go code is solid
go_test: ${GO} ${STATICCHECK} .USE
	${GO} test ${GOBUILDFLAGS} -timeout 10s ./...
	${GO} vet  ${GOBUILDFLAGS} ./...
.if ${PATH:S/${MGOPATH}//} == ${PATH}
	. ~/.bashrc; ${STATICCHECK} ./...
.else
	${STATICCHECK} ./...
.endif

# Build a thing with Go.
# Should be used in a rule like 
# ${THING}: <deps> go_build
go_build: go_test .USE
	${GO} build ${GOBUILDFLAGS} -o $@

.endif # .ifndef HAVE_GO_MK
