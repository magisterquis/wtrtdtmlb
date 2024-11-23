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

# Add in our binary and the passwords file
COPY --chmod=0755 local/passwordstore /passwordstore
COPY --chmod=0400 local/passwords.json /passwords.json
copy --chmod=0400 007.procpidroot.flag /

ENTRYPOINT ["/passwordstore", "passwords.json"]
