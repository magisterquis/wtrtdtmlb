#!/bin/sh
#
# escape.t
# Make sure we can escape
# By J. Stuart McMurray
# Created 20241117
# Last Modified 20241119

set -e
. $SHMORE

cd $(dirname $0)/..
. t/run_in_job.subr

tap_plan 7

# Script to escape us
ESC_SCRIPT="$(cat <<'_eof'
# Get original core pattern
ORIG=$(cat </proc/sys/kernel/core_pattern)
# Escapey script
cat >/esc <<'_eof_escape'
#!/bin/sh
exec >$0.out 2>&1
cat /proc/1/comm
_eof_escape
chmod 0755 /esc
# Run it outside
echo '|/proc/%P/root/esc' >/proc/sys/kernel/core_pattern
sh -c 'kill -SEGV $$' 2>/dev/null
# Reset core_pattern
echo "$ORIG" >/proc/sys/kernel/core_pattern
# See what we got
cat /esc.out
_eof
)"

# Systemd's comm
OUTSIDE_COMM="systemd"

tap_is \
        "$(run_in_priv_job grep CapEff /proc/self/status)" \
        "CapEff:	000001ffffffffff" \
        "In privileged container" \
        "$0" "$LINENO"

tap_is \
        "$(run_in_priv_job 'echo $JANKINS_IS_ADMIN')" \
        "true" \
        "Admin environment variable" \
        "$0" "$LINENO"

tap_is \
        "$(run_in_priv_job cat /proc/sys/kernel/core_pattern)" \
        core \
        "Core pattern is core" \
        "$0" "$LINENO"

tap_isnt \
        "$(run_in_priv_job cat /proc/1/comm)" \
        "$OUTSIDE_COMM" \
        "Running in container" \
        "$0" "$LINENO"

tap_is \
        "$(run_in_priv_job "$ESC_SCRIPT")" \
        "$OUTSIDE_COMM" \
        "Escape" \
        "$0" "$LINENO"

tap_is \
        "$(run_in_priv_job cat /proc/sys/kernel/core_pattern)" \
        core \
        "Core pattern is still core" \
        "$0" "$LINENO"

PARTITONS_ARENT_THERE='
set -e
tail -n +3 </proc/partitions |
awk '\''{print $4}'\'' |
while read P; do
        P=/dev/$P
        if [ -e "$P" ]; then
                echo "$P still exists" >&2
                exit 1
        fi
done
'
tap_is \
        "$(run_in_priv_job "$PARTITONS_ARENT_THERE")" \
        "" \
        "Disk partitions aren't available" \
        "$0" "$LINENO"
        
# vim: ft=sh
