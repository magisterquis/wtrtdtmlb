#!/bin/sh
# httpcheckerstart.sh
# Start our container going
# By J. Stuart McMurray
# Created 20241116
# Last Modified 20241116

set -e

# Remove this file
rm $0

# Start sshd
/etc/init.d/ssh start

# Start the HTTP Checker itself
exec /codeserver
