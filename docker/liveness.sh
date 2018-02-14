#!/bin/sh

# path to script
path=$(cd -P -- "$(dirname -- "$0")" && printf '%s\n' "$(pwd -P)")

# clamd
if clamdscan ${path}/eicar.com | grep -q 'Infected files: 1'; then
    echo "clamd running successfully"
else
    echo "clamd not running"
    exit 1
fi

exit 0