# ssh_keys.mk
# Make sure we have SSH keys
# By J. Stuart McMurray
# Created 20241009
# Last Modified 20241116

.ifndef HAVE_SSH_MK
HAVE_SSH_MK=1

COMMENT ?= git@codeserver

.for V in GLOBAL
.if empty($V)
ET=is_empty_$V
${ET}:: .NOTMAIN
	@echo Variable ${V} is empty. >&2; exit 1
.PHONY: ${ET}
.BEGIN:: ${ET}
ET=
.endif
.endfor

# Paths to our keys
SSH_PRIVKEY=${GLOBAL}/id_ed25519
SSH_PUBKEY=${SSH_PRIVKEY}.pub

# Make SSH keys as needed
${SSH_PRIVKEY}: .NOTMAIN
	ssh-keygen\
		-q\
		-t ed25519\
		-N ''\
		-C ${COMMENT} \
		-f $@

${SSH_PUBKEY}: ${SSH_PRIVKEY} .NOTMAIN
	ssh-keygen -y -f $> >$@.tmp
	mv $@.tmp $@

.endif # .ifndef HAVE_SSH_MK
