#!/usr/bin/env perl
#
# run_archive.t
# Test running a normal archive
# By J. Stuart McMurray
# Created 20241103
# Last Modified 20241116

use warnings;
use strict;
use Test::More;

use Cwd;

# Archive to run
my $txtar = <<'_eof';
-- Makefile --
test:
	prove --nocount --quiet
-- t/test.t --
#!/bin/sh
WANT=kittens
GOT=$(cat data)
echo 1..1
if [ "$GOT" = "$WANT" ]; then
        echo ok 1 - File contents
else
        echo not ok 1 - File contents
fi
-- data --
kittens
_eof
if ($txtar =~ /'/) {
        die "single quotes in the txtar break this fragile test";
}

# Be in our source directory.
my $srcdir = $0 =~ s,t/[^/]+\.t,,r;
die "could not get source directory from $0" if ($0 eq $srcdir);
if ("" ne $srcdir) {
        die "could not cd to $srcdir" unless chdir $srcdir;
}

# Expect output, per-line
my @pats = (
        'Unpacked archive to /tmp/jankins_jobs/files-127.0.0.1.*',
        'Starting build\.\.\.',
        'prove --nocount --quiet',
        't/test\.t \.\. ok',
        'All tests successful\.',
        'Files=1, Tests=1,\s+\d+ wallclock secs '.
                '\( \d+\.\d+ usr \+\s+\d+\.\d+ sys =\s+\d+\.\d+ CPU\)',
        'Result: PASS',
        'Finished in [0-9.]+[a-z]+',
);
my @res = map{ qr{^$_$} } @pats;

# Test it?
$ENV{TXTAR} = $txtar;
my $cmd = "curl ".
        "--silent ".
        "--insecure ".
        "--user \$(cat default_username):\$(cat default_password) ".
        "--data-urlencode 'archive=$txtar' ".
        "https://127.0.0.1:443";
my @output = `$cmd 2>&1`;

# Make sure we got enough output.
is @output, @res, "Number of lines of output";

# Make sure each line is correct.
for my $i (0 .. $#output) {
        my $l = $i + 1;
        chomp $output[$i];
        like $output[$i], $res[$i], "Output line $l";
}

done_testing;
