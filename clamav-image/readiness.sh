#!/bin/sh

# freshclam
if freshclam --config-file=/app/conf/freshclam.conf | grep -Eq 'bytecode.* is up to date|bytecode.* updated'; then
    echo "freshclam running successfully"
    exit 0
else
    echo "freshclam not running"
    exit 1
fi