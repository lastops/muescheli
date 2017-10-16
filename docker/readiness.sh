#!/bin/sh

# freshclam
#if freshclam | grep -q 'bytecode.* is up to date'; then
#    echo "freshclam running successfully"
#else
#    echo "freshclam not running"
#    exit 1
#fi

# clamd
if clamdscan eicar.com | grep -q 'Infected files: 1'; then
    echo "clamd running successfully"
else
    echo "clamd not running"
    exit 1
fi

exit 0