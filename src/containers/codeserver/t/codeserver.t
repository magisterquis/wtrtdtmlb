#!/bin/sh
#
# codeserver.t
# Sanity-check the codeserver
# By J. Stuart McMurray
# Created 20241116
# Last Modified 20241123

set -e

. $SHMORE # TAP library

# Be in our own directory
cd $(dirname $0)/..

# We'll need to do all of this from a docker container, because networking.
CURL="docker run \
        --network jankins \
        --quiet \
        --rm alpine/curl \
        --silent --max-time 3 http://codeserver"

# We'll also run GIT from a docker container, because networking.  Entrypoint
# is $1.
KP=$(realpath $(pwd)/../../../global/id_ed25519)
SC="ssh -o StrictHostKeyChecking=accept-new -q"
GIT() {
        EP=$1
        shift
        docker run \
                --entrypoint "$EP" \
                --env "GIT_SSH_COMMAND=$SC" \
                --network jankins \
                --quiet \
                --rm \
                --volume "$KP:/root/.ssh/id_ed25519" \
                alpine/git \
                "$@" 2>&1
}

tap_plan 5

tap_is \
        "$($CURL)" \
        '<!doctype html>
<meta name="viewport" content="width=device-width">
<pre>
<a href=".ssh/">.ssh/</a>
<a href="central_iac.git/">central_iac.git/</a>
<a href="curlrevshell.git/">curlrevshell.git/</a>
<a href="dtffmacac.git/">dtffmacac.git/</a>
<a href="httpd_botnet_controller.git/">httpd_botnet_controller.git/</a>
<a href="mqd.git/">mqd.git/</a>
<a href="wtrtdtmlb.git/">wtrtdtmlb.git/</a>
</pre>' \
        "List repos" \
        "$0" "$LINENO"

subtest() {
        tap_plan 2
        for RN in central_iac curlrevshell; do
                GIT git clone --quiet git@codeserver:$RN.git
                tap_ok $? "Clone $RN via ssh" "$0" $LINENO
        done
}
tap_subtest "Clone repos" "subtest" "$0" $LINENO

tap_is \
        "$(GIT sh -c 'git clone --quiet git@codeserver:central_iac.git && '\
'egrep -hro "[a-z0-9]+://[^:]+:[^@]+@[a-z0-9.]+" central_iac | sort -u')" \
        "http://jankins_p:l0ng_p4ssw0rds_4r3_n1c3@passwordstore" \
        "Passwordserver password" \
        "$0" "$LINENO"

tap_is \
        "$(GIT sh -c '
                git clone --quiet git@codeserver:central_iac.git &&
                ls central_iac
        ')" \
                '004.repo.flag
Makefile
README.md
digitalocean
get_creds.subr
go.mod
go.sum
src
staticcheck.conf' \
        "IaC Repo Files exist" \
        "$0" "$LINENO"

tap_is \
        "$(GIT $SC git@codeserver whoami)" \
        "git" \
        "Shell in codeserver" \
        "$0" $LINENO

# vim: ft=sh
