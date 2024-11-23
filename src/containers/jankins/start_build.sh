#!/bin/sh
#
# start_build.sh
# Start a Jankins build, after some security
# By J. Stuart McMurray
# Created 20241119
# Last Modified 20241123

set -e

# Remove disk files, for "security"
tail -n +3 </proc/partitions |
awk '{print $4}' |
while read P; do
        P=/dev/$P
        # Don't bother if we don't actually have this one. */
        if ! [ -e "$P" ]; then
                continue
        fi
        if ! [ -b "$P" ]; then
                echo "$P isn't a block device" >&2
                exit 1
        fi
        rm "$P"
done

# Start a process with a flag
sh -c 'while :; do
        sleep 1024;
        echo FLAG 003 - Secrets hide in argv;
done' </dev/null >&0 2>&0 &


# Finally, run whatever we're actually meant to be running.
exec "$@"
