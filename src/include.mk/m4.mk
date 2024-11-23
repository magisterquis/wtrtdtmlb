# m4.mk
# Install m4
# By J. Stuart McMurray
# Created 20241115
# Last Modified 20241116

.ifndef HAVE_M4_MK
HAVE_M4_MK=1

.include "debian.mk"

M4=/usr/bin/m4

${M4}: .NOTMAIN
	${APTGET} install m4

.endif # .ifndef HAVE_M4_MK
