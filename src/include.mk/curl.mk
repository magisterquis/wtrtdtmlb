# curl.mk
# Grab curl
# By J. Stuart McMurray
# Created 20241102
# Last Modified 20241102

.ifndef HAVE_CURL_MK
HAVE_CURL_MK=1

.include "debian.mk"

CURL = /usr/bin/curl

${CURL}: .NOTMAIN
	${APTGET} install ca-certificates curl

.endif # .ifndef HAVE_CURL_MK
