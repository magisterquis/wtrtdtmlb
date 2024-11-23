# m4.mk
# Install m4
# By J. Stuart McMurray
# Created 20241119
# Last Modified 20241119

.ifndef HAVE_ED_MK
HAVE_ED_MK=1

.include "debian.mk"

ED=/usr/bin/ed

${ED}: .NOTMAIN
	${APTGET} install ed

.endif # .ifndef HAVE_ED_MK
