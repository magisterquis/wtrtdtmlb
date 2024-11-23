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

tap_plan 1
tap_is "$(run_in_job 'whoami')" "root" "Command execution" "$0" "$LINENO"

# vim: ft=sh
