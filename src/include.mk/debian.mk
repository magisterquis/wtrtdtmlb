# debian.mk
# Debian-specific things
# By J. Stuart McMurray
# Created 20241102
# Last Modified 20241123

.ifndef HAVE_DEBIAN_MK
HAVE_DEBIAN_MK=1

STRINGS=/usr/bin/strings

# Quiter installation of things.
APTGET = DEBIAN_FRONTEND=noninteractive >/dev/null apt-get -y -qq\
	 -o DPkg::Lock::Timeout=-1

# Make sure we have a package list before installing things
.BEGIN::
	@${APTGET} update

${STRINGS}: .NOTMAIN
	${APTGET} install binutils

.endif # .ifndef HAVE_DEBIAN_MK
