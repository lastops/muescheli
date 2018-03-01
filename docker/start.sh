#!/bin/sh

# check proxy environment and replace in freshclam.conf
if [ $http_proxy ]; then
    host=$(echo $http_proxy | sed -n "s/^https\?:\/\/\(.*\):.*\?\/\?/\1/p")
    port=$(echo $http_proxy | sed -n "s/^https\?:\/\/.*:\(.*\?\)\/\?/\1/p")

    sed -i "s/^#HTTPProxyServer.*/HTTPProxyServer $host/g" /etc/clamav/freshclam.conf
    sed -i "s/^#HTTPProxyPort.*/HTTPProxyPort $port/g" /etc/clamav/freshclam.conf
fi

# start freshclam and clam in daemon mode
freshclam -d --no-dns &
clamd