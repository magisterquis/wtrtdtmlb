#!/bin/sh
#
# passwordstore.t
# Sanity-check the passwordstore
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
        curl --silent --max-time 3"

# get_password gets a password from the password store.
get_password() {
        $CURL \
                --user "$(cat default_username):$(cat default_password)" \
                "http://passwordstore/$1"
}

tap_plan 5

# Can we hit the container?
tap_is \
        "$( $CURL http://passwordstore )" \
        'flag
jankins_admin_password
jankins_admin_username
lsys_ssh_key
lsys_ssh_username
random_thing_password
random_thing_username' \
        "Creds list"\
        "$0" "$LINENO"

# Is password auth required?
tap_is \
        "$($CURL http://passwordstore/dummy)" \
        "Unauthorized" \
        "Auth required" \
        "$0" "$LINENO"

# Can we get a password?
tap_is \
        "$(get_password lsys_ssh_username)" \
        "root" \
        "Get ssh username" \
        "$0" "$LINENO"

# Can we get the Jankins admin password?
subtest() {
        tap_plan 4
        JAU="$(cat ../jankins/default_admin_username)"
        JAP="$(cat ../jankins/default_admin_password)"
        tap_isnt "$JAU" "" "Jankins admin username" "$0" $LINENO
        tap_isnt "$JAP" "" "Jankins admin password" "$0" $LINENO
        tap_is \
                "$(get_password jankins_admin_username)" \
                "$JAU" \
                "Get Jankins admin username" \
                "$0" "$LINENO"
        tap_is \
                "$(get_password jankins_admin_password)" \
                "$JAP" \
                "Get Jankins admin password" \
                "$0" "$LINENO"
}
tap_subtest "Jankins creds" "subtest" "$0" $LINENO

# Make sure the wrong password is rejected
tap_is \
        "$($CURL --user "dummy:dummy" "http://passwordstore/dummy")" \
        "Unauthorized" \
        "Wrong creds rejected" \
        "$0" "$LINENO"

# vim: ft=sh
