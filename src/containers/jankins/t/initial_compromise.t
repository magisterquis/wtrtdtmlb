#!/bin/sh
#
# initial_compromise.t
# Make sure initial compromise works
# By J. Stuart McMurray
# Created 20241116
# Last Modified 20241119

set -e
. $SHMORE

cd $(dirname $0)/..
. t/run_in_job.subr

tap_plan 10

tap_isnt "$(pgrep jankins)" "" "Jankins is running" "$0" "$LINENO"
tap_is   "$(run_in_job pwd)" "/code" "Command execution" "$0" "$LINENO"
tap_like \
        "$(run_in_job curl -sm1 https://icanhazip.com)" \
        "(?s)^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$" \
        "Comms to internet" \
        "$0" "$LINENO"
tap_is \
        "$(run_in_job 'tr "\0" " " </proc/1/cmdline; echo')" \
        "/sbin/docker-init -- /start_build.sh sh -c exec bmake " \
        "In container" \
        "$0" $LINENO
tap_isnt \
        "$(run_in_job grep CapEff /proc/self/status)" \
        "CapEff:	000001ffffffffff" \
        "Unprivileged container" \
        "$0" "$LINENO"
tap_is \
        "$(run_in_job head -n 1 /root/.ssh/id_ed25519)" \
        "-----BEGIN OPENSSH PRIVATE KEY-----" \
        "Have SSH key" \
        "$0" "$LINENO"
tap_like \
        "$(run_in_job egrep ^codeserver /root/.ssh/known_hosts)" \
        '^codeserver ssh-ed25519' \
        "Have codeserver pubkey fingerprint" \
        "$0" "$LINENO"
tap_is \
        "$(run_in_job which curl)" \
        "/usr/bin/curl" \
        "Have curl" \
        "$0" $LINENO
tap_is \
        "$(run_in_job 'echo $BXT_FD3_SSH_KEY')" \
        ">>Removed from password store<<" \
        "Breadcrumb for password store" \
        "$0" $LINENO
tap_like \
        "$(run_in_job \
                'curl -svm.2 codeserver:{22,80,81} 2>&1 | grep Connected '
        )" \
        '(?s)^'\
'\* Connected to codeserver \(172\.18\.0\.\d+\) port 22 \(#\d+\)\n'\
'\* Connected to codeserver \(172\.18\.0\.\d+\) port 80 \(#\d+\)$' \
        "Code server portscan" \
        "$0" $LINENO

# vim: ft=sh
