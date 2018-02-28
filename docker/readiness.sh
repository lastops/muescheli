#!/bin/sh

# freshclam
if freshclam | grep -Eq 'bytecode.* is up to date|bytecode.* updated'; then
    echo "freshclam running successfully"
else
    echo "freshclam not running"
    exit 1
fi

exit 0