#!/bin/sh
# 
# get_creds.subr
# Grabs creds from the password store
# By J. Stuart McMurray
# Created 20241116
# Last Modified 20241117
 
get_creds() {
        curl \
                --silent      \
                --max-time 10 \
                http://m4_incnonl(../passwordstore/default_username):m4_incnonl(../passwordstore/default_password)@passwordstore/$1"
}
