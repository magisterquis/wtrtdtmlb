#!/bin/sh
# user-data-debian.sh
# Cloud-init user-data for Debian
# By J. Stuart McMurray
# Created 20241007
# Last Modified 20241102

set -e

echo "Starting user-data script"

# Change SSH to port 2
echo "Port 2" >> /etc/ssh/sshd_config
service sshd restart

# Add a 1GB swapfile
fallocate -l 1G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo '/swapfile none swap sw 0 0' >>/etc/fstab

# Update ALL the things
for i in update upgrade autoremove; do
        apt-get -y -qq "$i"
done

# Install just enough to get the rest of the repo here
apt-get -y -qq install rsync bmake
apt-get -y -qq remove bash-completion
