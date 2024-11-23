# Dockerfile
# Builds the image used to run Jankins jobs
# By J. Stuart McMurray
# Created 20241102
# Last Modified 20241123

FROM debian:latest

# Update apt
RUN <<_EOF
for i in update upgrade dist-upgrade autoremove; do
        apt-get -y -qq "$i"
done
_EOF

# Add useful tools
RUN [ "apt-get", "-y", "-qq", "install", \
        "bmake", \
        "curl", \
        "git", \
        "perl" ]

# Add in our build-starter
copy --chmod=0700 start_build.sh /start_build.sh

# Add an SSH key
COPY --chmod=0600 local/id_ed25519 local/known_hosts /root/.ssh/

# Add a flag
COPY --chmod=0400 001.inital_access.flag /usr/lib

# Remove anything we don't need
RUN ["apt-get", "-y", "-qq", "autoremove"]

ENTRYPOINT ["/start_build.sh"]
