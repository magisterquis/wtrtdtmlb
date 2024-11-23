# Makefile
# Build the Codeserver
# By J. Stuart McMurray
# Created 20241115
# Last Modified 20241123

DOCKER_NAME=codeserver
DOCKER_NETWORK=jankins
LOCAL=local
CODESERVERBIN=${LOCAL}/codeserver
REPOSDIR=${LOCAL}/repos
INFRAREPO=central_iac.git
REPOS=${REPOSDIR}/${INFRAREPO}
GET_CREDS=${LOCAL}/get_creds.subr

# This is like .BEGIN, but even earlier in the flow, for -q
.BEGIN::
	@mkdir -p ${LOCAL} ${REPOSDIR}

.include "docker.mk"
.include "git.mk"
.include "go.mk"
.include "m4.mk"
.include "ssh.mk"

# Set it running
main: ${DOCKER_CONTAINER_ID}
.PHONY: main

# Start the container going
${DOCKER_CONTAINER_ID}: start_container

# Cheese a handful of repos
.for R in curlrevshell dtffmacac mqd httpd_botnet_controller wtrtdtmlb
REPOS += ${REPOSDIR}/$R.git
${REPOSDIR}/$R.git: ${GIT}
	if [ -d $@.tmp ]; then rm -r $@.tmp; fi
	git clone --bare --quiet https://github.com/magisterquis/${@F} $@.tmp
	mv $@.tmp $@
.endfor

# Build the docker image
${DOCKER_IMAGE_ID}: ${CODESERVERBIN} start.sh ${REPOS}\
	${LOCAL}/authorized_keys build_image

# Allow our SSH key in
${LOCAL}/authorized_keys: ${SSH_PUBKEY}
	cp $> $@

# Build the HTTP server
${CODESERVERBIN}: *.go go_build

# And add our own repo, or at least another with a file with creds plus a
# flag.
${REPOSDIR}/${INFRAREPO}: ${GIT} ${GET_CREDS} 004.repo.flag
	if [ -d $@.full ]; then rm -r $@.full; fi
	git clone --quiet https://github.com/magisterquis/wtrtdtmlb.git $@.full
	find $@.full -name Makefile -delete
	cp ${>:N*/git} $@.full
	(cd $@.full &&\
		git add . &&\
		GIT_COMMITTER_NAME=infrabot\
		GIT_COMMITTER_EMAIL=infrabot\
			git commit --quiet -am "Initial commit")
	rm -rf $@ && mv $@.full/.git $@ && rm -r $@.full

${GET_CREDS}: ${M4} ../../include.m4/inc.m4 get_creds.subr.m4\
	../passwordstore/default_username ../passwordstore/default_password
	${M4} -PEE ${>:M*.m4} >$@.tmp
	mv $@.tmp $@

clean::
	rm -rf ${LOCAL}
.PHONY: clean
