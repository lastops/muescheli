#!/bin/sh

# check proxy environment and replace in freshclam.conf
if [ $http_proxy ]; then
    host=$(echo $http_proxy | sed -n "s/^https\?:\/\/\(.*\):.*\?\/\?/\1/p")
    port=$(echo $http_proxy | sed -n "s/^https\?:\/\/.*:\(.*\?\)\/\?/\1/p")

    sed -i "s/#HTTPProxyServer myproxy.com/HTTPProxyServer $host/g" /etc/clamav/freshclam.conf
    sed -i "s/#HTTPProxyPort 1234/HTTPProxyPort $port/g" /etc/clamav/freshclam.conf
fi

freshclam -d --no-dns &
clamd