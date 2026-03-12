#!/bin/sh
# fastmotd — shown at login
# Installed to /etc/profile.d/fastmotd.sh
if [ -x /usr/bin/fastmotd ]; then
    /usr/bin/fastmotd
fi
