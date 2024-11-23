Jankins
=======
Like Jenkins, but just the fun bits.

Setup
-----
1. Have Docker available
2. bmake -I ../../include.mk

Running a Job
-------------
Use the web interface or curl or whatever to upload a txtar archive to `/run`
and wait for output.  Creds are in [`default_username`](./default_username)
 and [`default_password`](./default_password).

The archive will be unpacked into `/code` in a Debian 12 docker container which
has the following packages installed
- BSD Make
- Curl
- Git
- Perl 

The container will run with the working directory set to `/code` and run either
`bmake`, if the archive has no comment, or if the `-allow-custom` flag was
given and the archive has a comment, the archive comment.
