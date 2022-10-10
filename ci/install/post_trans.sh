#!/bin/sh

# Generated by: gromit policy
# Generated on: Mon Oct 10 17:46:26 UTC 2022

if command -V systemctl >/dev/null 2>&1; then
    if [ ! -f /lib/systemd/system/tyk-gateway.service ]; then
        cp /opt/tyk-gateway/install/inits/systemd/system/tyk-gateway.service /lib/systemd/system/tyk-gateway.service
    fi
else
    if [ ! -f /etc/init.d/tyk-gateway ]; then
        cp /opt/tyk-gateway/install/inits/sysv/init.d/tyk-gateway /etc/init.d/tyk-gateway
    fi
fi
