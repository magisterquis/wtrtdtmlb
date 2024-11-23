# rsyslog.mk
# Install rsyslog
# By J. Stuart McMurray
# Created 20241019
# Last Modified 20241115

.ifndef HAVE_RSYSLOG_MK
HAVE_RSYSLOG_MK=1

.include "debian.mk"

# Make sure we have the necessary variables
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

SYSLOGD     ?= /usr/sbin/rsyslogd
SYSLOGD_PID  = ${GLOBAL}/syslog.pid

${SYSLOGD}: .NOTMAIN
	${APTGET} install rsyslog
	while ! pgrep rsyslogd >/dev/null; do sleep .1; done
 
${SYSLOGD_PID}: ${SYSLOGD} .NOTMAIN
	while ! pgrep rsyslogd >/dev/null; do sleep .1; done
	pgrep rsyslogd >$@.tmp
	mv $@.tmp $@

.endif # .ifndef HAVE_RSYSLOG_MK
