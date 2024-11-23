# Dockerfile
# Builds the codeserver
# By J. Stuart McMurray
# Created 20241115
# Last Modified 20241119

FROM debian:latest

# Update apt
RUN <<_EOF
for i in update upgrade dist-upgrade autoremove; do
        apt-get -y -qq "$i"
done
_EOF

# Add useful tools
RUN [ "apt-get", "-y", "-qq", "install", "git", "openssh-server", "curl" ]

# Add git user with shell
RUN [ "useradd", "-r", "-d", "/git", "git" ]
RUN [ "mkdir", "-p", "/git/.ssh" ]
RUN [ "chown", "-R", "git:git", "/git" ]
RUN [ "chmod", "-R", "0700",    "/git" ]

# Add ssh key
copy --chown=git:git --chmod=0600 \
        local/authorized_keys /git/.ssh/authorized_keys

# Add repos
COPY --chown=git:git local/repos /git

# Add the code and start it running
COPY --chown=root:root --chmod=0700 local/codeserver start.sh /
CMD ["/start.sh"]
