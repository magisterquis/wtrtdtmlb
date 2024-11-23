# git.mk
# Grab git
# By J. Stuart McMurray
# Created 20241116
# Last Modified 20241116

.ifndef HAVE_GIT_MK
HAVE_GIT_MK=1

.include "debian.mk"

GIT = /usr/bin/git

${GIT}: .NOTMAIN
	${APTGET} install git

.endif # .ifndef HAVE_GIT_MK

