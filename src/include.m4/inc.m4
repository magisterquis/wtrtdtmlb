m4_dnl inc.m4
m4_dnl Macro to include a file without a newline
m4_dnl By J. Stuart McMurray
m4_dnl Created 20241116
m4_dnl Last Modified 20241116
m4_define(m4_nonl, `m4_patsubst(`$1', `
', `')')m4_dnl
m4_define(m4_incnonl, `m4_nonl(m4_include($1))')m4_dnl
