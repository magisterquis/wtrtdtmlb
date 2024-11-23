#!/bin/sh
#
# credentials.t
# Get creds from a git repo
# By J. Stuart McMurray
# Created 20241116
# Last Modified 20241116

set -e
. $SHMORE

cd $(dirname $0)/..
. t/run_in_job.subr

tap_plan 1

# Get a list of repos
REPOS=$(run_in_job 'curl -sv $CODESERVER_HTTP')
tap_isnt "$REPOS" "" "Repo list" "$0" $LINENO


# vim: ft=sh
