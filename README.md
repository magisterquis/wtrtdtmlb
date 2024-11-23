What's the Red Team doing to my Linux Box?
==========================================
Material for a [talk](https://docs.google.com/presentation/d/1jg3IHtXxMkbiQq935PzRmc8brnBKgAvxrpSzKluz0Io).

Put the contents of this Archive on a Debian 12 box, grab bmake, cd in, and run
`bmake`.  Tests at the end should pass.

Probably best to use a separate box as it kinda takes over :/

This code.  It's not great..

Infrastructure
--------------
Developed and Tested on Debian 12.  The [`digitalocean`](./digitalocean)
directory has a [Makefile](./digitalocean/Makefile) for spinning up and setting
up a DigitalOcean Droplet.  It may require modification if your VPC UUID isn't
the same as mine.  Apologies.

Initial Access
--------------
A service will be listening for HTTPS connections on port 443 with a
self-signed certificate.

Username: `jankins`
Password: `s3cr3t_p4ssw0rd`

Flags
-----
There are nine flags in the environment, all of the form `FLAG \d{3} - .*`.

Start with just the [initial access](#Initial-Access).  Works out better to
delete the source repo.

Scraping flags from the disk is cheating.  And won't give you all the flags.
