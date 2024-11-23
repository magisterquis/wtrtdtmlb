# docker.mk
# Makefile things for Docker containers
# By J. Stuart McMurray
# Created 20241009
# Last Modified 20241123

.ifndef HAVE_DOCKER_MK
HAVE_DOCKER_MK=1

# The following variables ned to be defined before sourcing this makefile:
# DOCKER_NAME - Name of the docker container as well as its image
# LOCAL       - Temporary per-project storage directory
# GLOBAL      - Temporary global storage directory

# We'll need the below for what we do, but they should be defined elsewhere.
# We'd use .poison empty, but NetBSD's Make doesn't have it. :(
.for V in DOCKER_NAME LOCAL GLOBAL
.if empty($V)
ET=is_empty_$V
${ET}:: .NOTMAIN
	@echo Variable ${V} is empty. >&2; exit 1
.PHONY: ${ET}
.BEGIN:: ${ET}
ET=
.endif
.endfor

.include "curl.mk"
.include "debian.mk"
.include "rsyslog.mk"

# Variables which don't necessarily need to have been defined
DOCKER_NETWORK ?= bridge

# Files
DOCKER              = /usr/bin/docker
DOCKER_IMAGE_ID     = ${GLOBAL}/docker_image_id_${DOCKER_NAME}
DOCKER_CONTAINER_ID = ${LOCAL}/docker_container_id
DOCKER_NETWORK_ID   = ${GLOBAL}/docker_network_${DOCKER_NETWORK}

# Docker command bits
DOCKER_IS_RUNNING       = ( [ -n "$$(${GET_DOCKER_CONTAINER_ID})" ] )
DOCKER_IMAGE_EXISTS     = ( [ -e ${DOCKER} ] && [ -n "$$(${DOCKER} image ls\
				  --quiet ${DOCKER_NAME})" ] )
DOCKER_STOP             = ! ${DOCKER_IS_RUNNING} || {\
				${DOCKER} kill ${DOCKER_NAME};\
				while ${DOCKER_IS_RUNNING}; do sleep .1; done;\
			}; rm -f ${DOCKER_CONTAINER_ID}
GET_DOCKER_NETWORK_ID   = ${DOCKER} network ls --quiet\
				--filter name=${DOCKER_NETWORK:Q}
GET_DOCKER_CONTAINER_ID = ${DOCKER} ps --quiet --filter name=${DOCKER_NAME:Q}

# On startup, remove files indicating we have the image built and the container
# running if we don't have them yet.
docker_begin: .NOTMAIN
	@if ! ${DOCKER_IMAGE_EXISTS}; then\
		rm -vf ${DOCKER_IMAGE_ID} >&2;\
	fi
	@if ! [ -e ${DOCKER} ] || ! ${DOCKER_IS_RUNNING}; then\
		rm -f ${DOCKER_CONTAINER_ID};\
	fi
	@if ! [ -e ${DOCKER} ] || [ -z $$(${GET_DOCKER_NETWORK_ID}) ];\
	then\
		rm -f ${DOCKER_NETWORK_ID};\
	fi
.PHONY: docker_begin
.BEGIN:: docker_begin

# Install docker, for when we just need to test something and can't remember
# the path to the docker binary.
install_docker: ${DOCKER} .NOTMAIN
.PHONY: install_docker

${DOCKER}: .NOTMAIN /etc/apt/sources.list.d/docker.list
${DOCKER}: /etc/apt/keyrings/docker.asc
	${APTGET}\
		install\
		docker-ce\
		docker-ce-cli\
		containerd.io\
		docker-buildx-plugin\
		docker-compose-plugin
	${DOCKER} run --quiet --rm hello-world >/dev/null
	touch $@

/etc/apt/sources.list.d/docker.list: /etc/apt/keyrings/docker.asc .NOTMAIN
	{\
		echo -n "deb [arch=$$(dpkg --print-architecture) ";\
		echo -n "signed-by=/etc/apt/keyrings/docker.asc] ";\
		echo -n "https://download.docker.com/linux/debian ";\
		echo -n "$$(. /etc/os-release && echo "$$VERSION_CODENAME") ";\
		echo "stable";\
	} > $@
	${APTGET} update

/etc/apt/keyrings/docker.asc: ${CURL} .NOTMAIN
	install -m 0755 -d ${@:H}
	curl -fsSL -o $@ https://download.docker.com/linux/debian/gpg
	chmod a+r $@

# Build the image.  This should be used in a rule like
# ${DOCKER_IMAGE_ID}: <deps> build_image
# The deps will already include Dockerfile
build_image: ${DOCKER} Dockerfile .USE
	${DOCKER} build --quiet --tag ${DOCKER_NAME} . >$@.tmp
	mv $@.tmp $@

# Make sure a network exists.
${DOCKER_NETWORK_ID}: ${DOCKER} .NOTMAIN
	if [ -z "$$(${GET_DOCKER_NETWORK_ID})" ]; then\
		docker network create ${DOCKER_NETWORK:Q} >$@.tmp;\
	fi
	${GET_DOCKER_NETWORK_ID} >$@.tmp
	mv $@.tmp $@

# (Re)start the container.  This should be used in a rule like
# ${DOCKER_CONTAINER_ID}: <deps> start_container
# The optional variable DOCKER_RUN_ARGS may be used to pass additional
# arguments to docker run.
start_container: ${DOCKER} ${DOCKER_IMAGE_ID} ${DOCKER_NETWORK_ID} \
${SYSLOGD_PID} .USE
	${DOCKER_STOP}
	${DOCKER} run\
		${DOCKER_RUN_ARGS}\
		--detach\
		--init\
		--log-driver syslog\
		--log-opt tag=${DOCKER_NAME}\
		--name ${DOCKER_NAME}\
		--network ${DOCKER_NETWORK}\
		--quiet\
		--rm\
		${DOCKER_NAME} >$@.tmp
	mv $@.tmp $@

# Clean up things we may have made on clean
clean:: .NOTMAIN ${DOCKER}
	# Stop Container
	${DOCKER_STOP}
	# Remove the image
	if ${DOCKER_IMAGE_EXISTS}; then\
		${DOCKER} image rm ${DOCKER_NAME} >/dev/null;\
	fi
	rm -f ${DOCKER_IMAGE_ID}
	if [ -f ${DOCKER_NETWORK_ID} ] && [ "map[]" = "$$(\
		${DOCKER} network inspect\
			--format '{{.Containers}}'\
			${DOCKER_NETWORK:Q}\
	)" ]; then\
		docker network rm ${DOCKER_NETWORK} &&\
		rm ${DOCKER_NETWORK_ID};\
	fi

.endif # .ifndef HAVE_DOCKER_MK
