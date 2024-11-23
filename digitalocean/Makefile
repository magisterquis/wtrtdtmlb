# Makefile
# Manage DigitalOcean Droplets
# By J. Stuart McMurray
# Created 20241007
# Last Modified 20241120

# Things we use for droplets
PROJECTNAME    = wtrtdtmlb
CODESYNC       = ${LOCAL}/lastcodesync
DROPLETADDR    = ${LOCAL}/serveraddr
DROPLETINFO    = ${LOCAL}/droplet.json
DROPLETNAME   ?= ${PROJECTNAME}
DROPLETTAG    ?= ${PROJECTNAME}
LOCAL          = ./local
PROJECTID      = 0d1fc7b4-6d6c-4947-946e-3aec42f6cde3
REMOTESRCDIR   = ${PROJECTNAME}
RUNCMD         = ./run.sh
SSHARGS        = ssh -i ${SSHKEY} -oStrictHostKeyChecking=accept-new -p 2
SSHCMD         = ${LOCAL}/sshcmd
SSHKEYNAME    ?= ${PROJECTNAME}
SSHKEY        ?= ${HOME}/.ssh/id_ed25519_${SSHKEYNAME}
SSHKEYID       = ${LOCAL}/sshkeyid
SSHPORT        = 2
TOSYNC        != find .. \! -name '*.swp' \! -path '../digitalocean/*'
USERDATA      ?= ./user-data-debian.sh
VPCUUID        = e5e6c064-06d6-44e2-af9f-086197e3aa08

DROPLET      !!= cat ${DROPLETADDR}
SSH          !!= cat ${SSHCMD}
.poison empty SSH
.poison empty DROPLET

.BEGIN:
.if ! exists(${LOCAL})
	mkdir -p ${LOCAL}
.endif

# Run the command in run.sh on the droplet, which will be make test if it
# doesn't exist
run: ${RUNCMD} ${SSHCMD} ${CODESYNC}
	{ echo 'mkdir -p ${REMOTESRCDIR} && cd ${REMOTESRCDIR}';\
		tee /dev/stderr <${RUNCMD}; } | ${SSH} -T
.PHONY: run

# Just SSH to the Droplet
ssh: ${SSHCMD}
	${SSH}
.PHONY: ssh

# Rsync up this repo.  sync syncs regardless of file changes.
force_sync: _openrsync 
sync: ${CODESYNC}
${CODESYNC}: ${TOSYNC} _openrsync
_openrsync: ${SSHCMD} ${DROPLETADDR} .USE
	openrsync --rsh "${SSHARGS}"\
		--archive\
		--delete\
		--exclude '*.swp'\
		--exclude '/.git*'\
		--exclude */local\
		--exclude /digitalocean\
		--exclude /global\
		--verbose\
		../ ${DROPLET}:${REMOTESRCDIR} 2>&1 |\
		tee ${CODESYNC} |\
		awk '!/^openrsync: warning: .*: skipping excluded file/'
.PHONY: force_sync sync

# Create a file with a command to run, expected to usually be hand-edited
${RUNCMD}:
	echo "bmake" > ${RUNCMD}

# Make sure we have an SSH key
${SSHKEY}:
	@echo "${SSHKEY} not found.  Please create with" >&2
	@echo >&2
	@echo "ssh-keygen -f ${SSHKEY}"
	@echo >&2
	@exit 1

# Work out which SSH key to use
${SSHKEYID}: ${SSHKEY}
	doctl compute ssh-key list --output json |\
		jq '.[] | select("${SSHKEYNAME}" == .name) | .id' >$@.tmp
	@[[ -s $@.tmp ]] || {\
		echo >&2;\
		echo "${SSHKEY} not known to DigitalOcean as\
			${SSHKEYNAME}" >&2;\
		echo "Please upload ${SSHKEY} to DigitalOcean with" >&2;\
		echo >&2;\
		echo "doctl compute ssh-key create ${SSHKEYNAME}\
			--public-key '$$(cat ${SSHKEY}.pub)'" >&2;\
		echo >&2;\
		exit 1;\
	}
	mv $@.tmp $@

# Roll an SSH command from the info returned from creating a droplet.
${SSHCMD}: ${DROPLETADDR} ${SSHKEY}
	echo -n "${SSHARGS} " >$@.tmp
	cat ${DROPLETADDR} >> $@.tmp
	while $$(cat $@.tmp) -o ConnectTimeout=15 /bin/true 2>&1 |\
		egrep -q 'ssh: connect to host [[:digit:].]+ port ${SSHPORT}:\
		Connection refused'; do\
		sleep 1;\
	done
	$$(cat $@.tmp) cloud-init status --wait
	mv $@.tmp $@

# Work out the user@host for the droplet
${DROPLETADDR}: ${DROPLETINFO}
	jq -r '"root@\(\
		.[0].networks.v4[] | select("public" == .type) | .ip_address\
	)"' $> >$@

# Creates a new droplet and JSON'd metadata in ${DROPLETINFO}.
${DROPLETINFO}: ${SSHKEYID} ${USERDATA}
	doctl compute droplet create\
		--enable-monitoring\
		--image debian-12-x64\
		--output json\
		--project-id ${PROJECTID}\
		--region fra1\
		--size s-1vcpu-512mb-10gb\
		--ssh-keys "$$(cat ${SSHKEYID})"\
		--tag-names ${DROPLETTAG}\
		--user-data-file ${USERDATA}\
		--vpc-uuid ${VPCUUID}\
		--wait\
		${DROPLETNAME} >$@.tmp || { cat $@.tmp && exit 1; }
	mv $@.tmp $@

distclean: clean
	rm -rf ${LOCAL} ${RUNCMD}
clean:
	doctl compute droplet delete --force --tag-name ${DROPLETTAG}
	rm -f ${SSHCMD} ${DROPLETINFO} ${DROPLETADDR} ${LOCAL}/*.tmp
.PHONY: distclean clean
