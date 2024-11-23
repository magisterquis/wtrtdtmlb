# Makefile
# Get this whole thing running
# By J. Stuart McMurray
# Created 20241102
# Last Modified 20241123

HOSTDISKFLAG=/usr/include/005.hostdisk.flag

# Make sure we have a global config directory, for things shared between source
# directories, like docker networks.
GLOBAL=${.CURDIR}/global
.BEGIN::
	@mkdir -p ${GLOBAL}

PROVE     = SHMORE=$$(pwd)/src/include.t/shmore.subr\
		    prove\
		    -Isrc/include.t\
		    --recurse
SUBMAKES != find src -mindepth 2 -name Makefile -type f
.MAKEFLAGS: -I ${.CURDIR}/src/include.mk GLOBAL=${GLOBAL:Q}

.include "ed.mk"
	
# We can't actually build things in parallel, because jankins relies on the
# other two containers.  The containers themeslevs can be built in parallel,
# though.
.NOTPARALLEL:

# Install a couple of handy utilities and make sure our thing works.
all: ${ED} ${STRINGS} test ${HOSTDISKFLAG}

# Test ALL the tests!
test: recursive_makes ${HOSTDISKFLAG}
	${PROVE}

# Make ALL the makes!
recursive_makes::
.PHONY: recursive_makes

# Clean ALL the cleans!
clean:: recursive_cleans .NOTMAIN
	rm -rf ${GLOBAL}
.PHONY: clean

# Submake ALL the submakes!
.for S in ${SUBMAKES:H}
clean_${S:T}: .NOTMAIN
	${MAKE} -C $S clean
.PHONY: clean_${S:T}
recursive_cleans:: clean_${S:T}

# Subclean ALL the subcleans!
${S:T}:: .NOTMAIN
	${MAKE} -C $S
.PHONY: ${S:T}
recursive_makes:: ${S:T}

# Subtest ALL the subtests!
test_${S:T}: ${S:T} .NOTMAIN
	${PROVE} $S

.endfor

# Jankins containers need the the other two running
jankins:: codeserver passwordstore

# Leave flags outside the containers
${HOSTDISKFLAG}: src/flags/005.hostdisk.flag
	cp $> $@
