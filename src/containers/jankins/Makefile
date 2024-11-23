# Makefile
# Build Jankins
# By J. Stuart McMurray
# Created 20240323
# Last Modified 20241117

PROGNAME=jankins
PROGBIN=${LOCAL}/${PROGNAME}
DOCKER_NAME=${PROGNAME}
DOCKER_NETWORK=${PROGNAME}
LOCAL=local
JANKINS_PID=${PROGBIN}.pid
LOGFILE=${LOCAL}/jankins.log
LOCAL_KEY=${LOCAL}/id_ed25519
KNOWN_HOSTS=${LOCAL}/known_hosts
LEAKED_CREDS=${LOCAL}/ci.conf

.include "curl.mk"
.include "docker.mk"
.include "go.mk"
.include "m4.mk"
.include "ssh.mk"

SRCS != find . \
	-name '*.go' -o\
	-name '*.html' -o\
	-name '*.txtar' -o\
	-name 'default_*'

.BEGIN:: .NOTMAIN
	@mkdir -p ${LOCAL}
	@pgrep ${PROGNAME} >${JANKINS_PID} || rm -f ${JANKINS_PID}

# Make sure jankins is running.  curl is for tests eventually.
main: ${JANKINS_PID} ${CURL} ${LEAKED_CREDS}
.PHONY: main
${JANKINS_PID}: ${PROGBIN} ${DOCKER_IMAGE_ID} ${DOCKER_NETWORK_ID}
	while pgrep ${PROGNAME} >/dev/null; do pkill ${PROGNAME}; done
	PATH=$$PATH:${PROGBIN:H} ${PROGBIN:T} >>${LOGFILE} 2>&1 &
	pgrep ${PROGNAME} >$@.tmp
	mv $@.tmp $@

# Build the docker image
${DOCKER_IMAGE_ID}: build_image ${LOCAL_KEY} ${KNOWN_HOSTS} start_build.sh\
	.NOTMAIN

# Local SSH key, for cloning things
${LOCAL_KEY}: ${SSH_PRIVKEY}
	cp $> $@

# Known Hosts, for authenticating the codeserver
${KNOWN_HOSTS}: ../codeserver/local/docker_container_id ${DOCKER}\
	${DOCKER_NETWORK_ID} ${LOCAL_KEY}
	${DOCKER} run\
		--entrypoint sh\
		--network ${DOCKER_NETWORK}\
		--quiet\
		--rm\
		--volume ./${LOCAL_KEY}:/root/.ssh/${LOCAL_KEY:T}\
		alpine/git:latest\
		-c 'ssh\
				-N\
				-f\
				-o HashKnownHosts=no\
				-o StrictHostKeyChecking=accept-new\
				-q\
				git@codeserver </dev/null &&\
			cat /root/.ssh/known_hosts' >$@.tmp
	mv $@.tmp $@

# Creds which we'll find to start this whole thing
${LEAKED_CREDS}: ${M4} ../../include.m4/inc.m4 ci.conf.m4\
	default_username default_password
	${M4} -PEE ${>:M*.m4} >$@.tmp
	mv $@.tmp $@

# Build jankins
${PROGBIN}: ${SRCS} ${DOCKER_IMAGE_ID} go_build .NOTMAIN

# Stop the server and delete local things
clean::
	# Stop the process if it's going
	while pgrep ${PROGNAME} >/dev/null; do\
		pkill ${PROGNAME} >/dev/null;\
	done
	rm -f ${JANKINS_PID}
	rm -rf ${LOCAL}
