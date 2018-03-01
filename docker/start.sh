#!/bin/sh

# check proxy environment and replace in freshclam.conf
if [ $http_proxy ]; then
    host=$(echo $http_proxy | sed -n "s/^https\?:\/\/\(.*\):.*\?\/\?/\1/p")
    port=$(echo $http_proxy | sed -n "s/^https\?:\/\/.*:\(.*\?\)\/\?/\1/p")

    sed -i "s/^#HTTPProxyServer.*/HTTPProxyServer $host/g" /app/conf/freshclam.conf
    sed -i "s/^#HTTPProxyPort.*/HTTPProxyPort $port/g" /app/conf/freshclam.conf
fi

# start freshclam and clam in daemon mode
freshclam -d --no-dns --config-file=/app/conf/freshclam.conf &
clamd --config-file=/app/conf/clamd.conf