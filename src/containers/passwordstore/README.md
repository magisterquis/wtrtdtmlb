Password Store
==============
A simple readonly key/value store meant for passwords.

The program expects a single argument, a file containing a JSON object with
password names and values.

Serves up `/`, to list password names, and `/<name>` to retrieve passwords.

The JSON file is re-read every request; it may be changed while the server is
running.
