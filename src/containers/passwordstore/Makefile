# Makefile
# Build the Password store
# By J. Stuart McMurray
# Created 20241115
# Last Modified 20241116

PROGNAME=passwordstore
LOCAL=local
DOCKER_NAME=${PROGNAME}
DOCKER_NETWORK=jankins
PROGBIN=${LOCAL}/${PROGNAME}
PASSWORDS=${LOCAL}/passwords.json

# This is like .BEGIN, but even earlier in the flow, for -q
.BEGIN::
	@mkdir -p ${LOCAL}

.include "docker.mk"
.include "go.mk"
.include "m4.mk"

# curl is for tests eventually
main: ${DOCKER_CONTAINER_ID}
.PHONY: main

# Start the container going
${DOCKER_CONTAINER_ID}: start_container

# Build the docker image
${DOCKER_IMAGE_ID}: ${PROGBIN} ${PASSWORDS} 007.procpidroot.flag build_image

# Build the binary
${PROGBIN}: 008.bakedin.flag ${LOCAL}/009.mem.flag.hex *.go default_* go_build

${PASSWORDS}: ${M4} ../../include.m4/inc.m4 passwords.json.m4\
	../jankins/default_admin_username ../jankins/default_admin_password\
	${LOCAL}/lsys_id_ed25519
	${M4} -PEE ${>:M*.m4} >$@.tmp
	mv $@.tmp $@

${LOCAL}/009.mem.flag.hex: 009.mem.flag
	perl -pe '$$_=unpack("H*", $$_)' $> > $@.tmp
	mv $@.tmp $@

${LOCAL}/lsys_id_ed25519:
	ssh-keygen\
		-q\
		-t ed25519\
		-N ''\
		-C "root@lsys"\
		-f $@

# Clean ALL the things
clean::
	rm -rf ${LOCAL}
